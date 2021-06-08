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
	db, err := sqlx.Connect("postgres", "name=test")
	if err != nil {
		fmt.Println(err, "err")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	//d := NewPostgreVault(db)
	//err = d.SaveData([]byte("key"), []byte("value"))
	//
	//require.NoError(t, err)
	//data, err := d.ReadData([]byte("key"))
	//require.NoError(t, err)
	//fmt.Println(data)
}

func TestPostgreVault_SaveData(t *testing.T) {

	db, err := sqlx.Connect("postgres", "name=test")
	if err != nil {
		fmt.Println(err, "err")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	d := NewPostgreVault(db)
	err = d.SaveData([]byte("key"), []byte("value"))
	require.NoError(t, err)

}

func TestPostgreVault_ReadData(t *testing.T) {
	db, err := sqlx.Connect("postgres", "name=test")
	if err != nil {
		fmt.Println(err, "err")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	d := NewPostgreVault(db)
	err = d.SaveData([]byte("key"), []byte("value"))
	require.NoError(t, err)

	data, err := d.ReadData([]byte("key"))
	require.NoError(t, err)
	fmt.Println(data)
}
