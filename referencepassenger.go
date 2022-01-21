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

type Passenger struct { // map this type to the record created in the table
	PassengerID int    //varchar 5
	FirstName   string //varchar 30
	LastName    string //varchar 30
	MobileNo    string //varchar 8
	Email       string //varchar 40
}

var passengers map[string]Passenger
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

func CreatePassenger(PassengerID string, FirstName string, LastName string, MobileNo string, Email string) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Passenger_db")

	query := fmt.Sprintf("INSERT INTO Passenger(PassengerID, FirstName, LastName, MobileNo, Email) VALUES ('%s', '%s', '%s', '%s', '%s')",
		PassengerID, FirstName, LastName, MobileNo, Email)

	_, err = db.Query(query)

	if err != nil {

		panic(err.Error())
	}

	defer db.Close()
}

func UpdatePassenger(PassengerID string, FirstName string, LastName string, MobileNo string, Email string) {
	query := fmt.Sprintf("UPDATE Passenger SET FirstName = '%s', LastName = '%s', MobileNo = '%s', Email = '%s' WHERE PassengerID = '%s'",
		FirstName, LastName, MobileNo, Email, PassengerID)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func DeletePassenger(db *sql.DB, passengerID int) string { //not to be implemented into the FE
	query := fmt.Sprintf("DELETE FROM Passenger WHERE ID='%d'", passengerID)

	_, err := db.Query(query)
	var errMsg string

	if err != nil {
		errMsg = "Account does not exist"
	}
	return errMsg
}

//----------------------------------- functions for http --------------------------------------------

func CreatePassengerAccount(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {

			// read the string sent to the service
			var passenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &passenger)

				if passenger.FirstName == "" || passenger.LastName == "" || passenger.MobileNo == "" || passenger.Email == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply Passenger " +
							"information " + "in JSON format" + " & Provide all required Fields"))
					return
				}

				// check if course exists; add only if
				// course does not exist
				if _, ok := passengers[params["passengerid"]]; !ok {
					passengers[params["passengerid"]] = passenger
					CreatePassenger(params["passengerid"], passenger.FirstName, passenger.LastName, passenger.MobileNo, passenger.Email)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - This Passenger has made an account: " +
						params["passengerid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Duplicate Passenger ID"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " +
					"in JSON format"))
			}
		}
	}
}

func UpdatePassengerInfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "PUT" {

			// read the string sent to the service
			var passenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &passenger)

				if passenger.FirstName == "" || passenger.LastName == "" || passenger.MobileNo == "" || passenger.Email == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply Passenger " +
							"information " + "in JSON format" + " & Provide all required Fields"))
					return
				}

				// check if course exists; add only if
				// course does not exist
				if _, ok := passengers[params["passengerid"]]; !ok {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(
						"404 - No such Passenger ID"))
				} else {
					passengers[params["passengerid"]] = passenger
					UpdatePassenger(params["passengerid"], passenger.FirstName, passenger.LastName, passenger.MobileNo, passenger.Email)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - This Passenger has updated their account!: " +
						params["passengerid"]))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply passenger information " +
					"in JSON format"))
			}
		}
	}
}

func DeletePassengerAccount(w http.ResponseWriter, r *http.Request) { // not to be implemented into FE
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)
	var passengerid int
	fmt.Sscan(params["passengerid"], &passengerid)

	var passenger Passenger

	errMsg := DeletePassenger(db, passengerid)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Passenger Account Deleted: " + passenger.Email))
	}
}

func main() {

	// instantiate errors
	passengers = make(map[string]Passenger)

	//Database code
	_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Passenger_db") //Connecting to the db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	db = _db
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/passengers/{passengerid}", CreatePassengerAccount).Methods("POST", "GET", "PUT", "DELETE")

	fmt.Println("Passenger microservice API operating on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
