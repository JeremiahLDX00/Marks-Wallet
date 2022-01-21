package main

//Start here
import (
	"database/sql"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ListSearchToken struct { // map this type to the record created in the table
	StudentID string //varchar 9
	CM        int    //int 3
	CSF       int    //int 3
	DP        int    //int 3
	PRG1      int    //int 3
	DB        int    //int 3
	ID        int    //int 3
	OSNF      int    //int 3
	PRG2      int    //int 3
	OOAD      int    //int 3
	WEB       int    //int 3
	PFD       int    //int 3
	SDD       int    //int 3
}

var listSearchToken map[string]ListSearchToken
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

//-------------------------------- Functions for DB ----------------------------------------

func GetTokens(db *sql.DB, StudentID string) ListSearchToken {
	query := fmt.Sprintf("Select * FROM ListSearchToken WHERE StudentID = '%s'", StudentID)
	results, err := db.Query(query)
	//handle error
	if err != nil {
		panic(err.Error)
	}
	var listSearchToken ListSearchToken
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&listSearchToken.StudentID, &listSearchToken.CM,
			&listSearchToken.CSF, &listSearchToken.DP, &listSearchToken.PRG1, &listSearchToken.DB, &listSearchToken.ID,
			&listSearchToken.OSNF, &listSearchToken.PRG2, &listSearchToken.OOAD, &listSearchToken.WEB, &listSearchToken.PFD, &listSearchToken.SDD)
		if err != nil {
			panic(err.Error())
		}
	}
	return listSearchToken
}

//------------------------------- Functions for HTTP ----------------------------------
func GetTokeninfo(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	email := r.URL.Query().Get("studentID")

	var listSearchToken ListSearchToken
	var errMsg string

	listSearchToken, errMsg = Login(db, email)
	if errMsg == "StudentID does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No Student with that ID found"))
	} else {
		json.NewEncoder(w).Encode(listSearchToken)
	}
}

func main() {

	listSearchTokens = map[string]ListSearchToken{}

	//Database code
	_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/MarksWalletDB") //Connecting to the db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	db = _db
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/drivers", GetTokeninfo).Methods("GET")

	fmt.Println("Driver microservice API operating on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
