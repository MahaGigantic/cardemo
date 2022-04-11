/*
Copyright 2022 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("CarDemoappUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate CarDemoappUser wallet contents: %s\n", err)
			os.Exit(1)
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
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("cardemo")

	result, err := contract.EvaluateTransaction("QueryAllCars")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
	// Call createNewCar Function and supply paramters like manufacturerId string, carId string, carMake string, carModel string, carColor string, manufacturingDate string, manufacturerPrice int, role string
	result, err = contract.SubmitTransaction("createNewCar", "MOrg03", "M105", "2022", "MOrg03CM101", "White", time.Now().String(), "450000", "manufacturer")
	if err != nil {
		fmt.Printf("Failed to submit  createNewCar transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call QueryCar Function and by supplying CarID paramter
	result, err = contract.EvaluateTransaction("QueryCar", "M105")
	if err != nil {
		fmt.Printf("Failed to evaluate QueryCar transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call ShipToDealer Function and supply paramters like carId string, dealerId string, shippingPrice int, role string
	result, err = contract.SubmitTransaction("ShipToDealer", "M105", "D101", "12000", "manufacturer")
	if err != nil {
		fmt.Printf("Failed to submit ShipToDealer transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call QueryCar Function and by supplying CarID paramter
	result, err = contract.EvaluateTransaction("QueryCar", "M105")
	if err != nil {
		fmt.Printf("Failed to evaluate QueryCar transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call ReceiveDelivery Function and supply paramters like carId string, role string
	result, err = contract.SubmitTransaction("ReceiveDelivery", "M105", "dealer")
	if err != nil {
		fmt.Printf("Failed to submit ReceiveDelivery transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call QueryCar Function and by supplying CarID paramter
	result, err = contract.EvaluateTransaction("QueryCar", "M105")
	if err != nil {
		fmt.Printf("Failed to evaluate QueryCar transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call SellToCustomer Function and supply paramters like carId string, consumerId string, customerPrice int, role string
	result, err = contract.SubmitTransaction("SellToCustomer", "M105", "CUST103", "950000", "dealer")
	if err != nil {
		fmt.Printf("Failed to submit SellToCustomer transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	// Call QueryCar Function and by supplying CarID paramter
	result, err = contract.EvaluateTransaction("QueryCar", "M105")
	if err != nil {
		fmt.Printf("Failed to evaluate QueryCar transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

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
