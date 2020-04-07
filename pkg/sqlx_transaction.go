package pkg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SQLXTransaction struct {
	db *sqlx.DB
}

func NewSQLXTransaction(db *sqlx.DB) *SQLXTransaction {
	return &SQLXTransaction{
		db: db,
	}
}

func (repo *SQLXTransaction) CreateTransaction(ctx context.Context, options sql.TxOptions) (*sqlx.Tx, error) {
	tx, err := repo.db.BeginTxx(ctx, &options)
	if err != nil {
		return nil, err
	}

	return tx, err
}

func (repo *SQLXTransaction) RollbackTransaction(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func (repo *SQLXTransaction) CommitTransaction(tx *sqlx.Tx) error {
	return tx.Commit()
}
