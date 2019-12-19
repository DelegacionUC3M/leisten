package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// ItemType enumeration to represent the type of object
type ItemType int

// Different tipes of object that can be loaned
const (
	Laboratorio  ItemType = 0
	Papeleria    ItemType = 1
	Herramientas ItemType = 2
	Electronica  ItemType = 3
)

// Item represents an object.
type Item struct {
	gorm.Model
	Name               string `gorm:"not_null"`
	Amount             int    `gorm:"not_null"`
	Type               string
	State              string
	LoanDays           int
	PenaltyCoefficient float64
	MaxLoans           int
	Description        string
}

// Loan to represent student loans
type Loan struct {
	gorm.Model
	Item        Item `gorm:"foreignkey:ItemID"`
	ItemID      int
	Nia         int       `gorm:"not_null"`
	Amount      int       `gorm:"not_null"`
	LoanDate    time.Time `gorm:"not_null"`
	RefundDate  time.Time
	Finished    bool `gorm:"default:false"`
	Description string
}

// Penalty represents a sanction given to a user
type Penalty struct {
	gorm.Model
	Nia          int
	Loan         Loan `gorm:"foreignkey:LoanID"`
	LoanID       int
	sanctionDate time.Time `gorm:"not_null"`
	PenaltyDate  time.Time
	Finished     bool `gorm:"default:false"`
	Description  string
}
