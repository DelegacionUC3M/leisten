package api

import (
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"

	models "github.com/DelegacionUC3M/leisten/models"
	// "strconv"
)

// Handler contains a pointer to the database
type Handler struct {
	DB *gorm.DB
}

// ListItem returns the item with the requested id
func (API *Handler) ListItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	itemID := mux.Vars(r)["itemID"]

	var (
		count int
		item  models.Item
	)

	API.DB.Model(&models.Item{}).
		Where("id = ?", itemID).Count(&count)

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{"message": "Item does not exist."})
	} else {
		API.DB.First(&item, itemID)

		w.WriteHeader(http.StatusOK)
		payload, _ := json.Marshal(item)
		w.Write(payload)
	}
}

// DeleteItem deletes the item with the requested id from the database
func (API *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	itemID := mux.Vars(r)["itemID"]

	var count int
	API.DB.Model(&models.Item{}).
		Where("id = ?", itemID).Count(&count)

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{"message": "Item does not exist."})
	} else {
		API.DB.Where("id = ?", itemID).Delete(&models.Item{})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Item deleted correctly"})
	}
}

// GetAllItems returns all items available in the database
func (API *Handler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		itemsList []models.Item
		item      models.Item
	)

	itemsRows, err := API.DB.Model(models.Item{}).Rows()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Could not fetch items"})
		return
	}
	for itemsRows.Next() {
		API.DB.ScanRows(itemsRows, &item)
		itemsList = append(itemsList, item)
	}

	w.WriteHeader(http.StatusOK)
	payload, _ := json.Marshal(itemsList)
	w.Write(payload)
}

// GetDepletedItems returns all items whose amount is 0
func (API *Handler) GetDepletedItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var itemsList []models.Item

	API.DB.Where("amount = 0").Find(&itemsList)
	w.WriteHeader(http.StatusOK)
	payload, _ := json.Marshal(itemsList)
	w.Write(payload)
}

// CreateItems creates a new item and inserts it into the database
func (API *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data models.Item

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Could not create item"})
	} else {
		API.DB.Create(&data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]uint{"id": data.ID})
	}
}

// UpdateItem updates the item with the given id based on the payload
// The payload must be a full new object
//
// Gorm only updates the attributes that are not null
func (API *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	itemID := mux.Vars(r)["itemID"]

	var (
		count       int
		item        models.Item
		itemChanges models.Item
	)

	API.DB.Model(&item).
		Where("id = ?", itemID).Count(&count)

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Item does not exist"})
	} else {
		err := json.NewDecoder(r.Body).Decode(&itemChanges)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Body should be an item"})
		} else {
			API.DB.Model(&item).Updates(itemChanges)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Item updated correctly"})
		}
	}
}
