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

type Driver struct { // map this type to the record created in the table
	DriverID  int    //varchar 5
	FirstName string //varchar 30
	LastName  string // varchar 30
	MobileNo  string //varchar 8
	Email     string //varchar 40
	LicenseNo string //varchar 5
}

var drivers map[string]Driver
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
func GetDriver(db *sql.DB, Email string) Driver {
	query := fmt.Sprintf("Select * FROM Driver WHERE Email = '%s'", Email)
	results, err := db.Query(query)
	//handle error
	if err != nil {
		panic(err.Error)
	}
	var driver Driver
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&driver.DriverID, &driver.FirstName,
			&driver.LastName, &driver.MobileNo, &driver.Email, &driver.LicenseNo)
		if err != nil {
			panic(err.Error())
		}
	}
	return driver
}

func CreateDriver(DriverID string, FirstName string, LastName string, MobileNo string, Email string, LicenseNo string) {

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Driver_db")

	query := fmt.Sprintf("INSERT INTO Driver(DriverID, FirstName, LastName, MobileNo, Email, LicenseNo) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')",
		DriverID, FirstName, LastName, MobileNo, Email, LicenseNo)

	_, err = db.Query(query)

	if err != nil {

		panic(err.Error())
	}

	defer db.Close()
}

func UpdateDriver(db *sql.DB, DriverID int, d Driver) {
	query := fmt.Sprintf("UPDATE Driver SET FirstName = '%s', LastName = '%s', MobileNo = '%s', Email = '%s', LicenseNo = '%s' WHERE DriverID = '%d'",
		d.FirstName, d.LastName, d.MobileNo, d.Email, d.LicenseNo, d.DriverID)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

func DeleteDriver(db *sql.DB, DriverID int) string { //not to be implemented into the FE
	query := fmt.Sprintf("DELETE FROM Driver WHERE ID='%d'", DriverID)

	_, err := db.Query(query)
	var errMsg string

	if err != nil {
		errMsg = "Account does not exist"
	}
	return errMsg
}

func Login(db *sql.DB, email string) (Driver, string) {
	query := fmt.Sprintf("SELECT * FROM Driver where Email = '%s'", email)

	results := db.QueryRow(query)

	var driver Driver
	var errMsg string

	switch err := results.Scan(&driver.DriverID, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.Email, &driver.LicenseNo); err {
	case sql.ErrNoRows:
		errMsg = "Account does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return driver, errMsg
}

//----------------------------------- functions for http --------------------------------------------

func CreateDriverAccount(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {

			// read the string sent to the service
			var driver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &driver)

				if driver.FirstName == "" || driver.LastName == "" || driver.MobileNo == "" || driver.Email == "" || driver.LicenseNo == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply Driver " +
							"information " + "in JSON format" + " & Provide all required Fields"))
					return
				}

				// check if course exists; add only if
				// course does not exist
				if _, ok := drivers[params["driverid"]]; !ok {
					drivers[params["driverid"]] = driver
					CreateDriver(params["driverid"], driver.FirstName, driver.LastName, driver.MobileNo, driver.Email, driver.LicenseNo)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - This Driver has made an account: " +
						params["driverid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Duplicate Driver ID"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply Driver information " +
					"in JSON format"))
			}
		}
	}
}

func GetDriverInfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	email := r.URL.Query().Get("email")

	var driver Driver
	var errMsg string

	driver, errMsg = Login(db, email)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		json.NewEncoder(w).Encode(driver)
	}
}

func UpdateDriverInfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)
	var driverid int
	fmt.Sscan(params["driverid"], &driverid)

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		var driver Driver
		json.Unmarshal([]byte(reqBody), &driver)

		if driver.FirstName == "" || driver.LastName == "" || driver.MobileNo == "" || driver.Email == "" || driver.LicenseNo == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please include all required information "))
		} else {
			UpdateDriver(db, driverid, driver)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Account details has been successfully updated"))
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply information in JSON format"))
	}
}

func DeleteDriverAccount(w http.ResponseWriter, r *http.Request) { // not to be implemented into FE
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	params := mux.Vars(r)
	var driverid int
	fmt.Sscan(params["driverid"], &driverid)

	var driver Driver

	errMsg := DeleteDriver(db, driverid)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Driver Account Deleted: " + driver.Email))
	}
}

func main() {

	drivers = map[string]Driver{}

	//Database code
	_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Driver_db") //Connecting to the db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	db = _db
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/drivers/{driverid}", CreateDriverAccount).Methods("POST")
	router.HandleFunc("/drivers", GetDriverInfo).Methods("GET")
	router.HandleFunc("/drivers/{driverid}", UpdateDriverInfo).Methods("PUT")
	router.HandleFunc("/drivers/{driverid}", DeleteDriverAccount).Methods("DELETE")

	fmt.Println("Driver microservice API operating on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
