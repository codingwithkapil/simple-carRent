package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// This function initializes the router and returns it for testing purposes.
func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars", AddCar).Methods("POST")
	router.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	router.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")
	return router
}

func TestGetCars(t *testing.T) {
	cars = []Car{
		{Model: "Tesla M3", Registration: "BTS812", Mileage: 6003, Condition: "available"},
	}

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cars", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("GET /cars returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"model":"Tesla M3","registration":"BTS812","mileage":6003,"condition":"available"}]`
	got := rec.Body.String()
	if strings.TrimSpace(got) != expected {
		t.Errorf("GET /cars returned unexpected body: got %v want %v", got, expected)
	}
}

func TestAddCar(t *testing.T) {
	newCar := []byte(`{"model":"BMW X5","registration":"XYZ123","mileage":5000,"condition":"available"}`)
	req, _ := http.NewRequest("POST", "/cars", bytes.NewBuffer(newCar))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("POST /cars returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseCar map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &responseCar)

	if responseCar["model"] != "BMW X5" {
		t.Errorf("POST /cars returned unexpected model: got %v want %v", responseCar["model"], "BMW X5")
	}
}

func TestRentCar(t *testing.T) {
	// Initialize test data
	cars = []Car{
		{Model: "Tesla M3", Registration: "BTS812", Mileage: 6003, Condition: "available"},
		{Model: "BMW X5", Registration: "XYZ123", Mileage: 5000, Condition: "available"},
	}

	req, err := http.NewRequest("POST", "/cars/XYZ123/rentals", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	router.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("POST /cars/{registration}/rentals returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Add assertions for expected response data and behavior
}

func TestReturnCar(t *testing.T) {
	// Initialize test data
	cars = []Car{
		{Model: "Tesla M3", Registration: "BTS812", Mileage: 6003, Condition: "rented"},
		{Model: "BMW X5", Registration: "XYZ123", Mileage: 5000, Condition: "available"},
	}

	req, err := http.NewRequest("POST", "/cars/BTS812/returns", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")
	router.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("POST /cars/{registration}/returns returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
