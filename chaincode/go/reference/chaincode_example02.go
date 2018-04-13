
package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/pem"
	"crypto/x509"
	"strings"
)

var logger = shim.NewLogger("SimpleChaincode")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init")

	_, args := stub.GetFunctionAndParameters()
	var a, b string    // Entities
	var aVal, bVal int // Asset holdings
	var aSub, bSub string // Asset holdings
	var err error

	if len(args) != 4 {
		return pb.Response{Status:403, Message:"Incorrect number of args. Expecting 4"}
	}

	// Initialize the chaincode
	a = args[0]
	aVal, err = strconv.Atoi(args[1])
	aSub = "testa"
	if err != nil {
		return pb.Response{Status:403, Message:"Expecting integer value for asset holding"}
	}
	b = args[2]
	bVal, err = strconv.Atoi(args[3])
	bSub = "testb"
	if err != nil {
		return pb.Response{Status:403, Message:"Expecting integer value for asset holding"}
	}
	logger.Debugf("aVal, bVal = %d aSub,bSub = %s", aVal, bVal,aSub,bSub)

	// Write the state to the ledger
	err = stub.PutState(a, []byte(strconv.Itoa(aVal)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("asubscriptions", []byte(aSub))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(b, []byte(strconv.Itoa(bVal)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("bsubscriptions", []byte(bSub))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("security", []byte("security1"))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("security1",[]byte("."))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Invoke")

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	name, org := getCreator(creatorBytes)

	logger.Debug("transaction creator " + name + "@" + org)

	function, args := stub.GetFunctionAndParameters()
	if function == "move" {
		// Make payment of x units from a to b
		return t.move(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemented in invoke
		return t.query(stub, args)
	} else if function == "subscribe" {
		// the old "Query" is now implemented in invoke
		return t.subscribe(stub, args)
	} else if function == "unsubscribe" {
		// the old "Query" is now implemented in invoke
		return t.unsubscribe(stub, args)
	} else if function == "addsecurity" {
		// the old "Query" is now implemented in invoke
		return t.addsecurity(stub, args)
	}

	return pb.Response{Status:403, Message:"Invalid invoke function name."}
}

// Transaction makes payment of x units from a to b
func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var a, b string    // Entities
	var aVal, bVal int // Asset holdings
	var x int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of hjuh args. Expecting 3")
	}

	a = args[0]
	b = args[1]

	// Get the state from the ledger
	aBytes, err := stub.GetState(a)
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("Entity not found")
	}
	aVal, _ = strconv.Atoi(string(aBytes))

	bBytes, err := stub.GetState(b)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if bBytes == nil {
		return shim.Error("Entity not found")
	}
	bVal, _ = strconv.Atoi(string(bBytes))

	// Perform the execution
	x, err = strconv.Atoi(args[2])
	if err != nil {
		return pb.Response{Status:403, Message:"Invalid transaction amount, expecting a integer value"}
	}
	aVal = aVal - x -10
	bVal = bVal + x
	logger.Debug("aVal = %d, bVal = %d\n", aVal, bVal)

	// Write the state back to the ledger
	err = stub.PutState(a, []byte(strconv.Itoa(aVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(b, []byte(strconv.Itoa(bVal)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Subscribe to a security
func (t *SimpleChaincode) subscribe(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var a string    // Entities
	var aSub string // Security to Subscribe
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of hjuh args. Expecting 2")
	}
	
	a = args[0]
	aSub = args[1]
	// Get the state from the ledger
	aBytes, err := stub.GetState(a+"subscriptions")
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("Entity not found")
	}
	
	securityBytes, err := stub.GetState("security")
	if err != nil {
		return shim.Error(err.Error())
	}
	if securityBytes == nil {
		return shim.Error("Entity not found")
	}
	aSubString := string(aBytes)
	securityString :=string(securityBytes)
	
	if(strings.Contains(aSubString,aSub+",")){
		return shim.Error("Already Subscribed")
	}
	if(!strings.Contains(securityString,aSub+",")){
		return shim.Error("Invalid Security")
	}
	logger.Debug("aSub = %s", aSub)

	// Write the state back to the ledger
	err = stub.PutState(a+"subscriptions",append(append([]byte(aSub),[]byte(",")...),aBytes...))
	if err != nil {
		return shim.Error(err.Error())
	}
	// add subs to security
	securitySubBytes, err := stub.GetState(aSub)
		if err != nil {
			err1 := stub.PutState(aSub,[]byte(a))
	if err1 != nil {
		return shim.Error(err1.Error())
	}
	} else {
		//securitySubString :=string(securitySubBytes)
		err2 := stub.PutState(aSub,append(append([]byte(a),[]byte(",")...),securitySubBytes...))
			if err2 != nil {
		return shim.Error(err2.Error())
	}
	}
	return shim.Success(nil)
}
//add new security
// Subscribe to a security
func (t *SimpleChaincode) addsecurity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var a string    // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number args")
	}
	
	a = "security"
	security := args[0]

	securityBytes, err := stub.GetState("security")
	if err != nil {
		return shim.Error(err.Error())
	}
	if securityBytes == nil {
		return shim.Error("Entity not found")
	}
	securityString :=string(securityBytes)

	if(strings.Contains(securityString,security+",")){
		return shim.Error("Security Already Added")
	}
	
	logger.Debug("Security = %s", security)

	// Write the state back to the ledger
	err = stub.PutState(a,append(append([]byte(security),[]byte(",")...),securityBytes...))
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(nil)
}
// Subscribe to a security
func (t *SimpleChaincode) unsubscribe(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var a string    // Entities
	var aSub string // Asset holdings
	//var x int          // Transaction value
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of hjuh args. Expecting 2")
	}
	
	a = args[0]
	aSub = args[1]
	// Get the state from the ledger
	aBytes, err := stub.GetState(a+"subscriptions")
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("Entity not found")
	}
	aSubString := string(aBytes)
	if(!strings.Contains(aSubString,aSub+",")){
		return shim.Error("You are not Subscribed to this security")
	}
	aSubString = strings.Replace(aSubString, aSub+",", "", 1)

	// Write the state back to the ledger
	err = stub.PutState(a+"subscriptions",[]byte(aSubString))
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(nil)
}

// deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return pb.Response{Status:403, Message:"Incorrect number of args"}
	}

	a := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(a)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// read value
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var a string // Entities
	var err error

	//if len(args) != 1 {
	//	return pb.Response{Status:403, Message:"Incorrect number of args"}
	//}

	a = args[0]

	// Get the state from the ledger
	valBytes, err := stub.GetState(a)
	if err != nil {
		return shim.Error(err.Error())
	}

	if valBytes == nil {
		return shim.Error("Entity not found")
	}

	return shim.Success(valBytes)
}

var getCreator = func (certificate []byte) (string, string) {
	data := certificate[strings.Index(string(certificate), "-----"): strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode([]byte(data))
	cert, _ := x509.ParseCertificate(block.Bytes)
	organization := cert.Issuer.Organization[0]
	commonName := cert.Subject.CommonName
	logger.Debug("commonName: " + commonName + ", organization: " + organization)

	organizationShort := strings.Split(organization, ".")[0]

	return commonName, organizationShort
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Error(err.Error())
	}
}
