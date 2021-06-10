package storage

import (
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
)

func TestPostgreVault_SaveData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)

		migrateUp()
		defer migrateDown()

		err = db.Ping()
		require.NoError(t, err)

		d := NewPostgreVault(db)
		err = d.SaveData([]byte("k1234"), []byte("value1234"))
		require.NoError(t, err)
	})
	t.Run("success update if try set repeated key in db", func(t *testing.T) {
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)

		migrateUp()
		defer migrateDown()

		err = db.Ping()
		require.NoError(t, err)

		d := NewPostgreVault(db)
		err = d.SaveData([]byte("k1234"), []byte("value1234"))
		require.NoError(t, err)

		err = d.SaveData([]byte("k1234"), []byte("value"))
		require.NoError(t, err)

		data, err := d.ReadData([]byte("k1234"))
		require.NoError(t, err)
		require.EqualValues(t, "value", string(data))
	})
	t.Run("key nil error if try set nil value into db", func(t *testing.T) {
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)

		migrateUp()
		defer migrateDown()

		err = db.Ping()
		require.NoError(t, err)

		d := NewPostgreVault(db)
		err = d.SaveData([]byte("k1234"), []byte("value1234"))
		require.NoError(t, err)

		err = d.SaveData([]byte("k1234"), []byte(""))
		require.NoError(t, err)

		data, err := d.ReadData([]byte("k1234"))
		require.Error(t, err, "postgre: key not found ")
		require.EqualValues(t, []byte(nil), string(data))
	})
	t.Run("error if set nil key", func(t *testing.T) {
		key := ""
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)
		require.NoError(t, err)

		migrateUp()
		defer migrateDown()

		d := NewPostgreVault(db)
		err = d.SaveData([]byte(key), []byte("value1234"))
		require.Error(t, err)
		require.EqualValues(t, "postgre: key can't be nil ", err.Error())
	})
}

func TestPostgreVault_ReadData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)
		require.NoError(t, err)

		migrateUp()
		defer migrateDown()

		err = db.Ping()
		require.NoError(t, err)
		d := NewPostgreVault(db)
		err = d.SaveData([]byte("k1234"), []byte("value1234"))
		require.NoError(t, err)

		data, err := d.ReadData([]byte("k1234"))
		require.NoError(t, err)
		require.EqualValues(t, "value1234", string(data))
	})
	t.Run("error if get by nil key", func(t *testing.T) {
		key := ""
		connStr := "user=postgres password=postgres  sslmode=disable"
		db, err := sqlx.Connect("postgres", connStr)
		require.NoError(t, err)

		migrateUp()
		defer migrateDown()

		d := NewPostgreVault(db)
		data, err := d.ReadData([]byte(key))
		require.Error(t, err)
		require.EqualValues(t, "postgre: key can't be nil ", err.Error())
		require.EqualValues(t, []byte(nil), data)
	})
}

func migrateUp() {
	m, err := migrate.New(
		"file://../../../scripts/migrations",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	err = m.Up()
	if err != nil {
		fmt.Println(err)
	}
}

func migrateDown() {
	mdown, err := migrate.New(
		"file://../../../scripts/migrations",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	err = mdown.Down()
	if err != nil {
		fmt.Println(err)
	}
}
