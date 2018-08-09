package service

import (
	"context"
	"fmt"

	"github.com/ONSdigital/go-ns/server"
	"github.com/gorilla/mux"
	"github.com/time-sheet/config"
	"github.com/time-sheet/service/handlers"
	"github.com/time-sheet/store"
)

var (
	httpServer   *server.Server
	serverErrors chan error
)

// TimeSheetAPI manages stored data
type TimeSheetAPI struct {
	DataStore store.DataStore
	Router    *mux.Router
}

// CreateTimeSheetAPI manages all the routes configured to API
func CreateTimeSheetAPI(cfg config.Configuration, dataStore store.DataStore, errorChan chan error) {
	router := mux.NewRouter()
	Routes(cfg, router, dataStore)

	httpServer = server.New(cfg.BindAddr, router)

	// Disable this here to allow main to manage graceful shutdown of the entire app.
	httpServer.HandleOSSignals = false

	go func() {
		fmt.Printf("Starting api...")
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Printf("api http server returned error: [%v]\n", err)
			errorChan <- err
		}
	}()
}

// Routes represents a list of endpoints that exist with this api
func Routes(cfg config.Configuration, router *mux.Router, dataStore store.DataStore) *TimeSheetAPI {

	api := TimeSheetAPI{
		DataStore: dataStore,
		Router:    router,
	}

	timeSheetStore := handlers.Store{Storer: api.DataStore.Backend}
	// TODO api.Router.HandleFunc("/assignments", GetAssignments).Methods("GET")
	// TODO api.Router.HandleFunc("/assignments/{assignment}", GetAssignment).Methods("GET")
	api.Router.HandleFunc("/clock-in", timeSheetStore.PostStart).Methods("POST")
	// TODO api.Router.HandleFunc("/clock-out", PostEnd).Methods("POST")
	// TODO api.Router.HandleFunc("/years", GetYears).Methods("GET")
	// TODO api.Router.HandleFunc("/years/{year}", GetYear).Methods("GET")
	// TODO api.Router.HandleFunc("/years/{year}/months", GetMonths).Methods("GET")
	// TODO api.Router.HandleFunc("/years/{year}/months/{month}", GetMonth).Methods("GET")
	// TODO api.Router.HandleFunc("/years/{year}/months/{month}/days", GetDays).Methods("GET")
	// TODO api.Router.HandleFunc("/year/{year}/month/{month}/day/{day}", GetDay).Methods("Get")
	// TODO Edit time for date (currently not in spec) api.Router.HandleFunc("/year/{year}/month/{month}/day/{day}", PutDay).Methods("PUT")
	// TODO api.Router.HandleFunc("/year/{year}/month/{month}/day/{day}/assignments/{assignment}", PutDay).Methods("POST")

	return &api
}

// Close represents the graceful shutting down of the http server
func Close(ctx context.Context) error {
	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}

	fmt.Println("graceful shutdown of http server complete")
	return nil
}
