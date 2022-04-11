/*
SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Smart Contract Chaincode provides functions to manage a Car manufacturer to an owner delivary
type CarChainCode struct {
	contractapi.Contract
}

/*
 Car structure to record the world state
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

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *CarChainCode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{ManufacturerId: "MOrg01", CarId: "M101", DealerId: "D101", ConsumerId: "CUST101", CarMake: "2022", CarModel: "MOrg01CM101", CarColor: "Red", Status: "SOLD", ManufacturingDate: "2022/01/01", ShippingDate: "2022/02/01", DeliveryDate: "2022/02/20", SoldOnDate: "2022/04/20", ManufacturerPrice: 350000, ShippingPrice: 10000, CustomerPrice: 550000},
		Car{ManufacturerId: "MOrg01", CarId: "M102", DealerId: "D102", ConsumerId: "CUST102", CarMake: "2022", CarModel: "MOrg01CM102", CarColor: "Blue", Status: "SOLD", ManufacturingDate: "2022/01/01", ShippingDate: "2022/02/01", DeliveryDate: "2022/02/20", SoldOnDate: "2022/04/20", ManufacturerPrice: 360000, ShippingPrice: 10000, CustomerPrice: 600000},
		Car{ManufacturerId: "MOrg02", CarId: "M103", DealerId: "D102", ConsumerId: "CUST103", CarMake: "2022", CarModel: "MOrg02CM103", CarColor: "Blue", Status: "SOLD", ManufacturingDate: "2022/01/01", ShippingDate: "2022/02/01", DeliveryDate: "2022/02/20", SoldOnDate: "2022/04/20", ManufacturerPrice: 360000, ShippingPrice: 10000, CustomerPrice: 630000},
		Car{ManufacturerId: "MOrg02", CarId: "M104", DealerId: "D101", ConsumerId: "CUST101", CarMake: "2022", CarModel: "MOrg01CM101", CarColor: "Red", Status: "SOLD", ManufacturingDate: "2022/01/01", ShippingDate: "2022/02/01", DeliveryDate: "2022/02/20", SoldOnDate: "2022/04/20", ManufacturerPrice: 350000, ShippingPrice: 10000, CustomerPrice: 550000},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed Car data to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *CarChainCode) createNewCar(ctx contractapi.TransactionContextInterface, manufacturerId string, carId string, carMake string, carModel string, carColor string, manufacturingDate string, manufacturerPrice int, role string) error {

	if role != "manufacturer" {
		return fmt.Errorf("Failed to put new Car to world state due to unauthorized user")
	}
	car := Car{
		ManufacturerId: manufacturerId,
		CarId:          carId,
		CarMake:        carMake,
		CarModel:       carModel,
		CarColor:       carColor,
		Status:         "CREATED",

		ManufacturingDate: manufacturingDate,
		ManufacturerPrice: manufacturerPrice,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carId, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *CarChainCode) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *CarChainCode) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// Manufecturer ship the car to dealer. This method updates the shipment details for given carId in world state
func (s *CarChainCode) ShipToDealer(ctx contractapi.TransactionContextInterface, carId string, dealerId string, shippingPrice int, role string) error {

	if role != "manufacturer" {
		return fmt.Errorf("Failed to put to world state due to unauthorized user")
	}
	car, err := s.QueryCar(ctx, carId)
	if err != nil {
		return err
	}
	car.DealerId = dealerId
	car.Status = "SHIPPED"
	car.ShippingDate = time.Now().Format("2022-04-10 15:04:05")
	car.ShippingPrice = shippingPrice
	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carId, carAsBytes)
}

// Delear received the shipment and updates the delivery details for given carId in world state
func (s *CarChainCode) ReceiveDelivery(ctx contractapi.TransactionContextInterface, carId string, role string) error {
	if role != "dealer" {
		return fmt.Errorf("Failed to put to world state due to unauthorized user")
	}
	car, err := s.QueryCar(ctx, carId)
	if err != nil {
		return err
	}
	car.Status = "READY_FOR_SALE"
	car.DeliveryDate = time.Now().Format("2022-04-10 15:04:05")
	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carId, carAsBytes)
}

// Delear sell the car to customer and updates the sell details for given carId in world state
func (s *CarChainCode) SellToCustomer(ctx contractapi.TransactionContextInterface, carId string, consumerId string, customerPrice int, role string) error {
	if role != "dealer" {
		return fmt.Errorf("Failed to put to world state due to unauthorized user")
	}

	car, err := s.QueryCar(ctx, carId)
	if err != nil {
		return err
	}
	car.Status = "SOLD"
	car.ConsumerId = consumerId
	car.SoldOnDate = time.Now().Format("2022-04-10 15:04:05")
	car.CustomerPrice = customerPrice

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carId, carAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(CarChainCode))

	if err != nil {
		fmt.Printf("Error while creating Car Chain Code: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Car Chain Code: %s", err.Error())
	}
}
