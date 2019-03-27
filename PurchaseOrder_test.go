
/*
Copyright IBM Corp. 2016 All Rights Reserved.
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
package main

import (
"encoding/json"
"fmt"
"reflect"
"testing"

"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
res := stub.MockInit("1", args)
if res.Status != shim.OK {
	fmt.Println("Init failed", string(res.Message))
	t.FailNow()
}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
res := stub.MockInvoke("1", args)
if res.Status != shim.OK {
	fmt.Println("Invoke", args, "failed", string(res.Message))
	t.FailNow()
}
}

func checkQuery(t *testing.T, stub *shim.MockStub, value map[string]interface{}, args [][]byte) {
res := stub.MockInvoke("1", args)

if res.Status != shim.OK {
	fmt.Println("Query failed", string(res.Message))
	t.FailNow()
}
if res.Payload == nil {
	fmt.Println("Query failed to get value")
	t.FailNow()
}

str := string(res.Payload)
var valueMap map[string]interface{}
json.Unmarshal([]byte(str), &valueMap)

fmt.Println("After query: ", string(res.Payload))
fmt.Println("=======")
valueString, _ := json.Marshal(value)

fmt.Println("Parameter:", string(valueString))

eq := reflect.DeepEqual(value, valueMap)
if eq {
	fmt.Println("They're equal.")
} else {
	fmt.Println("Query failed")
	t.FailNow()
}
fmt.Println("---------------------")
}


func TestEndToEndWorkflow(t *testing.T) {
scc := new(PurchaseOrder)
stub := shim.NewMockStub("PurchaseOrder", scc)

// Lane 1, Port to Door, test cases
checkInit(t, stub, [][]byte{[]byte("init")})
//invoking createPO
checkInvoke(t, stub, [][]byte{[]byte("createPO"), []byte("{\"POID\": \"1234\"}")})
checkInvoke(t, stub, [][]byte{[]byte("createPO"), []byte("{\"POID\": \"12345\"}")})
//invoking updateStatus
checkInvoke(t, stub, [][]byte{[]byte("updateStatus"), []byte("{\"POID\": \"1234\",\"status\": \"active\"}")})
checkInvoke(t, stub, [][]byte{[]byte("updateStatus"), []byte("{\"POID\": \"12345\",\"status\": \"active\"}")})
//invoking updateQuantity
checkInvoke(t, stub, [][]byte{[]byte("updateQuantity"), []byte("{\"POID\": \"1234\",\"quantity\": \"10\"}")})
checkInvoke(t, stub, [][]byte{[]byte("updateQuantity"), []byte("{\"POID\": \"12345\",\"quantity\": \"20\"}")})
//invoking updateproductName
checkInvoke(t, stub, [][]byte{[]byte("updateproductName"), []byte("{\"POID\": \"1234\",\"productName\": \"Product_A\"}")})
checkInvoke(t, stub, [][]byte{[]byte("updateproductName"), []byte("{\"POID\": \"12345\",\"productName\": \"Product_B\"}")})
//invoking customer
checkInvoke(t, stub, [][]byte{[]byte("updateCustomer"), []byte("{\"POID\": \"12345\",\"customer\": \"Jyoti\"}")})

//invoking deletePO
checkInvoke(t, stub, [][]byte{[]byte("deletePO"), []byte("{\"POID\": \"1234\"}")})

str := "{\"POID\":\"12345\"}"



	var valueMap map[string]interface{}
	json.Unmarshal([]byte(str), &valueMap)

	checkQuery(t, stub, valueMap, [][]byte{[]byte("getPODetails"), []byte("{\"POID\": \"12345\",\"customer\": \"Jyoti\"}")})
}



