package api

import (
	"encoding/json"
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	models "github.com/DelegacionUC3M/leisten/models"
	"strconv"
)

// APIHandler contains a pointer to the database
type APIHandler struct {
	DB *gorm.DB
}

// GetItems returns the item with the given id
func (API *APIHandler) ListItem(c *gin.Context) {
	itemID := c.Param("id")

	var (
		count int
		item  models.Item
	)

	API.DB.Model(&models.Item{}).
		Where("id = ?", itemID).Count(&count)

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Item does not exist",
		})
	} else {
		API.DB.First(&item, itemID)
		c.JSON(http.StatusOK, item)
	}
}

// GetAllItems returns all available items
//
// With no query parameters returns all objects in the database
// Due to limitations of Gin we cannot have an entry point like /items/depleted to
// query objects with Amount = 0. The user has to create that specific query
// The solution I propose is to create a map of all the query parameters
// If the parameter is not nil, add it to the list and perform the query
func (API *APIHandler) GetAllItems(c *gin.Context) {

	var itemsList []models.Item
	queryHolder := make(map[string]interface{})

	amount := c.Query("amount")

	if amount != "" {
		queryHolder["amount"], _ = strconv.Atoi(amount)
	}

	API.DB.Where(queryHolder).Find(&itemsList)

	// TODO: Think about what this endpoint should return
	if len(itemsList) > 0 {
		c.JSON(http.StatusOK, itemsList)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find items with the requirements",
		})
	}
}

// CreateItems inserts new items into the database
func (API *APIHandler) CreateItems(c *gin.Context) {
	var data models.Item

	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	API.DB.Create(&data)
	c.JSON(http.StatusOK, gin.H{
		"message": "New item created",
		"id":      data.ID,
	})
}

// DeleteItem deletes one item from the database
func (API *APIHandler) DeleteItem(c *gin.Context) {
	itemID := c.Param("id")

	var count int
	API.DB.Model(&models.Item{}).
		Where("id = ?", itemID).Count(&count)

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Item does not exist",
		})
	} else {
		API.DB.Where("id = ?", itemID).Delete(&models.Item{})
		c.JSON(http.StatusOK, gin.H{
			"message": "Item deleted correctly",
		})
	}
}

// UpdateItem updates the parameter of the given item
func (API *APIHandler) UpdateItem(c *gin.Context) {
	itemID := c.Param("id")

	var (
		count int
		item  models.Item
	)

	API.DB.Model(&item).
		Where("id = ?", itemID).Count(&count)

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The item does not exist",
		})
	} else {
		// TODO: Handle update
		// item.Amount = 100
		// API.DB.Save(&item)
		c.JSON(http.StatusOK, item)
	}

}

// GetLoans returns all available loans
func (API *APIHandler) GetAllLoans(c *gin.Context) {
	loansRows, err := API.DB.Model(&models.Loan{}).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get loans",
		})
	}
	defer loansRows.Close()

	var loan models.Loan
	var loanList []models.Loan

	for loansRows.Next() {
		API.DB.ScanRows(loansRows, &loan)
		loanList = append(loanList, loan)
	}

	c.JSON(http.StatusOK, loanList)
}

// CreateLoan creates an student loan for the given user id
func (API *APIHandler) CreateLoan(c *gin.Context) {
	var data models.Loan

	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	API.DB.Create(&data)
	c.JSON(http.StatusOK, gin.H{
		"message": "New loan created",
		"id":      data.ID,
	})
}
