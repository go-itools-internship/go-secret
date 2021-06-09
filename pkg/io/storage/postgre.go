package storage

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgreVault struct {
	db *sqlx.DB
}

var schema = `CREATE TABLE postgres (
	key text NOT NULL,
	value text NOT NULL,
	PRIMARY KEY (key)
);`

type store struct {
	key   string
	value string
}

// NewPostgreVault create new postgreSQL  client
func NewPostgreVault(p *sqlx.DB) *postgreVault {
	pv := &postgreVault{
		db: p,
	}
	return pv
}

// SaveData put data in postgres storage by key and encoded value
// 	key to set in postgres storage
// 	encoded value to storage
func (r *postgreVault) SaveData(key, encodedValue []byte) error {
	ctx := context.Background()
	tx, _ := r.db.BeginTxx(ctx, nil)
	defer func() {
		if err := r.db.Close(); err != nil {
			fmt.Println("postgre: can't close db: ", err.Error())
		}
	}()
	r.db.MustBegin()
	tx.MustExec("INSERT INTO postgres (key , value) VALUES ($1,$2)", string(key), string(encodedValue))
	err := tx.Commit() // TODO fix error handler
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
	}
	return nil
}

// ReadData get data from postgres storage by key
// 	key to get value for pair key-value from postgres storage
func (r *postgreVault) ReadData(key []byte) ([]byte, error) {
	ctx := context.Background()
	var val []string
	err := r.db.SelectContext(ctx, &val, "SELECT value FROM postgres WHERE key=$1 LIMIT 1", string(key))
	defer func() {
		if err := r.db.Close(); err != nil {
			fmt.Println("postgre: can't close db: ", err.Error())
		}
	}()
	fmt.Println(val) //TODO for testing (delete before pull request)

	return []byte(val[0]), err
}
