package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/stretchr/testify/require"
)

func TestNewPostgreVault(t *testing.T) {
	connStr := "user=postgres password=postgres  sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(t, err)
	r := db.MustExec(schema)
	require.NotEmpty(t, r)
	err = db.Ping()
	require.NoError(t, err)
}

func TestPostgreVault_SaveData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)
		if err != nil {
			fmt.Println(err, "err")
		}
		err = db.Ping()
		require.NoError(t, err)

		d := NewPostgreVault(db)
		err = d.SaveData([]byte("k123"), []byte("value123"))
		require.NoError(t, err)
	})
}

func TestPostgreVault_ReadData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)
		if err != nil {
			fmt.Println(err, "err")
		}
		err = db.Ping()
		require.NoError(t, err)
		d := NewPostgreVault(db)
		data, err := d.ReadData([]byte("k123"))
		require.NoError(t, err)
		require.EqualValues(t, "value123", string(data))
	})

}
