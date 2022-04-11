/*
Copyright 2022 IBM All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

/*
 Car structure to to store the world state
*/
type Car struct {
	ManufacturerId    string `json:"manufacturerId"`
	CarId             string `json:"carId"`
	DealerId          string `json:"dealerId"`
	ConsumerId        string `json:"consumerId"`
	CarMake           string `json:"carMake"`
	CarModel          string `json:"carModel"`
	CarColor          string `json:"carColor"`
	Status            string `json:"status"`
	ManufacturingDate string `json:"manufacturingDate"`
	ShippingDate      string `json:"shippingDate"`
	DeliveryDate      string `json:"deliveryDate"`
	SoldOnDate        string `json:"soldOnDate"`
	ManufacturerPrice int    `json:"manufacturerPrice"`
	ShippingPrice     int    `json:"shippingPrice"`
	CustomerPrice     int    `json:"customerPrice"`
}

/* let's declare a global Car array
// that we can then populate
// to simulate a world state
*/
var cars []Car

func GetContract(w http.ResponseWriter) *gateway.Contract {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Fprintf(w, "Failed to create wallet: %s\n", err)
	}

	if !wallet.Exists("CarDemoappUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Fprintf(w, "Failed to populate CarDemoappUser wallet contents: %s\n", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "CarDemoappUser"),
	)
	if err != nil {
		fmt.Fprintf(w, "Failed to connect to gateway: %s\n", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Fprintf(w, "Failed to get network: %s\n", err)
	}

	contract := network.GetContract("cardemo")

	return contract
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("CarDemoappUser", identity)
	if err != nil {
		return err
	}
	return nil
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CarDemo API V1!")
	fmt.Fprintf(w, "Welcome to the Car Manufacturer to an owner LifeCycle Smart Contract!")
}

func returnAllCars(w http.ResponseWriter, r *http.Request) {
	contract := GetContract(w)
	result, err := contract.EvaluateTransaction("QueryAllCars")
	if err != nil {
		fmt.Fprintf(w, "Failed to evaluate transaction: %s\n", err)
	}
	json.NewEncoder(w).Encode(result)
}

func returnSingleCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	contract := GetContract(w)

	// Call QueryCar Function and by supplying CarID paramter
	result, err := contract.EvaluateTransaction("QueryCar", key)
	if err != nil {
		fmt.Fprintf(w, "Failed to evaluate QueryCar transaction: %s\n", err)
	}
	json.NewEncoder(w).Encode(result)
}

func _createNewCar(w http.ResponseWriter, r *http.Request) {
	// get the body of the POST request
	// unmarshal this into a new Car struct
	// append this to our cars array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newCar Car
	json.Unmarshal(reqBody, &newCar)
	// update our global cars array to include
	// our new Car
	cars = append(cars, newCar)
	contract := GetContract(w)
	result, err := contract.SubmitTransaction("createNewCar", newCar.ManufacturerId, newCar.CarId, newCar.CarMake, newCar.CarModel, newCar.CarColor, newCar.ManufacturingDate, string(newCar.ManufacturerPrice), "manufacturer")
	if err != nil {
		fmt.Fprintf(w, "Failed to submit  createNewCar transaction: %s\n", err)
	}
	fmt.Fprintf(w, string(result))
}

func _shipToDealer(w http.ResponseWriter, r *http.Request) {
	// get the body of the POST request
	// unmarshal this into a new Car struct
	// append this to our cars array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newCar Car
	json.Unmarshal(reqBody, &newCar)
	// update our global cars array to include
	// our new Car
	cars = append(cars, newCar)
	contract := GetContract(w)
	result, err := contract.SubmitTransaction("ShipToDealer", newCar.CarId, newCar.DealerId, string(newCar.ShippingPrice), "manufacturer")
	if err != nil {
		fmt.Fprintf(w, "Failed to submit  createNewCar transaction: %s\n", err)
	}
	fmt.Fprintf(w, string(result))
}

func _receiveDelivery(w http.ResponseWriter, r *http.Request) {
	// get the body of the POST request
	// unmarshal this into a new Car struct
	// append this to our cars array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newCar Car
	json.Unmarshal(reqBody, &newCar)
	// update our global cars array to include
	// our new Car
	cars = append(cars, newCar)
	contract := GetContract(w)

	// Call ReceiveDelivery Function and supply paramters like carId string, role string
	result, err := contract.SubmitTransaction("ReceiveDelivery", newCar.CarId, "dealer")
	if err != nil {
		fmt.Fprintf(w, "Failed to submit  ReceiveDelivery transaction: %s\n", err)
	}
	fmt.Fprintf(w, string(result))
}

func _sellToCustomer(w http.ResponseWriter, r *http.Request) {
	// get the body of the POST request
	// unmarshal this into a new Car struct
	// append this to our cars array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newCar Car
	json.Unmarshal(reqBody, &newCar)
	// update our global cars array to include
	// our new Car
	cars = append(cars, newCar)
	contract := GetContract(w)

	// Call SellToCustomer Function and supply paramters like carId string, consumerId string, customerPrice int, role string
	result, err := contract.SubmitTransaction("SellToCustomer", newCar.CarId, newCar.ConsumerId, string(newCar.CustomerPrice), "dealer")
	if err != nil {
		fmt.Fprintf(w, "Failed to submit  SellToCustomer transaction: %s\n", err)
	}
	fmt.Fprintf(w, string(result))
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", welcome)
	myRouter.HandleFunc("/getCars", returnAllCars)
	myRouter.HandleFunc("/getCar/{id}", returnSingleCar)
	myRouter.HandleFunc("/create", _createNewCar).Methods("POST")
	myRouter.HandleFunc("/ship", _shipToDealer).Methods("POST")
	myRouter.HandleFunc("/receive", _receiveDelivery).Methods("POST")
	myRouter.HandleFunc("/sell", _sellToCustomer).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
