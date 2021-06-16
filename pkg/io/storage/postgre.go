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
	fmt.Println("postgres-context ", ctx)
	if bytes.Equal(key, []byte("")) {
		return errors.New("postgres: key can't be nil ")
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	fmt.Println("transaction-begin", tx)
	if err != nil {
		return fmt.Errorf("postgres: can't begin transaction  %w", err)
	}
	if bytes.Equal(encodedValue, []byte("")) {
		_, err = tx.ExecContext(ctx, "DELETE FROM postgres WHERE key=$1;", string(key))
		if err != nil {
			rErr := tx.Rollback()
			if rErr != nil {
				return fmt.Errorf("postgres: can't rollback %v %v", err, rErr)
			}
			return fmt.Errorf("postgres: can't delete data %w", err)
		}
	} else {
		data, err := tx.ExecContext(ctx, "INSERT INTO postgres (key , value) VALUES ($1,$2) ON CONFLICT (key) DO UPDATE SET value=$2;", string(key), hex.EncodeToString(encodedValue))
		fmt.Println("try insert", data)
		if err != nil {
			rErr := tx.Rollback()
			if rErr != nil {
				return fmt.Errorf("postgres: can't rollback %v %v", err, rErr)
			}
			return fmt.Errorf("postgres: can't insert data %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil {
			return fmt.Errorf("postgres: can't rollback %v %v", err, rErr)
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
		return nil, errors.New("postgres: key can't be nil ")
	}
	var val []struct {
		Value string `db:"value"`
	}
	err := r.db.SelectContext(ctx, &val, "SELECT value FROM postgres WHERE key=$1 LIMIT 1;", string(key))
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	if len(val) == 0 {
		return nil, errors.New("postgres: key not found ")
	}
	value, err := hex.DecodeString(val[0].Value)
	if err != nil {
		return nil, fmt.Errorf("postgres: cant't decode value %w", err)
	}
	return value, nil
}
