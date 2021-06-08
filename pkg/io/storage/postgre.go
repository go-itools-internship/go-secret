package storage

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgreVault struct {
	db *sqlx.DB
}

var schema = `CREATE TABLE vault (
	key text,
	value text,
)`

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

func (r *postgreVault) SaveData(key, encodedValue []byte) error {
	tx := r.db.MustBegin()
	tx.MustExec("INSERT INTO vault (key, value) values ($1,$2)", key, encodedValue)
	err := tx.Commit()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *postgreVault) ReadData(key []byte) ([]byte, error) {
	tx := r.db.MustBegin()
	data := store{}
	err := tx.Select(&data, "SELECT value from vault WHERE key=$1", string(key))
	fmt.Println(data.value)
	err = tx.Commit()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return nil, err
		}
	}
	return []byte(data.value), err
}
