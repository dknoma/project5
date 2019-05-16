package gamedata

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type TradeFulfillment struct {
	Id          int32     `json:"id"`
	RequestId   int32     `json:"requestId"`
	Seller      int32     `json:"seller"`
	Buyer       int32     `json:"buyer"`
	Item        Equipment `json:"item"`
	SellerYield float64   `json:"sellerYield"`
	BuyerYield  float64   `json:"buyerYield"`
	MinerYield  float64   `json:"minerYield"`
}

type TradeFulfillments struct {
	Fulfillments map[int32]TradeFulfillment `json:"fulfillments"`
	mux          sync.Mutex
}

func (fulfillments *TradeFulfillments) InitTradeFulfillments() {
	fulfillments.Fulfillments = make(map[int32]TradeFulfillment)
}

func NewFulfillment(id, requestId, sellerId, buyerId int32, item Equipment, sellerYield, buyerYield, minerYield float64) TradeFulfillment {
	newFulfillment := TradeFulfillment{Id: id, RequestId: requestId, Seller: sellerId, Buyer: buyerId,
		Item: item, SellerYield: sellerYield, BuyerYield: buyerYield, MinerYield: minerYield}
	return newFulfillment
}

func (fulfillments *TradeFulfillments) AddFulfillment(fulfillment TradeFulfillment) {
	fulfillments.mux.Lock()
	defer fulfillments.mux.Unlock()
	fulfillments.Fulfillments[fulfillment.Id] = fulfillment
}

func (fulfillments *TradeFulfillments) RemoveFulfillment(id int32) {
	fulfillments.mux.Lock()
	defer fulfillments.mux.Unlock()
	delete(fulfillments.Fulfillments, id)
}

func (fulfillments *TradeFulfillments) GetFulfillment(id int32) TradeFulfillment {
	fulfillments.mux.Lock()
	defer fulfillments.mux.Unlock()
	return fulfillments.Fulfillments[id]
}

func (fulfillment *TradeFulfillment) EncodeFulfillmentToJson() (string, error) {
	jsonBytes, err := json.Marshal(fulfillment)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	mptJson := string(jsonBytes)
	mptJson = mptJson[1 : len(mptJson)-1]
	// Get the requests item and encode it into JSON format
	item, err := fulfillment.Item.EncodeEquipmentToJson()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	jsonOut := fmt.Sprintf("{\"id\": %v,\"requestId\":%v,\"sellerId\": %v,\"buyerId\": %v,\"item\": %v,\"sellerYield\": %v,\"buyerYield\": %v,\"minerYield\": %v}",
		fulfillment.Id, fulfillment.RequestId, fulfillment.Seller, fulfillment.Buyer, item, fulfillment.SellerYield, fulfillment.BuyerYield, fulfillment.MinerYield)
	isValid := json.Valid([]byte(jsonOut))
	if !isValid {
		fmt.Println(err.Error())
		return "", err
	}
	return jsonOut, nil
}

func DecodeFulfillmentFromJSON(jsonString string) (TradeFulfillment, error) {
	jsonBytes := []byte(jsonString)
	// Unmarshal the json bytes into a new key:value map
	var fulfillmentMap map[string]interface{}
	err := json.Unmarshal(jsonBytes, &fulfillmentMap)
	if err != nil {
		fmt.Println(err.Error())
		return TradeFulfillment{}, err
	}
	// Create new TradeFulfillment to insert unmarshalled values into
	var f TradeFulfillment
	f.Id = int32(fulfillmentMap["id"].(float64))
	f.RequestId = int32(fulfillmentMap["requestId"].(float64))
	f.Seller = int32(fulfillmentMap["sellerId"].(float64))
	f.Buyer = int32(fulfillmentMap["buyerId"].(float64))

	// decode equipment json to equipment struct
	//eqpString := fulfillmentMap["item"].(string)
	eqpMap := fulfillmentMap["item"].(map[string]interface{})
	eqp, err := DecodeEquipmentFromMap(eqpMap)
	//eqp, err := DecodeEquipmentFromJson(eqpString)
	if err != nil {
		fmt.Println(err.Error())
		return TradeFulfillment{}, err
	}

	f.Item = eqp
	f.SellerYield = fulfillmentMap["sellerYield"].(float64)
	f.BuyerYield = fulfillmentMap["buyerYield"].(float64)
	f.MinerYield = fulfillmentMap["minerYield"].(float64)
	return f, nil
}

// {"tradeFulfillment": {“sellerId”: 6690,“buyerId”: 6684,“item”: {"name": "sword",“id”: 2,“owner”: 1,“description”: “This is a sword I got from a slime.”,"stats" : {"level": 1,"atk": 5,“def”: 5}},“sellerYield”: {“currency”: 1000},“buyerYield”: {“currency”: -1000},“minerYield”: {“currency”: 10}}}
//func (fulfillments *TradeFulfillments) TryRemoveFulfillments(fulfilledTradesJson string) bool {
//	fulfillments.mux.Lock()
//	defer fulfillments.mux.Unlock()
//	// Do decode
//	fulfillmentsToRemove, success := DecodeFulfillmentJsonArrayToInterface(fulfilledTradesJson)
//	if !success {
//		return false
//	}
//	fmt.Printf("b4 reduced fulfillments %v\n", fulfillments.Fulfillments)
//	// try delete
//	for id := range fulfillmentsToRemove.Fulfillments {
//		delete(fulfillments.Fulfillments, id)
//	}
//	fmt.Printf("reduced fulfillments %v\n", fulfillments.Fulfillments)
//	return true
//}

// When miner successfully creates a block, need to get the ids of fulfillments to remove
func GetFulfillmentIdsFromJson(fulfilledTradesJson string) []int32 {
	// TODO: get ids
	//  	 get request ids
	//  	 remove those requests
	//  	 remove from fulfillment cache
	fulfillmentsToRemove, err := DecodeFulfillmentJsonArrayToInterface(fulfilledTradesJson)
	if err != nil {
		return []int32{}
	}
	var ids []int32
	for id := range fulfillmentsToRemove.Fulfillments {
		ids = append(ids, id)
	}
	return ids
}

// Decode json string -> []interface{} -> TradeFulfillments
func DecodeFulfillmentJsonArrayToInterface(fulfilledTradesJson string) (TradeFulfillments, error) {
	var fulfillmentList []interface{}
	err := json.Unmarshal([]byte(fulfilledTradesJson), &fulfillmentList)
	if err != nil {
		fmt.Println(err.Error())
		return TradeFulfillments{}, err
	}
	// Convert interface array into fulfillment map
	var fulfillments TradeFulfillments
	fulfillments.InitTradeFulfillments()
	for _, fulfillment := range fulfillmentList {
		fulfillment := DecodeInterfaceToFulfillment(fulfillment.(map[string]interface{}))
		//fmt.Printf("fulfillment: %v,%v\n",i,fulfillment)
		fulfillments.Fulfillments[fulfillment.Id] = fulfillment
	}
	fmt.Printf("mapo %v\n", fulfillments)
	return fulfillments, nil
}

// Decode map[string]interface{} to TradeFulfillment
func DecodeInterfaceToFulfillment(fromMap map[string]interface{}) TradeFulfillment {
	var ful TradeFulfillment
	ful.Id = int32(fromMap["id"].(float64))
	ful.RequestId = int32(fromMap["requestId"].(float64))
	ful.Seller = int32(fromMap["sellerId"].(float64))
	ful.Buyer = int32(fromMap["buyerId"].(float64))
	eqpMap := fromMap["item"].(map[string]interface{})

	var e Equipment
	e.Name = eqpMap["name"].(string)
	e.Id = int32(eqpMap["id"].(float64))
	e.Owner = int32(eqpMap["owner"].(float64))
	e.Description = eqpMap["description"].(string)

	ful.Item = e
	ful.SellerYield = fromMap["sellerYield"].(float64)
	ful.BuyerYield = fromMap["buyerYield"].(float64)
	ful.MinerYield = fromMap["minerYield"].(float64)

	return ful
}

func (fulfillments *TradeFulfillments) EncodeFulfillmentsToJson() (string, error) {
	var jsonOut string
	sb := strings.Builder{}
	sb.WriteString("[")
	for _, v := range fulfillments.Fulfillments {
		jsonOut, err := v.EncodeFulfillmentToJson()
		if err != nil {
			return "", err
		}
		sb.WriteString(jsonOut)
		sb.WriteString(",")
	}
	jsonOut = sb.String()
	if len(jsonOut) > 2 {
		jsonOut = jsonOut[:len(jsonOut)-1]
	}
	jsonOut += "]"
	return jsonOut, nil
}
