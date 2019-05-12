package gamedata

import (
	"encoding/json"
	"fmt"
)

type TradeFulfillment struct {
	Id          int32     `json:"id"`
	Seller      int32     `json:"seller"`
	Buyer       int32     `json:"buyer"`
	Item        Equipment `json:"item"`
	SellerYield int32     `json:"sellerYield"`
	BuyerYield  int32     `json:"buyerYield"`
	MinerYield  int32     `json:"minerYield"`
}

type TradeFulfillments struct {
	Fulfillments map[int32]TradeFulfillment `json:"fulfillments"`
}

type FulfillmentList struct {
	FulfillmentList []string `json:"fulfillmentList"`
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
	jsonOut := fmt.Sprintf("{\"id\": \"%v\",\"seller\": %v,\"buyer\": %v,\"item\": %v,\"sellerYield\": %v,\"buyerYield\": %v,\"minerYield\": %v}",
		fulfillment.Id, fulfillment.Seller, fulfillment.Buyer, item, fulfillment.SellerYield, fulfillment.BuyerYield, fulfillment.MinerYield)
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
	f.SellerYield = int32(fulfillmentMap["sellerYield"].(float64))
	f.BuyerYield = int32(fulfillmentMap["buyerYield"].(float64))
	f.MinerYield = int32(fulfillmentMap["minerYield"].(float64))
	return f, nil
}

// {"tradeFulfillment": {“sellerId”: 6690,“buyerId”: 6684,“item”: {"name": "sword",“id”: 2,“owner”: 1,“description”: “This is a sword I got from a slime.”,"stats" : {"level": 1,"atk": 5,“def”: 5}},“sellerYield”: {“currency”: 1000},“buyerYield”: {“currency”: -1000},“minerYield”: {“currency”: 10}}}
func (fulfillments *TradeFulfillments) RemoveFulfillments(fulfilledTradesJson string) bool {
	//var fulfillmentList FulfillmentList
	var fulfillmentList []interface{}
	err := json.Unmarshal([]byte(fulfilledTradesJson), &fulfillmentList)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//hhh := fulfillmentList
	fmt.Printf("listo %v\n", fulfillmentList)

	//fmt.Printf("len: %v\n", len(fulfillmentList.FulfillmentList))
	//for i, fulfillment := range fulfillmentList.FulfillmentList {
	//	fmt.Printf("i, f: %v,%v\n", i, fulfillment)
	//}
	return true
}
