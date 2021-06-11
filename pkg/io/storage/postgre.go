package storage

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgreVault struct {
	db *sqlx.DB
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
	if bytes.Equal(key, []byte("")) {
		return errors.New("postgre: key can't be nil ")
	}
	tx, _ := r.db.BeginTxx(ctx, nil)
	r.db.MustBegin()
	if bytes.Equal(encodedValue, []byte("")) {
		tx.MustExecContext(ctx, "DELETE FROM postgres WHERE key=$1", string(key))
	} else {
		tx.MustExecContext(ctx, "INSERT INTO postgres (key , value) VALUES ($1,$2) ON CONFLICT (key) DO UPDATE SET value=$2", string(key), hex.EncodeToString(encodedValue))
	}
	err := tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("postgres: can't rollback %w", err)
		}
		return fmt.Errorf("postgres: can't commit %w", err)
	}
	return nil
}

// ReadData get data from postgres storage by key
// 	key to get value for pair key-value from postgres storage
func (r *postgreVault) ReadData(key []byte) ([]byte, error) {
	ctx := context.Background()
	if bytes.Equal(key, []byte("")) {
		return nil, errors.New("postgre: key can't be nil ")
	}
	var val []struct {
		Value string `db:"value"`
	}
	err := r.db.SelectContext(ctx, &val, "SELECT value FROM postgres WHERE key=$1 LIMIT 1", string(key))
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	if len(val) == 0 {
		return nil, errors.New("postgre: key not found ")
	}
	value, err := hex.DecodeString(val[0].Value)
	if err != nil {
		return nil, fmt.Errorf("postgres: cant't decode value %w", err)
	}
	return value, nil
}
