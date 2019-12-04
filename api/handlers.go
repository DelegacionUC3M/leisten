package api

import (
	"encoding/json"
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"

	models "github.com/DelegacionUC3M/leisten/models"
)

// APIHandler contains a pointer to the database
type APIHandler struct {
	DB *gorm.DB
}

// GetItems returns all available items
func (API *APIHandler) ListItems(c *gin.Context) {
	itemsRows, err := API.DB.Model(&models.Item{}).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get items",
		})
	}
	defer itemsRows.Close()

	var item models.Item
	var itemsList []models.Item

	for itemsRows.Next() {
		API.DB.ScanRows(itemsRows, &item)
		itemsList = append(itemsList, item)
	}

	c.JSON(http.StatusOK, itemsList)
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
	if count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Item does not exist",
		})
	} else {
		API.DB.Where("id = ?", itemID).Delete(&models.Item{})
		c.JSON(http.StatusOK, gin.H{
			"message": itemID,
		})
	}
}

// GetLoans returns all available loans
func (API *APIHandler) GetLoans(c *gin.Context) {
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
