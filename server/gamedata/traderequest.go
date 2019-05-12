package gamedata

import (
	"encoding/json"
	"fmt"
	"sync"
)

//- trade requests have their own id and the sellers id
//	- mpt = [tx id]trade request json
//		- option1: trade request/trade fulfillment share tx ids
//		- option2: req+id (req1, req12, etc…), and ful+id (ful1, ful12, etc…)
//			- the requests and fulfillments themselves only have int32 ids
//			- string id is purely for for mpt insertion
type TradeRequest struct {
	Id      int32     `json:"id"`
	Seller  int32     `json:"seller"`
	Item    Equipment `json:"item"`
	Demands Demands   `json:"demands"`
}

type Demands struct {
	Currency int32 `json:"currency"`
	// TODO: optionally, could have an []Items to allow users to request money + item(s)
	//		 NO CANCELING TRADE REQUESTS AT THIS MOMENT. This will be one of the later features if time allows
}

// Realistically some sort of db/SQL would store some gamedata maybe?
// Required to know where a transaction is located rather than storing the transaction itself
type RequestCache struct {
	TradeRequests map[int32]TradeRequest `json:"tradeRequests"` //	[trade request id]request block info (height at which its stored)
	mux           sync.Mutex
}

func (t *TradeRequest) EncodeRequestToJson() (string, error) {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	mptJson := string(jsonBytes)
	mptJson = mptJson[1 : len(mptJson)-1]
	// Get the requests item and encode it into JSON format
	item, err := t.Item.EncodeEquipmentToJson()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	jsonOut := fmt.Sprintf("{\"id\": \"%v\",\"seller\": %v,\"item\": %v,\"demands\": {\"currency\": %v}}",
		t.Id, t.Seller, item, t.Demands.Currency)
	isValid := json.Valid([]byte(jsonOut))
	if !isValid {
		fmt.Println(err.Error())
		return "", err
	}
	return jsonOut, nil
}

func (r *RequestCache) AddToRequestCache(req TradeRequest) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.TradeRequests[req.Id] = req
}
