package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Proteins     float64 `json:"proteins"`
	Fats         float64 `json:"fats"`
	Carbohydrates float64 `json:"carbohydrates"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./products.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						name TEXT,
						proteins FLOAT,
						fats FLOAT,
						carbohydrates FLOAT);`)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	// add more routes as necessary...

	fmt.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", router)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// insert into database
	_, err = db.Exec("INSERT INTO products (name, proteins, fats, carbohydrates) VALUES (?, ?, ?, ?)", product.Name, product.Proteins, product.Fats, 
product.Carbohydrates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return created product as JSON
	json.NewEncoder(w).Encode(&product)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	err := db.QueryRow("SELECT id, name, proteins, fats, carbohydrates FROM products WHERE id = ?", id).Scan(&product.ID, &product.Name, &product.Proteins,
&product.Fats, &product.Carbohydrates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// return product as JSON
	json.NewEncoder(w).Encode(&product)
}
