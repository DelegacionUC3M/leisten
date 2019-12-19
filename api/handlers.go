package api

import (
	"github.com/jinzhu/gorm"
)

// Handler contains a pointer to the database
type Handler struct {
	DB *gorm.DB
}
