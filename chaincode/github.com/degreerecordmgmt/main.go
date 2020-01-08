package main

import "github.com/hyperledger/fabric/core/chaincode/shim"

var _mainLogger = shim.NewLogger("SmartContractMain")

func main() {
	err := shim.Start(new(EduNetSmartContract))
	if err != nil {
		_mainLogger.Criticalf("Error starting  chaincode: %v", err)
	}
}
