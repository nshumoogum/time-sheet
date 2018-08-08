package store

// DataStore provides a datastore.Storer interface used to store, retrieve, remove or update datasets
type DataStore struct {
	Backend Storer
}

//go:generate moq -out mocks/datastore.go -pkg mocks . Storer

// Storer represents basic data access via Get, Remove and Upsert methods.
type Storer interface {
}
