package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	api "github.com/DelegacionUC3M/leisten/api"
	models "github.com/DelegacionUC3M/leisten/models"
	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	tomlFile = "config.toml"
)

type Database struct {
	Name     string
	User     string
	Password string
}

type Databases struct {
	Loans Database
}

func main() {
	var config Databases
	_, err := toml.DecodeFile(tomlFile, &config)

	fmt.Println(config.Loans)

	conn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", config.Loans.User, config.Loans.Password, config.Loans.Name)
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.Item{})

	r := gin.Default()
	APIHandler := new(api.APIHandler)
	APIHandler.DB = db

	itemsGroup := r.Group("/items")
	{
		itemsGroup.GET("/list", APIHandler.ListItems)
		itemsGroup.POST("/create", APIHandler.CreateItems)
		itemsGroup.GET("/delete/:id", APIHandler.DeleteItem)
	}

	loansGroup := r.Group("/loans")
	{
		loansGroup.GET("/list", APIHandler.GetLoans)
	}

	r.Run()
}
