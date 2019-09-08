package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
  "bytes"
  "encoding/json"
  "fmt"
  "strconv"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  sc "github.com/hyperledger/fabric/protos/peer"
)
import "time"

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the data structure
type SensorData struct {
  SensorId   string `json:"sensorid"`
  Date  string `json:"date"`
  Temperature string `json:"temp"`
}

/*
 * The Init method is called when the Smart Contract "temp_monitor" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
  return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "temp_monitor"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

  // Retrieve the requested Smart Contract function and arguments
  function, args := APIstub.GetFunctionAndParameters()
  // Route to the appropriate handler function to interact with the ledger appropriately
  if function == "querySensor" {
    return s.querySensor(APIstub, args)
  } else if function == "initLedger" {
    return s.initLedger(APIstub)
  } else if function == "createRecord" {
    return s.createRecord(APIstub, args)
  } else if function == "queryAllSensors" {
    return s.queryAllSensors(APIstub)
  } else if function == "queryHistoryOfSensor" {
    return s.queryHistoryOfSensor(APIstub, args)
  }

  return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) querySensor(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  sensorAsBytes, _ := APIstub.GetState(args[0])
  return shim.Success(sensorAsBytes)
}

func (s *SmartContract) queryHistoryOfSensor(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  resultsIterator, err := APIstub.GetHistoryForKey(args[0])
  if err != nil {
    return shim.Error(err.Error())
  }
  defer resultsIterator.Close()

  // buffer is a JSON array containing QueryResults
  var buffer bytes.Buffer
  buffer.WriteString("[")

  bArrayMemberAlreadyWritten := false
  for resultsIterator.HasNext() {
    queryResponse, err := resultsIterator.Next()
    if err != nil {
      return shim.Error(err.Error())
    }
    // Add a comma before array members, suppress it for the first array member
    if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
    }
    buffer.WriteString("{\"Key\":")
    buffer.WriteString("\"")
    buffer.WriteString(queryResponse.TxId)
    buffer.WriteString("\"")

    buffer.WriteString(", \"Record\":")
    // Record is a JSON object, so we write as-is
    buffer.WriteString(string(queryResponse.Value))
    buffer.WriteString("}")
    bArrayMemberAlreadyWritten = true
  }
  buffer.WriteString("]")

  fmt.Printf("- queryHistoryOfSensor:\n%s\n", buffer.String())

  return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

  return shim.Success(nil)
}

func (s *SmartContract) createRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

  if len(args) != 4 {
    return shim.Error("Incorrect number of arguments. Expecting 4")
  }

  var sensor = SensorData{SensorId: args[1], Date: args[2], Temperature: args[3]}

  sensorAsBytes, _ := json.Marshal(sensor)
  APIstub.PutState(args[0], sensorAsBytes)

  return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

  // Create a new Smart Contract
  err := shim.Start(new(SmartContract))
  if err != nil {
    fmt.Printf("Error creating new Smart Contract: %s", err)
  }
}
