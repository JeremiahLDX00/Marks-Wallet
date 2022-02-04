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

type ListTransactions struct {
	TransactionID   int
	StudentID       string
	ToStudentID     string
	TokenTypeName   string
	TransactionType string
	Amount          int
}

var db *sql.DB

//----------------------------- Function for GET All transactions ----------------------------

func GetTransactions(db *sql.DB, StudentID string) []ListTransactions {
	query := `SELECT TransactionID, StudentID, ToStudentID, TokenType.TokenTypeName, TransactionType, Amount 
	FROM Transactions
	INNER JOIN TokenType
	ON Transactions.TokenTypeID = TokenType.TokenTypeID
	WHERE StudentID = ?`
	results, err := db.Query(query, StudentID)
	//handle error
	if err != nil {
		panic(err.Error)
	}

	allTransactions := make([]ListTransactions, 0)
	index := 0
	for results.Next() {
		var allTransaction ListTransactions
		//map this type to the record in the table
		err = results.Scan(&allTransaction.TransactionID, &allTransaction.StudentID, &allTransaction.ToStudentID,
			&allTransaction.TokenTypeName, &allTransaction.TransactionType, &allTransaction.Amount)
		if err != nil {
			panic(err.Error())
		}
		allTransactions = append(allTransactions, allTransaction)
		//fmt.Println(allTransaction.TokenTypeName)
		index++
	}
	return allTransactions
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//GET request to display all transactions referencing StudentID
	var transaction Transactions
	if r.Method == "GET" {
		params := mux.Vars(r)
		reqBody, err := ioutil.ReadAll(r.Body)

		//defer the close till after function has been completed
		defer r.Body.Close()
		if err == nil {
			err := json.Unmarshal(reqBody, &transaction)
			if err != nil {
				println(string(reqBody))
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("Invalid information"))
				return
			}
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(GetTransactions(db, params["studentid"]))
		return
	}
}

//--------------------- Function for CREATING OR UPDATING Transaction (Send/receive token) -----------------------------

func SendReceiveToken(db *sql.DB, transactions Transactions) {

	query := fmt.Sprintf(`INSERT INTO Transactions(StudentID, ToStudentID, TokenTypeID, TransactionType, Amount)
	VALUES('%s','%s','%d','%s','%d')`, transactions.StudentID, transactions.ToStudentID, transactions.TokenTypeID,
		transactions.TransactionType, transactions.Amount)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func MakeTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		//POST New Transaction into Transaction table
		if r.Method == "POST" {
			//read the string sent to the service
			var newTransaction Transactions
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// covert JSON to object
				json.Unmarshal(reqBody, &newTransaction)
				fmt.Println(newTransaction.StudentID, newTransaction.ToStudentID, newTransaction.TokenTypeID,
					newTransaction.TransactionType, newTransaction.Amount)
				SendReceiveToken(db, newTransaction)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Transaction has been created"))
			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply all information required"))
			}
		}
	}
}

func main() {

	//Database code
	_db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/MarksWalletdb") //Connecting to the db
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
	router.HandleFunc("/api/v1/Transactions/viewall/{studentid}", GetAllTransactions).Methods("GET")
	router.HandleFunc("/api/v1/Transactions/maketransaction/{studentid}", MakeTransaction).Methods("POST")

	fmt.Println("Driver microservice API operating on port 9072")
	log.Fatal(http.ListenAndServe(":9072", handlers.CORS(headers, methods, origins)(router)))
}
