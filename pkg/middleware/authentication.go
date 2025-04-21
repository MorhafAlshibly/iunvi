package middleware

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	_ "github.com/microsoft/go-mssqldb"
)

type Authentication struct {
	audience string
	jwks     string
}

func WithAudience(audience string) func(*Authentication) {
	return func(input *Authentication) {
		input.audience = audience
	}
}

func WithJWKS(jwks string) func(*Authentication) {
	return func(input *Authentication) {
		input.jwks = jwks
	}
}

func NewAuthentication(options ...func(*Authentication)) *Authentication {
	auth := &Authentication{
		jwks: "https://login.microsoftonline.com/common/discovery/v2.0/keys",
	}
	for _, option := range options {
		option(auth)
	}
	return auth
}

// Struct to match JWKS response from Azure
type JWK struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Middleware function for authentication
func (a Authentication) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		// Parse and validate JWT
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Fetch JWKS from Azure AD
			keys, err := fetchJWKS(a.jwks)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
			}

			// Get the "kid" (Key ID) from the token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("kid not found in token header")
			}

			// Find the matching public key
			for _, key := range keys {
				if key.Kid == kid {
					return convertJWKToPublicKey(key)
				}
			}

			return nil, errors.New("matching public key not found")
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Validate claims
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Verify audience
		if !claims.VerifyAudience(a.audience, true) {
			http.Error(w, "Invalid audience", http.StatusUnauthorized)
			return
		}

		tenantDirectoryId := (*claims)["tid"].(string)
		userObjectId := (*claims)["oid"].(string)
		appRoles := (*claims)["roles"].([]interface{})
		if len(appRoles) == 0 {
			http.Error(w, "No roles found", http.StatusUnauthorized)
			return
		}
		var highestRole string
		for _, role := range appRoles {
			roleString := role.(string)
			if highestRole == "" || strings.HasPrefix(roleString, "Admin") {
				highestRole = roleString
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "accessToken", tokenString)

		// Set context TenantDirectoryId, UserObjectId and TenantRole
		ctx = context.WithValue(ctx, "TenantDirectoryId", tenantDirectoryId)
		ctx = context.WithValue(ctx, "UserObjectId", userObjectId)
		ctx = context.WithValue(ctx, "TenantRole", highestRole)

		tx := GetTx(ctx)

		// Set Tenant ID in SQL session context
		_, err = tx.ExecContext(ctx,
			"EXEC sp_set_session_context @key = N'TenantDirectoryId', @value = @TenantDirectoryId;",
			sql.Named("TenantDirectoryId", tenantDirectoryId),
		)
		if err != nil {
			http.Error(w, "Failed to set tenant context", http.StatusInternalServerError)
			return
		}
		_, err = tx.ExecContext(ctx,
			"EXEC sp_set_session_context @key = N'UserObjectId', @value = @UserObjectId;",
			sql.Named("UserObjectId", userObjectId),
		)
		if err != nil {
			http.Error(w, "Failed to set user context", http.StatusInternalServerError)
			return
		}
		_, err = tx.ExecContext(ctx,
			"EXEC sp_set_session_context @key = N'TenantRole', @value = @TenantRole;",
			sql.Named("TenantRole", highestRole),
		)
		if err != nil {
			http.Error(w, "Failed to set role context", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Fetches JWKS keys from Azure AD
func fetchJWKS(jwksURL string) ([]JWK, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Keys []JWK `json:"keys"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Keys, nil
}

// Converts a JWK key (modulus + exponent) to an RSA public key
func convertJWKToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode Base64 URL-encoded modulus (N)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %v", err)
	}

	// Decode Base64 URL-encoded exponent (E)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	// Convert exponent to an integer
	var e int
	if len(eBytes) == 3 {
		e = int(binary.BigEndian.Uint32(append([]byte{0}, eBytes...)))
	} else {
		e = int(binary.BigEndian.Uint32(eBytes))
	}

	// Construct RSA public key
	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}

	return pubKey, nil
}

type RoleBasedAccessControl struct{}

// Middleware to enforce specific roles
func (r RoleBasedAccessControl) Middleware(roles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*jwt.MapClaims)
		if !ok {
			http.Error(w, "Claims not found", http.StatusUnauthorized)
			return
		}

		roleSet := make(map[string]struct{})
		for _, role := range roles {
			roleSet[role] = struct{}{}
		}

		hasRole := false
		for _, userRole := range (*claims)["roles"].([]interface{}) {
			if _, ok := roleSet[userRole.(string)]; ok {
				hasRole = true
				break
			}
		}

		if !hasRole {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
