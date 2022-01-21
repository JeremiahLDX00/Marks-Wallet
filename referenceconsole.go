package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//----------------------------structs------------------------

type Passenger struct { // map this type to the record created in the table
	PassengerID int    //varchar 5
	FirstName   string //varchar 30
	LastName    string //varchar 30
	MobileNo    string //varchar 8
	Email       string //varchar 40
}

type Driver struct { // map this type to the record created in the table
	DriverID  int    //varchar 5
	FirstName string //varchar 30
	LastName  string // varchar 30
	MobileNo  string //varchar 8
	Email     string //varchar 40
	LicenseNo string //varchar 5
}

type Trip struct { // map this type to the record created in the table
	TripID      int    //varchar 5
	PickUpPC    string //varchar 6
	DropOffPC   string //varchar 6
	PassengerID string //varchar 5
	DriverID    string //varchar 5
	TripStatus  string //varchar 20
}

//connections
const DriverbaseURL = "http://localhost:5000/drivers"
const PassengerbaseURL = "http://localhost:5000/passengers"
const TripbaseURL = "http://localhost:5000/trip"
const key = "2c78afaf"

//---------------------------Menus---------------------------
func HomeMenu() {

	loop := 1
	userinput := ""

	for loop != 0 {
		fmt.Println("============= Welcome to Ride Sharing! =============")
		fmt.Println("[1] Register as a Passenger")
		fmt.Println("[2] Register as a Driver")
		fmt.Println("[0] Exit")
		fmt.Println("Please choose an option: ")
		fmt.Scanln(&userinput)

		if userinput == "0" {
			fmt.Println("Thank you for using ride share!")
			loop = 0
		} else if userinput == "1" {
			RegisterPassengerMenu()
		} else if userinput == "2" {
			RegisterDriverMenu()
		} else {
			fmt.Println("Please provide a valid input")
		}
	}
}

func RegisterPassengerMenu() {
	id := 1
	var FirstName string
	var LastName string
	var MobileNo string
	var Email string

	fmt.Println("============= Register as a Passenger =============")
	fmt.Println("Please enter your first name: ")
	fmt.Scanln(&FirstName)
	fmt.Println("Please enter your last name: ")
	fmt.Scanln(&LastName)
	fmt.Println("Please enter your mobile number: ")
	fmt.Scanln(&MobileNo)
	fmt.Println("Please enter your email address: ")
	fmt.Scanln(&Email)

	jsonData := map[string]string{"FirstName": "", "LastName": "", "MobileNo": "", "Email": ""}
	jsonData["FirstName"] = FirstName
	jsonData["LastName"] = LastName
	jsonData["MobileNo"] = MobileNo
	jsonData["Email"] = Email

	CreatePassengerAccount(fmt.Sprintf("%v", id), jsonData)
	id += 1

	MainPassengerMenu()
}

func MainPassengerMenu() {

	userinput := ""
	loop := 1

	for loop != 0 {
		fmt.Println("============= Passenger =============")
		fmt.Println("[1] Update account details")
		fmt.Println("[2] Create a Trip")
		fmt.Println("[3] Display all Trips")
		fmt.Println("[0] Exit")
		fmt.Println("Please choose an option: ")
		fmt.Scanln(&userinput)
		if userinput == "1" {
			UpdatePassengerMenu()
		} else if userinput == "2" {
			CreateTrip()
		} else if userinput == "3" {
			DisplayTrip()
		} else if userinput == "0" {
			fmt.Println("Thank you for using ride share!")
			loop = 0
		} else {
			fmt.Println("Invalid input, please try again")
		}

	}
}

func UpdatePassengerMenu() {

	id := ""
	var FirstName string
	var LastName string
	var MobileNo string
	var Email string

	fmt.Println("============= Update Passenger Details =============")
	fmt.Println("Please enter your ID:")
	fmt.Scanln(&id)
	fmt.Println("Please enter your first name: ")
	fmt.Scanln(&FirstName)
	fmt.Println("Please enter your last name: ")
	fmt.Scanln(&LastName)
	fmt.Println("Please enter your mobile number: ")
	fmt.Scanln(&MobileNo)
	fmt.Println("Please enter your email address: ")
	fmt.Scanln(&Email)

	jsonData := map[string]string{"PassengerID": "", "FirstName": "", "LastName": "", "MobileNo": "", "Email": ""}
	jsonData["PassengerID"] = id
	jsonData["FirstName"] = FirstName
	jsonData["LastName"] = LastName
	jsonData["MobileNo"] = MobileNo
	jsonData["Email"] = Email

	UpdatePassenger(fmt.Sprintf("%v", id), jsonData)

	MainPassengerMenu()
}

func RegisterDriverMenu() {
	id := 1
	var FirstName string
	var LastName string
	var MobileNo string
	var Email string
	var LicenseNo string

	fmt.Println("============= Register as a Driver =============")
	fmt.Println("Please enter your first name: ")
	fmt.Scanln(&FirstName)
	fmt.Println("Please enter your last name: ")
	fmt.Scanln(&LastName)
	fmt.Println("Please enter your mobile number: ")
	fmt.Scanln(&MobileNo)
	fmt.Println("Please enter your email address: ")
	fmt.Scanln(&Email)
	fmt.Println("Please enter your license number: ")
	fmt.Scanln(&LicenseNo)

	jsonData := map[string]string{"FirstName": "", "LastName": "", "MobileNo": "", "Email": "", "LicenseNo": ""}
	jsonData["FirstName"] = FirstName
	jsonData["LastName"] = LastName
	jsonData["MobileNo"] = MobileNo
	jsonData["Email"] = Email
	jsonData["LicenseNo"] = LicenseNo

	CreateDriverAccount(fmt.Sprintf("%v", id), jsonData)
	id += 1

	MainDriverMenu()
}

func UpdateDriverMenu() {
	id := ""
	var FirstName string
	var LastName string
	var MobileNo string
	var Email string
	var LicenseNo string

	fmt.Println("============= Update Driver Details =============")
	fmt.Println("Please enter your ID:")
	fmt.Scanln(&id)
	fmt.Println("Please enter your first name: ")
	fmt.Scanln(&FirstName)
	fmt.Println("Please enter your last name: ")
	fmt.Scanln(&LastName)
	fmt.Println("Please enter your mobile number: ")
	fmt.Scanln(&MobileNo)
	fmt.Println("Please enter your email address: ")
	fmt.Scanln(&Email)
	fmt.Println("Please enter your License number: ")
	fmt.Scanln(&LicenseNo)

	jsonData := map[string]string{"DriverID": "", "FirstName": "", "LastName": "", "MobileNo": "", "Email": "", "LicenseNo": ""}
	jsonData["DriverID"] = id
	jsonData["FirstName"] = FirstName
	jsonData["LastName"] = LastName
	jsonData["MobileNo"] = MobileNo
	jsonData["Email"] = Email
	jsonData["LicenseNo"] = LicenseNo

	UpdateDriver(fmt.Sprintf("%v", id), jsonData)

	MainPassengerMenu()
}

func MainDriverMenu() {
	userinput := ""
	loop := 1
	for loop != 0 {
		fmt.Println("============= Driver =============")
		fmt.Println("[1] Update Driver details")
		fmt.Println("[2] Start a Trip")
		fmt.Println("[3] End a Trip")
		fmt.Println("[0] Exit")
		fmt.Println("Please choose an option: ")
		fmt.Scanln()
		if userinput == "1" {
			UpdateDriverMenu()
		} else if userinput == "2" {
			StartTrip()
		} else if userinput == "3" {
			EndTrip()
		} else if userinput == "0" {
			fmt.Println("Thank you for using ride share!")
			loop = 0
		} else {
			fmt.Println("Invalid input, please try again")
		}
	}
	if userinput == "1" {
		UpdateDriverMenu()
	} else if userinput == "2" {
		StartTrip()
	} else if userinput == "3" {
		EndTrip()
	} else if userinput == "0" {
		fmt.Println("Thank you for using ride share!")
	} else {
		fmt.Println("Invalid input, please try again")
	}
}

//-------------------------------------- functions -------------------------------

func getPassengerAccount(code string) {
	url := PassengerbaseURL
	if code != "" {
		url = PassengerbaseURL + "/" + code + "?key=" + key
	}
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func CreatePassengerAccount(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(PassengerbaseURL+"/"+code+"?key="+key,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func CreateDriverAccount(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(DriverbaseURL+"/"+code+"?key="+key,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func UpdatePassenger(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(PassengerbaseURL+"/"+code+"?key="+key,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func UpdateDriver(code string, jsonData map[string]string) {
	jsonValue, _ := json.Marshal(jsonData)

	response, err := http.Post(DriverbaseURL+"/"+code+"?key="+key,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func DisplayTrip() {
	response, err := http.Get("localhost5000/trip")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
}

func CreateTrip() {

}

func StartTrip() {

}

func EndTrip() {

}

func main() {
	HomeMenu()
}
