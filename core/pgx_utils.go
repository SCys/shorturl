package core

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func PGXUpdateHelper(ctx context.Context, name string, target interface{}, data H, db *pgxpool.Pool) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	for key, name := range data {
		query := fmt.Sprintf("update %s set %s=$2 where id=$1", name, key)

		if _, err := tx.Exec(ctx, query, name); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}

	}

	err = tx.Commit(ctx)
	return err
}
