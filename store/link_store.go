package store

import (
	"database/sql"
)

type LinkStore struct {
	db *sql.DB
}

type Link struct {
	ID  int
	URL string
}

func NewLinkStore(db *sql.DB) *LinkStore {
	return &LinkStore{db}
}

func (ls *LinkStore) GetByID(id int) (*Link, error) {
	var link Link
	err := ls.db.QueryRow("SELECT id, url FROM urls WHERE id = $1", id).Scan(&link.ID, &link.URL)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (ls *LinkStore) GetByURL(url string) (*Link, error) {
	var link Link
	err := ls.db.QueryRow("SELECT id, url FROM urls WHERE url = $1", url).Scan(&link.ID, &link.URL)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (ls *LinkStore) PersistURL(url string) (*Link, error) {
	var newID int
	err := ls.db.QueryRow("INSERT INTO urls (url) VALUES($1) returning id", url).Scan(&newID)
	if err != nil {
		return nil, err
	}

	return &Link{
		ID:  newID,
		URL: url,
	}, nil
}
