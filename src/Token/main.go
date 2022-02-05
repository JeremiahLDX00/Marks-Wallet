package main

//Start here
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type TokenType struct { // map this type to the record created in the table
	TokenTypeID   int    //int 5
	TokenTypeName string //varchar 5
}

type Transactions struct {
	TransactionID   int    //int 3
	StudentID       string //varchar 9
	ToStudentID     string //varchar 9
	TokenTypeID     int    //int 5
	TransactionType string //varchar 30
	Amount          int    //int 3,
}

type TokenTypeBalance struct {
	TokenTypeID   int
	TokenTypeName string
	Balance       int
}

var db *sql.DB

//-------------------------------- Functions for GET Token ----------------------------------------

func GetStudentTokens(db *sql.DB, StudentID string) []TokenTypeBalance {
	query := `SELECT TokenType.TokenTypeID, TokenType.TokenTypeName, SUM(Amount) FROM Transactions 
	INNER JOIN TokenType 
	ON Transactions.TokenTypeID = TokenType.TokenTypeID 
	WHERE StudentID = ?
	GROUP BY StudentID, Transactions.TokenTypeID`
	results, err := db.Query(query, StudentID)
	//handle error
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error)
	}

	tokenTypeBalances := make([]TokenTypeBalance, 0)
	index := 0
	for results.Next() {
		var tokenTypeBalance TokenTypeBalance
		// map this type to the record in the table
		err = results.Scan(&tokenTypeBalance.TokenTypeID, &tokenTypeBalance.TokenTypeName, &tokenTypeBalance.Balance)
		if err != nil {
			panic(err.Error())
		}
		tokenTypeBalances = append(tokenTypeBalances, tokenTypeBalance)
		// fmt.Println(TokenTypeBalance.TokenTypeName)
		index++
	}
	return tokenTypeBalances
}

func GetAvailableTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//GET request to display tokens referencing StudentID
	var token TokenType
	if r.Method == "GET" {
		params := mux.Vars(r)
		reqBody, err := ioutil.ReadAll(r.Body)

		//defer the close till after function has finished running
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &token)
			if err != nil {
				println(string(reqBody))
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information"))
				return
			}
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(GetStudentTokens(db, params["studentid"]))
		return
	}
}

//--------------------------------- Functions for Search Token ---------------------------------------

func SearchTokens(db *sql.DB, TokenTypeName string) TokenType {
	query := `SELECT TokenTypeID, TokenTypeName FROM TokenType 
	WHERE TokenTypeName = ?`
	print("TokenTypeName: ", TokenTypeName)
	results, err := db.Query(query, TokenTypeName)
	//handle error
	if err != nil {
		panic(err.Error)
	}
	var searchToken TokenType
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&searchToken.TokenTypeID, &searchToken.TokenTypeName)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println(searchToken.TokenTypeID, searchToken.TokenTypeName)
	}
	return searchToken
}

func SearchForTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//GET request to display tokens referencing Token Name
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(SearchTokens(db, params["tokentypename"]))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No such token found"))
	}
}

func main() {

	//Database code
	_db, err := sql.Open("mysql", "root:password@tcp(markswalletdb:3306)/markswalletdb") //Connecting to the db
	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	db = _db
	defer db.Close()

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/Token/student/{studentid}", GetAvailableTokens).Methods("GET")
	router.HandleFunc("/api/v1/Token/search/{tokentypename}", SearchForTokens).Methods("GET")

	fmt.Println("Driver microservice API operating on port 9071")
	log.Fatal(http.ListenAndServe(":9071", handlers.CORS(headers, methods, origins)(router)))
}
