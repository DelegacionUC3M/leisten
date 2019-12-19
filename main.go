package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	api "github.com/DelegacionUC3M/leisten/api"
	models "github.com/DelegacionUC3M/leisten/models"
	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	tomlFile = "config.toml"
	port     = ":8000"
)

// Database holds the parameters necessary to connect to a database
type Database struct {
	Name     string
	User     string
	Password string
}

// Databases holds the content of the toml config file
type Databases struct {
	Loans Database
}

func main() {
	var config Databases
	_, err := toml.DecodeFile(tomlFile, &config)

	// conn := fmt.Sprintf("host=db_postgres user=%s password=%s dbname=%s sslmode=disable",
	conn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.Loans.User, config.Loans.Password, config.Loans.Name)

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.Item{}, &models.Loan{})

	Handler := new(api.Handler)
	Handler.DB = db

	r := mux.NewRouter().StrictSlash(true)

	itemRouter := r.PathPrefix("/prestamos/items").Subrouter()

	itemRouter.HandleFunc("/depleted", Handler.GetDepletedItems).Methods("GET")
	itemRouter.HandleFunc("/{itemID}", Handler.ListItem).Methods("GET")
	itemRouter.HandleFunc("/{itemID}", Handler.UpdateItem).Methods("PUT")
	itemRouter.HandleFunc("/{itemID}", Handler.DeleteItem).Methods("DELETE")
	itemRouter.HandleFunc("", Handler.CreateItem).Methods("POST")
	itemRouter.HandleFunc("", Handler.GetAllItems).Methods("GET")

	loanRouter := r.PathPrefix("/prestamos/loans").Subrouter()

	loanRouter.HandleFunc("", Handler.GetAllLoans).Methods("GET")
	loanRouter.HandleFunc("", Handler.CreateLoan).Methods("POST")

	fmt.Println("Starting server on ", port)
	log.Fatal(http.ListenAndServe(port, r))

}
