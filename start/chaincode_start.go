/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// https://github.com/ccooper21/learn-chaincode/blob/c996d333b0c36a257464d1a172edb6b1bb5213e4/start/chaincode_start.go
package main
import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim" //the code that interfaces your golang code with a peer.
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// Invoke is called when you want to call chaincode functions to do real work.
	// Invocation transactions will be captured as blocks on the chain.
	// The structure of Invoke is simple. It receives a function argument and based on this argument calls Go functions in the chaincode.
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)  //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
// You will use Query to read the value of your chaincode state's key/value pairs.
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {
		//read a variable
		fmt.Println("hi there " + function)           //error
		return nil, nil;
	} else if function == "read" {                            //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)   //error
	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// M
// ============================================================================================================================

// Init resets all the things. Init is called when you first deploy your chaincode.
// As the name implies, this function should be used to do any initialization your chaincode needs.
// In our example, we use Init to configure the initial state of one variable on the ledger.
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	// stores the second element in the args argument to the key "financial_instrument".
	err = stub.PutState("financial_instrument", []byte(args[0]))
	// This is done by using the shim function stub.PutState.
	// The first argument is the key as a string, and the second argument is the value as an array of bytes.
	if err != nil {
		// This function may return an error which our code inspects and returns if present.
		return nil, err
	}
	err = stub.PutState("ISIN", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0]   //rename for fun
	value = args[1]
	// This write function should look similar to the Init change you just did.
	// One major difference is that you can now set the key and value for PutState.
	// This function allows you to store any key/value pair you want into the blockchain ledger.
	err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//This shim function just takes one string argument. The argument is the name of the key to retrieve.
// Next, this function returns the value as an array of bytes back to Query, who in turn sends it back to the REST handler.
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}
	return valAsbytes, nil
}