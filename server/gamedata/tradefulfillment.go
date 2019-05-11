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
	f.Seller = int32(fulfillmentMap["seller"].(float64))
	f.Buyer = int32(fulfillmentMap["buyer"].(float64))
	// decode equipment json to equipment struct
	eqpString := fulfillmentMap["item"].(string)
	eqp, err := DecodeEquipmentFromJson(eqpString)
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

func (fulfillments *TradeFulfillments) RemoveFulfillments(fulfilledTradesJson string) {

}
