package storage

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	postgreURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	migration  = "file://../../../scripts/migrations"
)

func TestPostgreVault_SaveData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		migrateUp(t)
		defer migrateDown(t)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgreURL)
		require.NoError(t, err)
		defer disconnectPDB(db, t)

		err = db.Ping()
		require.NoError(t, err)

		d := NewPostgreVault(db)
		err = d.SaveData([]byte("k1234"), []byte("value1234"))
		require.NoError(t, err)
	})
	t.Run("success update if try set repeated key in db", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		migrateUp(t)
		defer migrateDown(t)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgreURL)
		require.NoError(t, err)
		defer disconnectPDB(db, t)

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		migrateUp(t)
		defer migrateDown(t)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgreURL)
		require.NoError(t, err)
		defer disconnectPDB(db, t)

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		migrateUp(t)
		defer migrateDown(t)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgreURL)
		require.NoError(t, err)
		defer disconnectPDB(db, t)

		d := NewPostgreVault(db)
		err = d.SaveData([]byte(key), []byte("value1234"))
		require.Error(t, err)
		require.EqualValues(t, "postgres: key can't be nil", err.Error())
	})
}

func TestPostgreVault_ReadData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		migrateUp(t)
		defer migrateDown(t)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgreURL)
		defer disconnectPDB(db, t)
		require.NoError(t, err)

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
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		migrateUp(t)
		defer migrateDown(t)

		db, err := sqlx.ConnectContext(ctx, "postgres", postgreURL)
		defer disconnectPDB(db, t)
		require.NoError(t, err)

		d := NewPostgreVault(db)
		data, err := d.ReadData([]byte(key))
		require.Error(t, err)
		require.EqualValues(t, "postgres: key can't be nil ", err.Error())
		require.EqualValues(t, []byte(nil), data)
	})
}

func migrateUp(t *testing.T) {
	m, err := migrate.New(
		migration,
		postgreURL)
	if err != nil {
		t.Log(err)
	}
	err = m.Up()
	if err != nil {
		t.Log(err)
	}
}

func migrateDown(t *testing.T) {
	m, err := migrate.New(
		migration,
		postgreURL)
	if err != nil {
		t.Log(err)
	}
	err = m.Down()
	if err != nil {
		t.Log(err)
	}
}

func disconnectPDB(pdb *sqlx.DB, t *testing.T) {
	err := pdb.Close()
	if err != nil {
		t.Log("can't disconnect postgres db")
	}
}
