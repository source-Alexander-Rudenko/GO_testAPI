package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func TestMain(m *testing.M) {
	a.Initialazer()

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("DELETE FROM products")
	if err != nil {
		log.Fatal(err)

	}
	_, err = a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNotExistedProduct(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/products/1/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		log.Default()
	}
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["eror"])
	}

}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{"name":"test product", "price": 11.33}`)
	req, _ := http.NewRequest("GET", "/products", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "test product" {
		t.Errorf("expected test product, got %v", m["name"])
	}
	if m["price"] != 11.33 {
		t.Errorf("expected price = '11.33', got %v", m["price"])
	}
	if m["id"] != 1.0 {
		t.Errorf("expected id '1', got %v", m["id"])
	}

}
func addProduct(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO pruducts(name, price) VALUES ($1, $2)", "products"+strconv.Itoa(i), (i+1.0)*10)
		if err != nil {
			log.Fatal("error in db inset")
		}

	}
}
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct(1)
	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct(1)
	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name": "test product updated", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}
func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct(1)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/products/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/products/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
