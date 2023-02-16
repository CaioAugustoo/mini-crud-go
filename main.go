package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		Id:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func InsertProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("insert into products(id, name, price) values(?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Id, product.Name, product.Price)

	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("update products set name = ?, price = ? where id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.Id)

	if err != nil {
		return err
	}

	return nil
}

func GetProductById(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("select id, name, price from products where id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	p := Product{}

	err = stmt.QueryRow(id).Scan(&p.Id, &p.Name, &p.Price)

	if err != nil {
		return nil, err
	}

	return &p, err
}

func GetAllProducts(db *sql.DB) (*[]Product, error) {
	rows, err := db.Query("select id, name, price from products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		p := Product{}

		err = rows.Scan(&p.Id, &p.Name, &p.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return &products, nil
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/root")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	myProduct := NewProduct("Nike sneaker", 500.00)

	InsertProduct(db, myProduct)
	UpdateProduct(db, &Product{Id: "3f162eeb-f7d7-4912-952a-b0f00393050d", Name: "Updated nike sneaker", Price: 1299})

	// p, err := GetProductById(db, "5edd8ed3-0b12-4e1d-abf8-ae37b083a27d")
	allProducts, err := GetAllProducts(db)

	if err != nil {
		panic(err)
	}

	data := json.NewEncoder(os.Stdout).Encode(allProducts)
	fmt.Println(data)
}
