package middleware

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/microsoft/go-mssqldb"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	committed  bool
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.committed = true
	r.ResponseWriter.WriteHeader(code)
}

type Transaction struct {
	db *sql.DB
}

func WithDB(db *sql.DB) func(*Transaction) {
	return func(input *Transaction) {
		input.db = db
	}
}

func NewTransaction(options ...func(*Transaction)) *Transaction {
	transaction := &Transaction{}
	for _, option := range options {
		option(transaction)
	}
	return transaction
}

func (t *Transaction) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := t.db.Begin()
		if err != nil {
			http.Error(w, "failed to start transaction", http.StatusInternalServerError)
			return
		}

		// Create a response recorder to track the HTTP status code
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Inject transaction into context
		ctx := context.WithValue(r.Context(), "tx", tx)
		r = r.WithContext(ctx)

		// Ensure rollback by default, commit only if everything is okay
		defer func() {
			if rec.statusCode >= 400 { // If an error response was written, rollback
				tx.Rollback()
				return
			}
			if rec.committed { // If already committed, don't commit again
				return
			}
			if err := tx.Commit(); err != nil {
				http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
			}
		}()

		// Handle panic safely
		defer func() {
			if rec := recover(); rec != nil {
				tx.Rollback()
				panic(rec) // rethrow panic after rollback
			}
		}()

		next.ServeHTTP(rec, r)
	})
}

func GetTx(ctx context.Context) *sql.Tx {
	return ctx.Value("tx").(*sql.Tx)
}
