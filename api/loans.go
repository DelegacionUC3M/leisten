package api

import (
	"encoding/json"
	// "fmt"
	// "github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"
	"net/http"

	models "github.com/DelegacionUC3M/leisten/models"
	// "strconv"
)

// GetAllLoans returns all loans in the database
func (API *Handler) GetAllLoans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		loansList []models.Loan
		loan      models.Loan
	)

	loansRows, err := API.DB.Model(models.Loan{}).Rows()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Could not fetch loans"})
		return
	}

	for loansRows.Next() {
		API.DB.ScanRows(loansRows, &loan)
		loansList = append(loansList, loan)
	}

	w.WriteHeader(http.StatusOK)
	payload, _ := json.Marshal(loansList)
	w.Write(payload)
}

// CreateLoan creates a new loan and inserts it into the database
func (API *Handler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data models.Loan

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Could not create loan"})
	} else {
		API.DB.Create(&data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]uint{"id": data.ID})
	}
}
