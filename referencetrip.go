package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//----------------------------------------------- Structure -----------------------------------------------
type Trip struct { // map this type to the record created in the table
	TripID      int    //varchar 5
	PickUpPC    string //varchar 6
	DropOffPC   string //varchar 6
	PassengerID string //varchar 5
	DriverID    string //varchar 5
	TripStatus  string //varchar 20
}

var db *sql.DB

func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//------------------------------------------ functions for DB ------------------------------------------
func GetPassenger(db *sql.DB, passengerID string) Trip {
	query := fmt.Sprintf("Select * FROM Passenger WHERE PassengerID = '%s'", passengerID)
	results, err := db.Query(query)
	//handle error
	if err != nil {
		panic(err.Error)
	}
	var trip Trip
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&trip.TripID, &trip.PickUpPC,
			&trip.DropOffPC, &trip.PassengerID, &trip.DriverID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		}
	}
	return trip
}

func GetTrip(db *sql.DB, TripID string) (Trip, string) {
	query := fmt.Sprintf("Select * FROM Trip WHERE TripID = '%s'", TripID)
	results, err := db.Query(query)
	//handle error
	if err != nil {
		panic(err.Error)
	}
	var trip Trip
	var errMsg string

	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&trip.TripID, &trip.PickUpPC,
			&trip.DropOffPC, &trip.PassengerID, &trip.DriverID, &trip.TripStatus)
		if err != nil {
			panic(err.Error())
		}
	}
	return trip, errMsg
}

func CreateTrip(db *sql.DB, t Trip) {
	query := fmt.Sprintf("INSERT INTO Trip (PickUpPC, DropOffPC, PassengerID, DriverID, TripStatus) VALUES ('%s','%s','%s','%s','%s')",
		t.PickUpPC, t.DropOffPC, t.PassengerID, t.DriverID, t.TripStatus)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func UpdateTrip(db *sql.DB, TripID int, t Trip) {
	query := fmt.Sprintf("UPDATE Trip SET DriverID = '%s', TripStatus = '%s' WHERE TripID = '%d'",
		t.DriverID, t.TripStatus, t.TripID)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//----------------------------------- functions for http --------------------------------------------

func CreateTripDetail(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		var trip Trip
		json.Unmarshal([]byte(reqBody), &trip)

		if trip.PickUpPC == "" || trip.DropOffPC == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please include all required information "))
		}

	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply information in JSON format"))
	}
}

func GetTripInfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)

	var trip Trip
	var errMsg string

	trip, errMsg = GetTrip(db, params["tripid"])
	if errMsg == "Trip does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		json.NewEncoder(w).Encode(trip)
	}
}

func UpdateTripInfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)
	var tripid int
	fmt.Sscan(params["tripid"], &tripid)

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		var trip Trip
		json.Unmarshal([]byte(reqBody), &trip)

		if trip.PickUpPC == "" || trip.DropOffPC == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please include all required information "))
		} else {
			UpdateTrip(db, tripid, trip)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Trip details has been successfully updated"))
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply information in JSON format"))
	}
}

func GetPassengerInfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	//id := r.URL.Query().Get("passengerID")

}

func main() {
	//Database code
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Trip_db") //Connecting to the db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/trip", CreateTripDetail).Methods("POST")
	router.HandleFunc("/trip", GetTripInfo).Methods("GET")
	router.HandleFunc("/trip/{tripid}", UpdateTripInfo).Methods("PUT")

	fmt.Println("Trip microservice API operating on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
