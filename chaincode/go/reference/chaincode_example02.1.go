
package main

import (
	"fmt"
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
	err = stub.PutState("security", []byte("security1,"))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("security1",[]byte(""))
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
		return t.subscribe(stub, args,  org)
	} else if function == "unsubscribe" {
		// the old "Query" is now implemented in invoke
		return t.unsubscribe(stub, args,  org)
	} else if function == "addsecurity" {
		// the old "Query" is now implemented in invoke
		return t.addsecurity(stub, args)
	} else if function == "mysubscriptions" {
		// the old "Query" is now implemented in invoke
		return t.mysubscriptions(stub,  org)
	} else if function == "issueCa" {
		// the old "Query" is now implemented in invoke
		return t.issueCa(stub,args)
	} else if function == "myCa" {
		// the old "Query" is now implemented in invoke
		return t.myCa(stub,org)
	} else if function == "myCaSecurity" {
		// the old "Query" is now implemented in invoke
		return t.myCaSecurity(stub,args,org)
	} else if function == "subscribe2" {
		// the old "Query" is now implemented in invoke
		return t.subscribe2(stub,args,org)
	} else if function == "subscribe3" {
		// the old "Query" is now implemented in invoke
		return t.subscribe3(stub,args,org)
	} else if function == "issueCa3" {
		// the old "Query" is now implemented in invoke
		return t.issueCa3(stub,args)
	} else if function == "mysubscribe3get" {
		// the old "Query" is now implemented in invoke
		return t.mysubscribe3get(stub,args,org)
	} else if function == "myCa3" {
		// the old "Query" is now implemented in invoke
		return t.myCa3(stub,args,org)
	} else if function == "myCaSecurity3" {
		// the old "Query" is now implemented in invoke
		return t.myCaSecurity3(stub,args,org)
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
func (t *SimpleChaincode) subscribe(stub shim.ChaincodeStubInterface, args []string, name string) pb.Response {
	var a string    // Entities
	var aSub string // Security to Subscribe
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	
	//a = args[0]
	a = name
	aSub = args[0]
	// Get the state from the ledger
	aBytes, err := stub.GetState(a+"subscriptions")
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		//return shim.Error("Entity not found")
	}
	
	securityBytes, err := stub.GetState("security")
	if err != nil {
		return shim.Error(err.Error())
	}
	if securityBytes == nil {
		return shim.Error(" security Entity not found")
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
// Subscribe to a security
func (t *SimpleChaincode) subscribe2(stub shim.ChaincodeStubInterface, args []string, name string) pb.Response {
	var a string    // Entities
	var aSub string // Security to Subscribe
	

	if len(args) != 1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	
	//a = args[0]
	a = name
	aSub = args[0]
	// Write the state back to the ledger
   stub.PutState(a+"subscriptions"+aSub,[]byte(aSub))
	return shim.Success(nil)
}
//____________________________________________________________________________
// Subscribe to a security
func (t *SimpleChaincode) subscribe3(stub shim.ChaincodeStubInterface, args []string, name string) pb.Response {
	var a string    // Entities
	var aSub string // Security to Subscribe
	

	if len(args) != 1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	
	//a = args[0]
	a = name
	aSub = args[0]
	// Write the state back to the ledger
	stub.PutState(a+"subscriptions"+aSub,[]byte(aSub))
   
	//txid := stub.GetTxID()
	compositeIndexName := "security~org"
	compositeIndexName2 := "org~security"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{aSub,a})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}
	
	compositeKey2, compositeErr2 := stub.CreateCompositeKey(compositeIndexName2, []string{a,aSub})
	if compositeErr2 != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr2.Error()))
	}
	// Save the composite key index sec
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}
	//save sub 
	compositePutErr2 := stub.PutState(compositeKey2, []byte{0x00})
	if compositePutErr2 != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr2.Error()))
	}
	return shim.Success([]byte(fmt.Sprintf("Successfully added %s%s to %s", a,aSub)))
}
//--------------
func (s *SimpleChaincode) issueCa3(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check we have a valid number of args
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, expecting 2")
	}
	aSub := args[0]
	caData := args[1]
	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("security~org", []string{aSub})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", aSub, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", aSub))
	}

	// Iterate through result set and compute final value
	var finalVal string
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}
		// Retrieve the delta value and operation
		value := keyParts[1]
		compositeIndexName := "caorg~security~data"
		compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, []string{"ca"+value,aSub,caData})
		if compositeErr != nil {
			return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", aSub, compositeErr.Error()))
		}
		err1 := APIstub.PutState(compositeKey,[]byte(caData))
		if err1 != nil {
		return shim.Error(err1.Error())
        }
		finalVal = finalVal +","+ value
	}

	return shim.Success([]byte(finalVal))
}
func (s *SimpleChaincode) mysubscribe3get(APIstub shim.ChaincodeStubInterface, args []string,name string) pb.Response {
	// Check we have a valid number of args

	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("org~security", []string{name})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}

	// Iterate through result set and compute final value
	var finalVal string
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		value := keyParts[1]
		
		finalVal = finalVal +","+ value
	}

	return shim.Success([]byte(finalVal))
}


//__________________________________________________________________________________________
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
// unSubscribe to a security
func (t *SimpleChaincode) unsubscribe(stub shim.ChaincodeStubInterface, args []string,name string) pb.Response {
	var a string    // Entities
	var aSub string // Asset holdings
	//var x int          // Transaction value
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	a = name//args[0]
	aSub = args[0]
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
	//remove security from security list
	securitySubBytes, err := stub.GetState(aSub)
	securityString :=string(securitySubBytes)
	securityString = strings.Replace(securityString, a+",", "", 1)
	err = stub.PutState(aSub,[]byte(securityString))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
// get mysecurity
func (t *SimpleChaincode) mysubscriptions(stub shim.ChaincodeStubInterface,name string) pb.Response {
	var a string    // Entities
	//var x int          // Transaction value
	var err error
	a = name//args[0]
	// Get the state from the ledger
	aBytes, err := stub.GetState(a+"subscriptions")
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("Entity not found")
	}
	return shim.Success(aBytes)
}
//-----------------------------
// get my ca
func (t *SimpleChaincode) myCa3(APIstub shim.ChaincodeStubInterface,args []string,name string) pb.Response {
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	//page,err := strconv.ParseInt(args[0],10,0)
	caName := "ca"+name
	//var cnt int64;
	//cnt = 100
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("caorg~security~data", []string{caName})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", caName, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("myca3--No variable by the name %s exists", caName))
	}

	// Iterate through result set and compute final value
	var finalVal string
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		value := keyParts[1]
		data := keyParts[2]
		finalVal = finalVal +","+ value +"::"+data
		//if i >= (page*cnt+cnt) { break }
	}

	return shim.Success([]byte(finalVal))
}
// get my ca security
func (t *SimpleChaincode) myCaSecurity3(stub shim.ChaincodeStubInterface,args []string,name string) pb.Response {
	var a string    // Entities
	//var x int          // Transaction value
	var err error
	
	if len(args) !=1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	a = name//args[0]
	aSub := args[0]
	compositeIndexName2 := "caorg~security"
	compositeKey2, compositeErr2 := stub.CreateCompositeKey(compositeIndexName2, []string{a,aSub})
	if compositeErr2 != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr2.Error()))
	}
	// Get the state from the ledger
	aBytes, err := stub.GetState(compositeKey2)
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("you are not subscribed to this security")
	}
	
	return shim.Success(aBytes)
}
//--------------------------------
// get my ca
func (t *SimpleChaincode) myCa(stub shim.ChaincodeStubInterface,name string) pb.Response {
	var a string    // Entities
	//var x int          // Transaction value
	var err error
	a = name//args[0]
	// Get the state from the ledger
	aBytes, err := stub.GetState(a+"-ca")
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("Entity not found")
	}
	return shim.Success(aBytes)
}
// get my ca security
func (t *SimpleChaincode) myCaSecurity(stub shim.ChaincodeStubInterface,args []string,name string) pb.Response {
	var a string    // Entities
	//var x int          // Transaction value
	var err error
	
	if len(args) !=1 {
		return shim.Error("Incorrect number of hjuh args. Expecting 1")
	}
	a = name//args[0]
	aSub := args[0]
	// Get the state from the ledger
	aBytes, err := stub.GetState(a+"-ca-"+aSub)
	if err != nil {
		return shim.Error(err.Error())
	}
	if aBytes == nil {
		return shim.Error("you are not subscribed to this security")
	}
	
	return shim.Success(aBytes)
}
// // issue CA
func (t *SimpleChaincode) issueCa(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	var aSub string // security
	//var x int          // Transaction value
	var err error
		
	if len(args) != 2 {
		return shim.Error("Incorrect number of hjuh args. Expecting 2")
	}
	aSub = args[0]
	caData := args[1]
	
	securityBytes, err := stub.GetState("security")
	securityString :=string(securityBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	if(!strings.Contains(securityString,aSub+",")){
		return shim.Error("Invalid Security")
	}
	securitySubBytes, err := stub.GetState(aSub)
	securitySubString :=string(securitySubBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	if securitySubBytes == nil {
		return shim.Error("No Subscribers")
	}
	if(len(securitySubString)<2){
		return shim.Error("No valid Subscribers")
	}
	result := strings.Split(securitySubString, ",")
	
	err2 := stub.PutState("ca-"+aSub,[]byte(caData))
	if err2 != nil {
	return shim.Error(err2.Error())
    }
	 // Display all elements.
	i :=0;
    for i = range result {
    /*err2 := stub.PutState("debug-"+aSub+strconv.Itoa(i),[]byte(strconv.Itoa(i)))
	if err2 != nil {
	return shim.Error(err2.Error())
    }*/
        aBytes, err := stub.GetState(result[i]+"-ca")
        if err != nil {
		return shim.Error(err.Error())
        }
		if(aBytes == nil){
		err1 := stub.PutState(result[i]+"-ca",[]byte(aSub))
		if err1 != nil {
		return shim.Error(err1.Error())
        }
		} else {
		stringAbytes := string(aBytes)
		err1 := stub.PutState(result[i]+"-ca",append(append([]byte(stringAbytes),","...),[]byte(aSub)...))
		if err1 != nil {
		return shim.Error(err1.Error())
        }
		}
		err2 = stub.PutState(result[i]+"-ca-"+aSub,[]byte(caData))
		if err2 != nil {
		return shim.Error(err2.Error())
        }
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

// // read value
// func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface,st name, compositePutErr.Error()))
// 	}
// 	//save sub 
// 	compositePutErr2 := stub.PutState(compositeKey2, []byte{0x00})
// 	if compositePutErr2 != nil {
// 		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr2.Error()))
// 	}
// 	return shim.Success([]byte(fmt.Sprintf("Successfully added %s%s to %s", a,aSub)))
// }
//--------------
// func (s *SimpleChaincode) issueCa3(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	// Check we have a valid number of args
// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of arguments, expecting 2")
// 	}

// 	aSub := args[0]
// 	caData := args[1]
// 	// Get all deltas for the variable
// 	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("security~org", []string{aSub})
// 	if deltaErr != nil {
// 		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", aSub, deltaErr.Error()))
// 	}
// 	defer deltaResultsIterator.Close()

// 	// Check the variable existed
// 	if !deltaResultsIterator.HasNext() {
// 		return shim.Error(fmt.Sprintf("No variable by the name %s exists", aSub))
// 	}

// 	// Iterate through result set and compute final value
// 	var finalVal string
// 	var i int
// 	for i = 0; deltaResultsIterator.HasNext(); i++ {
// 		// Get the next row
// 		responseRange, nextErr := deltaResultsIterator.Next()
// 		if nextErr != nil {
// 			return shim.Error(nextErr.Error())
// 		}

// 		// Split the composite key into its component parts
// 		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
// 		if splitKeyErr != nil {
// 			return shim.Error(splitKeyErr.Error())
// 		}
// 		// Retrieve the delta value and operation
// 		value := keyParts[1]
// 		compositeIndexName := "caorg~security"
// 		compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{"ca"+a,aSub})
// 		if compositeErr != nil {
// 			return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", aSub, compositeErr.Error()))
// 		}
// 		err1 := APIstub.PutState(compositeKey,[]byte(caData))
// 		if err1 != nil {
// 		return shim.Error(err1.Error())
//         }
// 		finalVal = finalVal +","+ value
// 	}

// 	return shim.Success([]byte(finalVal))
// }
// func (s *SimpleChaincode) mysubscribe3get(APIstub shim.ChaincodeStubInterface, args []string,name string) pb.Response {
// 	// Check we have a valid number of args
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments, expecting 1")
// 	}

// 	// Get all deltas for the variable
// 	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey("org~security~txID", []string{name})
// 	if deltaErr != nil {
// 		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
// 	}
// 	defer deltaResultsIterator.Close()

// 	// Check the variable existed
// 	if !deltaResultsIterator.HasNext() {
// 		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
// 	}

// 	// Iterate through result set and compute final value
// 	var finalVal string
// 	var i int
// 	for i = 0; deltaResultsIterator.HasNext(); i++ {
// 		// Get the next row
// 		responseRange, nextErr := deltaResultsIterator.Next()
// 		if nextErr != nil {
// 			return shim.Error(nextErr.Error())
// 		}

// 		// Split the composite key into its component parts
// 		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
// 		if splitKeyErr != nil {
// 			return shim.Error(splitKeyErr.Error())
// 		}

// 		// Retrieve the delta value and operation
// 		value := keyParts[1]
		
// 		finalVal = finalVal +","+ value
// 	}

// 	return shim.Success([]byte(finalVal))
// }


// //__________________________________________________________________________________________
// //add new security
// // Subscribe to a security
// func (t *SimpleChaincode) addsecurity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	var a string    // Entities
// 	var err error

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number args")
// 	}
	
// 	a = "security"
// 	security := args[0]

// 	securityBytes, err := stub.GetState("security")
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if securityBytes == nil {
// 		return shim.Error("Entity not found")
// 	}
// 	securityString :=string(securityBytes)

// 	if(strings.Contains(securityString,security+",")){
// 		return shim.Error("Security Already Added")
// 	}
	
// 	logger.Debug("Security = %s", security)

// 	// Write the state back to the ledger
// 	err = stub.PutState(a,append(append([]byte(security),[]byte(",")...),securityBytes...))
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
	
// 	return shim.Success(nil)
// }
// // unSubscribe to a security
// func (t *SimpleChaincode) unsubscribe(stub shim.ChaincodeStubInterface, args []string,name string) pb.Response {
// 	var a string    // Entities
// 	var aSub string // Asset holdings
// 	//var x int          // Transaction value
// 	var err error

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of hjuh args. Expecting 1")
// 	}
	
// 	a = name//args[0]
// 	aSub = args[0]
// 	// Get the state from the ledger
// 	aBytes, err := stub.GetState(a+"subscriptions")
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if aBytes == nil {
// 		return shim.Error("Entity not found")
// 	}
// 	aSubString := string(aBytes)
// 	if(!strings.Contains(aSubString,aSub+",")){
// 		return shim.Error("You are not Subscribed to this security")
// 	}
// 	aSubString = strings.Replace(aSubString, aSub+",", "", 1)

// 	// Write the state back to the ledger
// 	err = stub.PutState(a+"subscriptions",[]byte(aSubString))
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	//remove security from security list
// 	securitySubBytes, err := stub.GetState(aSub)
// 	securityString :=string(securitySubBytes)
// 	securityString = strings.Replace(securityString, a+",", "", 1)
// 	err = stub.PutState(aSub,[]byte(securityString))
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success(nil)
// }
// // get mysecurity
// func (t *SimpleChaincode) mysubscriptions(stub shim.ChaincodeStubInterface,name string) pb.Response {
// 	var a string    // Entities
// 	//var x int          // Transaction value
// 	var err error
// 	a = name//args[0]
// 	// Get the state from the ledger
// 	aBytes, err := stub.GetState(a+"subscriptions")
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if aBytes == nil {
// 		return shim.Error("Entity not found")
// 	}
// 	return shim.Success(aBytes)
// }
// // get my ca
// func (t *SimpleChaincode) myCa(stub shim.ChaincodeStubInterface,name string) pb.Response {
// 	var a string    // Entities
// 	//var x int          // Transaction value
// 	var err error
// 	a = name//args[0]
// 	// Get the state from the ledger
// 	aBytes, err := stub.GetState(a+"-ca")
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if aBytes == nil {
// 		return shim.Error("Entity not found")
// 	}
// 	return shim.Success(aBytes)
// }
// // get my ca security
// func (t *SimpleChaincode) myCaSecurity(stub shim.ChaincodeStubInterface,args []string,name string) pb.Response {
// 	var a string    // Entities
// 	//var x int          // Transaction value
// 	var err error
	
// 	if len(args) !=1 {
// 		return shim.Error("Incorrect number of hjuh args. Expecting 1")
// 	}
// 	a = name//args[0]
// 	aSub := args[0]
// 	// Get the state from the ledger
// 	aBytes, err := stub.GetState(a+"-ca-"+aSub)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if aBytes == nil {
// 		return shim.Error("you are not subscribed to this security")
// 	}
	
// 	return shim.Success(aBytes)
// }
// // issue CA
// func (t *SimpleChaincode) issueCa(stub shim.ChaincodeStubInterface,args []string) pb.Response {
// 	var aSub string // security
// 	//var x int          // Transaction value
// 	var err error
		
// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of hjuh args. Expecting 2")
// 	}
// 	aSub = args[0]
// 	caData := args[1]
	
// 	securityBytes, err := stub.GetState("security")
// 	securityString :=string(securityBytes)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if(!strings.Contains(securityString,aSub+",")){
// 		return shim.Error("Invalid Security")
// 	}
// 	securitySubBytes, err := stub.GetState(aSub)
// 	securitySubString :=string(securitySubBytes)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	if securitySubBytes == nil {
// 		return shim.Error("No Subscribers")
// 	}
// 	if(len(securitySubString)<2){
// 		return shim.Error("No valid Subscribers")
// 	}
// 	result := strings.Split(securitySubString, ",")
	
// 	err2 := stub.PutState("ca-"+aSub,[]byte(caData))
// 	if err2 != nil {
// 	return shim.Error(err2.Error())
//     }
// 	 // Display all elements.
// 	i :=0;
//     for i = range result {
//     /*err2 := stub.PutState("debug-"+aSub+strconv.Itoa(i),[]byte(strconv.Itoa(i)))
// 	if err2 != nil {
// 	return shim.Error(err2.Error())
//     }*/
//         aBytes, err := stub.GetState(result[i]+"-ca")
//         if err != nil {
// 		return shim.Error(err.Error())
//         }
// 		if(aBytes == nil){
// 		err1 := stub.PutState(result[i]+"-ca",[]byte(aSub))
// 		if err1 != nil {
// 		return shim.Error(err1.Error())
//         }
// 		} else {
// 		stringAbytes := string(aBytes)
// 		err1 := stub.PutState(result[i]+"-ca",append(append([]byte(stringAbytes),","...),[]byte(aSub)...))
// 		if err1 != nil {
// 		return shim.Error(err1.Error())
//         }
// 		}
// 		err2 = stub.PutState(result[i]+"-ca-"+aSub,[]byte(caData))
// 		if err2 != nil {
// 		return shim.Error(err2.Error())
//         }
//     }
// 	return shim.Success(nil)
// }

// // deletes an entity from state
// func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	if len(args) != 1 {
// 		return pb.Response{Status:403, Message:"Incorrect number of args"}
// 	}

// 	a := args[0]

// 	// Delete the key from the state in ledger
// 	err := stub.DelState(a)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	return shim.Success(nil)
// }

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