package middleware

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/microsoft/go-mssqldb"
)

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
		ctx := context.WithValue(r.Context(), "tx", tx)
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			}
		}()
		next.ServeHTTP(w, r.WithContext(ctx))
		tx.Commit()
	})
}

func GetTx(ctx context.Context) *sql.Tx {
	return ctx.Value("tx").(*sql.Tx)
}
