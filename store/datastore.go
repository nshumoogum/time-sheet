package store

import (
	"time"

	"github.com/time-sheet/service/models"
)

// DataStore provides a datastore.Storer interface used to store, retrieve, remove or update datasets
type DataStore struct {
	Backend Storer
}

//go:generate moq -out mocks/datastore.go -pkg mocks . Storer

// Storer represents basic data access via Get, Remove and Upsert methods.
type Storer interface {
	AddTimesheet(timesheet *models.Timesheet) error
	GetTimesheet(day int, month time.Month, year int) (*models.Timesheet, error)
	UpdateTimesheet(timesheet *models.Timesheet) error
}
