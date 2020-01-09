package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	id "github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var _scLogger = shim.NewLogger("EduNetSmartContract")
var _educationalInstitutes = make(map[string]bool)

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
	TCH       string `json:"totalCreditHours"`
	OfferedBy string `json:"offeredByDept"`
	ValidFrom string `json:"validFrom"` //Date in YYYY-MM-DD format
	ValidUpto string `json:"validUpto"` //Date in YYYY-MM-DD format
	CreatedBy string `json:"createdBy"`
}

//StudentDegree represents degree ownedby a student
type StudentDegree struct {
	Obj       string `json:"objType"`
	StudentID string `json:"studuuid"`
	DegreeID  string `json:"degreeuuid"`
	CGPA      string `json:"cgpa"`
	ValidFrom string `json:"validFrom"` //Date in YYYY-MM-DD format
	UUID      string `json:"uuid"`
	CreatedBy string `json:"createdBy"`
}

//EduNetSmartContract implments the degree and student management smart contract
type EduNetSmartContract struct {
}

// Init initializes chaincode.
func (sc *EduNetSmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_scLogger.Infof("Inside the init method ")
	_educationalInstitutes["IITJMSP"] = true
	_educationalInstitutes["IITKJPMSP"] = true

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

func (sc *EduNetSmartContract) saveEntry(stub shim.ChaincodeStubInterface, objToSave interface{}, id string, isOverwrite bool) pb.Response {
	key := id
	existingData, _ := stub.GetState(key)
	if existingData != nil && len(existingData) > 0 && !isOverwrite {
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
	case "registerStudent":
		response = sc.registerStudent(stub)
	case "approveStudent":
		response = sc.approveStudent(stub)
	case "queryById":
		response = sc.queryObjectByID(stub)
	case "modInstitueRegInfo":
		response = sc.upsertInstituteEntry(stub)
	case "modDegreeInfo":
		response = sc.upsertDegreeDetails(stub)
	case "registerDegree":
		response = sc.registerDegree(stub)
	case "searchDegreeByStudentID":
		response = sc.getDegreesByStudentID(stub)
	case "searchDegreesOffered":
		response = sc.getDegreesOffered(stub)
	default:
		response = shim.Error("Invalid action provoided")
	}
	return response
}
func (sc *EduNetSmartContract) registerStudent(stub shim.ChaincodeStubInterface) pb.Response {
	isFound, id := sc.getInvokerIdentity(stub)
	if !isFound {
		_scLogger.Errorf("Identity of the user could not be retrived")
		return shim.Error("Identity of the user could not be retrieved")
	}
	if id != "EDUNETMSP" {
		_scLogger.Errorf("Trxn not allowed")
		return shim.Error("Trxn not allowed")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Student details not provided")
	}
	var student StudentIDInfo
	err := json.Unmarshal([]byte(args[0]), &student)
	if err != nil {
		return shim.Error("Student JSON could not be parsed")
	}
	if len(strings.TrimSpace(student.Email)) == 0 {
		return shim.Error("EmailID could not be blank")
	}
	student.UUID = student.Email
	student.IsApproved = false
	student.ApprovedBy = ""
	student.Obj = "edunet.student.reginfo"
	return sc.saveEntry(stub, student, student.UUID, false)
}
func (sc *EduNetSmartContract) approveStudent(stub shim.ChaincodeStubInterface) pb.Response {
	isFound, id := sc.getInvokerIdentity(stub)
	if !isFound {
		_scLogger.Errorf("Identity of the user could not be retrived")
		return shim.Error("Identity of the user could not be retrieved")
	}
	if _, isFound := _educationalInstitutes[id]; isFound {
		_scLogger.Errorf("Trxn not allowed")
		return shim.Error("Trxn not allowed")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Student details not provided")
	}
	emailID := args[0]

	if len(strings.TrimSpace(emailID)) == 0 {
		return shim.Error("EmailID could not be blank")
	}
	recBytes, err := stub.GetState(emailID)
	if err != nil || recBytes == nil || len(recBytes) == 0 {
		return shim.Error("Error retrival of student record")
	}
	var student StudentIDInfo

	err = json.Unmarshal(recBytes, &student)
	if err != nil {
		return shim.Error("Student JSON could not be parsed")
	}
	student.IsApproved = true
	student.ApprovedBy = id
	return sc.saveEntry(stub, student, student.UUID, true)
}

func (sc *EduNetSmartContract) upsertInstituteEntry(stub shim.ChaincodeStubInterface) pb.Response {
	isFound, id := sc.getInvokerIdentity(stub)
	if !isFound {
		_scLogger.Errorf("Identity of the user could not be retrived")
		return shim.Error("Identity of the user could not be retrieved")
	}
	if _, isFound := _educationalInstitutes[id]; isFound {
		_scLogger.Errorf("Trxn not allowed")
		return shim.Error("Trxn not allowed")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 2 {
		return shim.Error("Institude details not provided")
	}
	name := args[0]
	addr := args[1]
	if len(strings.TrimSpace(name)) == 0 && len(strings.TrimSpace(addr)) == 0 {
		return shim.Error("Name or Address of the institute could not be blank")
	}
	var instDetails InstituteIDInfo
	instDetails.Obj = "edunet.inst.reginfo"
	instDetails.Name = name
	instDetails.Address = addr
	instDetails.UUID = id
	return sc.saveEntry(stub, instDetails, instDetails.UUID, true)

}

func (sc *EduNetSmartContract) registerDegree(stub shim.ChaincodeStubInterface) pb.Response {
	isFound, id := sc.getInvokerIdentity(stub)
	if !isFound {
		_scLogger.Errorf("Identity of the user could not be retrived")
		return shim.Error("Identity of the user could not be retrieved")
	}
	if _, isFound := _educationalInstitutes[id]; isFound {
		_scLogger.Errorf("Trxn not allowed")
		return shim.Error("Trxn not allowed")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Degree details not provided")
	}

	var degreeDetails StudentDegree
	degreeDetails.Obj = "edunet.student.degree"

	studentUUID := degreeDetails.StudentID
	degreeID := degreeDetails.DegreeID
	//check for student id

	if len(strings.TrimSpace(studentUUID)) == 0 {
		return shim.Error("StudentID could not be blank")
	}
	if len(strings.TrimSpace(degreeID)) == 0 {
		return shim.Error("DegreeID could not be blank")
	}
	recBytes, err := stub.GetState(studentUUID)
	if err != nil || recBytes == nil || len(recBytes) == 0 {
		return shim.Error("Error retrival of student record")
	}
	var student StudentIDInfo

	err = json.Unmarshal(recBytes, &student)
	if err != nil {
		return shim.Error("Student JSON could not be parsed")
	}
	degreeInfoBytes, err := stub.GetState(degreeID)
	if err != nil || degreeInfoBytes == nil || len(degreeInfoBytes) == 0 {
		return shim.Error("Error retrival of degree info record")
	}
	var degInfo DegreeOffered
	err = json.Unmarshal(degreeInfoBytes, &degInfo)
	if err != nil {
		return shim.Error("Degree details JSON could not be parsed")
	}
	//At this point all the inputs are fine
	degreeDetails.ValidFrom = sc.getTrxnTS(stub)
	degreeDetails.UUID = stub.GetTxID()
	degreeDetails.CreatedBy = id
	return sc.saveEntry(stub, degreeDetails, degreeDetails.UUID, false)

}
func (sc *EduNetSmartContract) upsertDegreeDetails(stub shim.ChaincodeStubInterface) pb.Response {
	isFound, id := sc.getInvokerIdentity(stub)
	if !isFound {
		_scLogger.Errorf("Identity of the user could not be retrived")
		return shim.Error("Identity of the user could not be retrieved")
	}
	if _, isFound := _educationalInstitutes[id]; isFound {
		_scLogger.Errorf("Trxn not allowed")
		return shim.Error("Trxn not allowed")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Degree details not provided")
	}

	var degreeDetails DegreeOffered
	err := json.Unmarshal([]byte(args[0]), &degreeDetails)
	if err != nil {
		return shim.Error("Unable to parse input json")
	}
	degreeDetails.Obj = "edunet.inst.degreeinfo"
	if len(degreeDetails.DegreeID) == 0 {
		//New degree to the added in the system
		degreeDetails.CreatedBy = id
		degreeDetails.DegreeID = stub.GetTxID()
		return sc.saveEntry(stub, degreeDetails, degreeDetails.DegreeID, false)
	}

	recBytes, err := stub.GetState(degreeDetails.DegreeID)
	if err != nil || recBytes == nil || len(recBytes) == 0 {
		return shim.Error("Error retrival of degree record")
	}
	var existingRecord DegreeOffered

	err = json.Unmarshal(recBytes, &existingRecord)
	if err != nil {
		return shim.Error("Degree details JSON could not be parsed")
	}
	if existingRecord.CreatedBy != id {
		return shim.Error("Can not update degree created by other institution")
	}
	existingRecord.Name = degreeDetails.Name
	existingRecord.TCH = degreeDetails.TCH
	existingRecord.Type = degreeDetails.Type
	existingRecord.ValidFrom = degreeDetails.ValidFrom
	existingRecord.ValidUpto = degreeDetails.ValidUpto
	existingRecord.OfferedBy = degreeDetails.OfferedBy
	return sc.saveEntry(stub, existingRecord, existingRecord.DegreeID, true)

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

func (sc *EduNetSmartContract) getDegreesByStudentID(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Invalid number of arguments")
	}
	key := args[0]
	query := map[string]interface{}{
		"objType":  "edunet.student.degree",
		"studuuid": key,
	}
	records := sc.retriveRecords(stub, query)
	data, _ := json.Marshal(records)
	return shim.Success(data)
}
func (sc *EduNetSmartContract) getDegreesOffered(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Invalid number of arguments")
	}
	key := args[0]
	query := map[string]interface{}{
		"objType":   "edunet.inst.degreeinfo",
		"createdBy": key,
	}
	records := sc.retriveRecords(stub, query)
	data, _ := json.Marshal(records)
	return shim.Success(data)
}

func (sc *EduNetSmartContract) retriveRecords(stub shim.ChaincodeStubInterface, query map[string]interface{}) []map[string]interface{} {
	records := make([]map[string]interface{}, 0)
	queryBytes, _ := json.Marshal(query)
	selectorString := fmt.Sprintf("{\"selector\":%s }", string(queryBytes))
	_scLogger.Info("Query Selector :" + selectorString)
	resultsIterator, _ := stub.GetQueryResult(selectorString)
	for resultsIterator.HasNext() {
		record := make(map[string]interface{})
		recordBytes, _ := resultsIterator.Next()
		err := json.Unmarshal(recordBytes.Value, &record)
		if err != nil {
			_scLogger.Infof("Unable to unmarshal data retived:: %v", err)
		}
		records = append(records, record)
	}
	return records
}
