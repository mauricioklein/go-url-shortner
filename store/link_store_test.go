package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestLinkStore_GetByID(t *testing.T) {
	db, err := testDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer cleanupDB(db)

	ls := NewLinkStore(db)

	// insert registry in database
	var id int
	db.QueryRow("INSERT INTO links (url) values ('http://www.google.com/') returning id").Scan(&id)

	// query the record
	link, err := ls.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, id, link.ID)
	assert.Equal(t, "http://www.google.com/", link.URL)
}

func TestLinkStore_GetByURL(t *testing.T) {
	db, err := testDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer cleanupDB(db)

	ls := NewLinkStore(db)

	// insert registry in database
	var id int
	db.QueryRow("INSERT INTO links (url) values ('http://www.google.com/') returning id").Scan(&id)

	// query the record
	link, err := ls.GetByURL("http://www.google.com/")
	assert.NoError(t, err)
	assert.Equal(t, id, link.ID)
	assert.Equal(t, "http://www.google.com/", link.URL)
}

func TestLinkStore_PersistURL(t *testing.T) {
	db, err := testDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer cleanupDB(db)

	ls := NewLinkStore(db)

	// persist url
	newLink, err := ls.PersistURL("http://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, "http://www.google.com", newLink.URL)

	// check if the record was actually persisted
	var id, url string
	db.QueryRow("SELECT id, url from links where id = $1", newLink.ID).Scan(&id, &url)
	assert.Equal(t, id, strconv.Itoa(newLink.ID))
	assert.Equal(t, url, newLink.URL)
}

func testDBConnection() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", "db", "postgres", "test")
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func cleanupDB(db *sql.DB) {
	db.Exec("TRUNCATE TABLE links")
}
