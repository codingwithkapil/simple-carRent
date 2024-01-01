package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Car struct {
	Model        string `json:"model"`
	Registration string `json:"registration"`
	Mileage      int    `json:"mileage"`
	Condition    string `json:"condition"`
}

var cars []Car

func GetCars(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cars)
}

func AddCar(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	_ = json.NewDecoder(r.Body).Decode(&newCar)
	for _, car := range cars {
		if car.Registration == newCar.Registration {
			http.Error(w, "Car already exists", http.StatusBadRequest)
			return
		}
	}
	cars = append(cars, newCar)
	json.NewEncoder(w).Encode(newCar)
}

func RentCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]
	for i, car := range cars {
		if car.Registration == registration {
			if car.Condition == "rented" {
				http.Error(w, "Car already rented", http.StatusBadRequest)
				return
			}
			cars[i].Condition = "rented"
			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}
	http.Error(w, "Car not found", http.StatusNotFound)
}

func ReturnCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]
	for i, car := range cars {
		if car.Registration == registration {
			if car.Condition != "rented" {
				http.Error(w, "Car is not rented", http.StatusBadRequest)
				return
			}
			// Here you would handle updating mileage based on the kilometers driven
			cars[i].Condition = "available"
			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}
	http.Error(w, "Car not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	// Sample data to start with
	cars = append(cars, Car{Model: "Tesla M3", Registration: "BTS812", Mileage: 6003, Condition: "available"})

	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars", AddCar).Methods("POST")
	router.HandleFunc("/cars/{registration}/rentals", RentCar).Methods("POST")
	router.HandleFunc("/cars/{registration}/returns", ReturnCar).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
