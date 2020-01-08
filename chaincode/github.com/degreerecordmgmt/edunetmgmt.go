package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	id "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var _scLogger = shim.NewLogger("EduNetSmartContract")

//StudentIDInfo represents student identity
type StudentIDInfo struct {
	Obj        string `json:"objType"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	DOB        string `json:"dob"` //Date in YYYY-MM-DD format
	UUID       string `json:"uuid"`
	ApprovedBy string `json:"approvedBy"`
	IsApproved bool   `json:"isApproved"`
}

//InstituteIDInfo represents participating institute identity
type InstituteIDInfo struct {
	Obj     string `json:"objType"`
	Name    string `json:"name"`
	Address string `json:"addr"`
	UUID    string `json:"uuid"`
}

//DegreeOffered represents degree details
type DegreeOffered struct {
	Obj       string `json:"objType"`
	DegreeID  string `json:"uuid"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	TCH       int    `json:"totalCreditHours"`
	OfferedBy string `json:"instuuid"`
	ValidFrom string `json:"validFrom"` //Date in YYYY-MM-DD format
	ValidUpto string `json:"validUpto"` //Date in YYYY-MM-DD format
	CreatedBy string `json:"createdBy"`
}

//StudentDegree represents degree ownedby a student
type StudentDegree struct {
	Obj       string  `json:"objType"`
	StudentID string  `json:"studuuid"`
	DegreeID  string  `json:"degreeuuid"`
	CGPA      float64 `json:"cgpa"`
	ValidFrom string  `json:"validFrom"` //Date in YYYY-MM-DD format
	UUID      string  `json:"uuid"`
	CreatedBy string  `json:"createdBy"`
}

//EduNetSmartContract implments the degree and student management smart contract
type EduNetSmartContract struct {
}

// Init initializes chaincode.
func (sc *EduNetSmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_scLogger.Infof("Inside the init method ")

	return shim.Success(nil)
}
func (sc *EduNetSmartContract) probe(stub shim.ChaincodeStubInterface) pb.Response {
	ts := ""
	_scLogger.Info("Inside probe method")
	tst, err := stub.GetTxTimestamp()
	if err == nil {
		ts = tst.String()
	}
	output := "{\"status\":\"Success\",\"ts\" : \"" + ts + "\" }"
	_scLogger.Info("Retuning " + output)
	return shim.Success([]byte(output))
}

func (sc *EduNetSmartContract) saveEntry(stub shim.ChaincodeStubInterface, objToSave interface{}, id string) pb.Response {
	key := id
	existingData, _ := stub.GetState(key)
	if existingData != nil && len(existingData) > 0 {
		return shim.Error("Object id already exists")
	}
	jsonBytesToStore, _ := json.Marshal(objToSave)
	stub.PutState(key, jsonBytesToStore)

	return shim.Success([]byte(jsonBytesToStore))
}

func (sc *EduNetSmartContract) queryObjectByID(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Invalid number of arguments")
	}
	key := args[0]
	data, err := stub.GetState(key)
	if err != nil {
		return shim.Success(nil)

	}

	return shim.Success(data)
}

//Invoke is the entry point for any transaction
func (sc *EduNetSmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var response pb.Response
	action, _ := stub.GetFunctionAndParameters()
	switch action {
	case "probe":
		response = sc.probe(stub)

	case "queryById":
		response = sc.queryObjectByID(stub)
	default:
		response = shim.Error("Invalid action provoided")
	}
	return response
}
func (sc *EduNetSmartContract) getInvokerIdentity(stub shim.ChaincodeStubInterface) (bool, string) {
	//Following id comes in the format X509::<Subject>::<Issuer>>
	/*enCert, err := id.GetX509Certificate(stub)
	if err != nil {
		return false, "Unknown."
	}*/

	mspID, err := id.GetMSPID(stub)
	if err != nil {
		return false, "Unknown."
	}
	return true, fmt.Sprintf("%s", mspID)

}
func (sc *EduNetSmartContract) getTrxnTS(stub shim.ChaincodeStubInterface) string {
	txTime, err := stub.GetTxTimestamp()
	if err != nil {
		return "0000.00.00.00.00.000"
	}
	var ts time.Time
	newTS := ts.Add(time.Duration(txTime.Seconds) * time.Second)
	return newTS.Format("2006.01.02.15.04.05.000")

}
func (sc *EduNetSmartContract) getOrganizationRole(stub shim.ChaincodeStubInterface) string {
	idOk, who := sc.getInvokerIdentity(stub)
	if !idOk {
		_scLogger.Error("Unable to retrive the invoker ID")
		return ""
	}
	key := fmt.Sprintf("PARTICIPANT_%s", who)
	_scLogger.Infof("User key %s", key)
	if roleJSON, err := stub.GetState(key); err == nil {
		_scLogger.Infof("User key %s", string(roleJSON))
		role := string(roleJSON)
		return role
	}
	_scLogger.Error("Unable to retrive the role , not registered")
	return ""

}
