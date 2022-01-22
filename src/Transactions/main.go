package main

//Start here
import (
	_ "github.com/go-sql-driver/mysql"
)


func main (

	router.HandleFunc("/api/v1/Transactions/", GetAllTransactions).Methods("GET")
	router.HandleFunc("/api/v1/Transactions/{studentid}", SendReceiveTokens).Methods("PUT")
	fmt.Println("Driver microservice API operating on port 9072")
	log.Fatal(http.ListenAndServe(":9072", router))
)