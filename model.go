package main

import (
	"database/sql"
	"errors"
)

type product struct {
	ID    int     `json: "id"`
	Name  string  `json: "name"`
	Price float64 `json: "price"`
}

func (p *product) getPruduct(db *sql.DB) error {
	return errors.New("not implemented")
}

func (p *product) updateProduct(db *sql.DB) error {
	return errors.New("not implemented")
}

func (p *product) udeleteProduct(db *sql.DB) error {
	return errors.New("not implemented")
}
func (p *product) createProduct(db *sql.DB) error {
	return errors.New("not implemented")
}

func getPruduct(db *sql.DB, start, count int) ([]product, error) {
	return nil, errors.New("not implemented")
}
