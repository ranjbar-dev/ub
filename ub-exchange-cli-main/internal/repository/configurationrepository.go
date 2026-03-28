package repository

import (
	"exchange-go/internal/configuration"
	"gorm.io/gorm"
)

type configurationRepository struct {
	db *gorm.DB
}

func NewConfigurationRepository(db *gorm.DB) configuration.Repository {
	return &configurationRepository{db}
}