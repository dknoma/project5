package server

import (
	"fmt"
	"github.com/dknoma/project5/server/gamedata"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var MyPort int32
var MyID int32 = 0

var BLOCKCHAIN_SERVER = "http://localhost:6686"

//var REGISTER_SERVER = BLOCKCHAIN_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = BLOCKCHAIN_SERVER + "/upload"
var SELF_ADDR = "http://localhost:"

var TradeRequests gamedata.RequestCache
var PendingTradeFulfillments gamedata.TradeFulfillments
var UserList gamedata.Users

var nextUserId = int32(0)
var nextTradeRequestId = int32(0)
var nextTradeFulfillmentId = int32(0)
var ifStarted bool

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
	UserList.InitUserList()
	TradeRequests.InitTradeRequests()
	PendingTradeFulfillments.InitTradeFulfillments()
}

// Register ID, download BlockChain, start HeartBeat
// dApp has special ID 6666 for peer list.
func Start(w http.ResponseWriter, r *http.Request) {
	// After register, send heartbeat every 5-10 seconds
	if ifStarted {
		return
	}
	ifStarted = true
	SELF_ADDR = fmt.Sprintf("%v%v", SELF_ADDR, MyPort)
}

// Register to TA's server, get an ID
func Register() {
	fmt.Printf("registering...\n")
	fmt.Printf("req: %v\n", MyID)
}

// Allow users to create an "account" (just a basic user)
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	clientId := int32(nextUserId)
	// Create a new user
	UserList.AddUser(clientId)
	nextUserId++
	fmt.Fprint(w, clientId)
}

// POST body
// 	"item": json
// 	"sellerId": id
// 	"demands": json
func CreateRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// Error occurred. Param was not an integer
		returnServerError(w, "reading body")
		return
	}
	parsedBodyValue, err := url.ParseQuery(string(body)) // Parse request body into a Value
	if err != nil {
		// Error occurred. Param was not an integer
		returnServerError(w, "query parsing")
		return
	}
	fmt.Printf("body value: %v\n", parsedBodyValue) // Print out the parsed body value
	// Convert to id from the map to an int, to int32
	id, err := strconv.Atoi(parsedBodyValue["seller"][0])
	sellerId := int32(id)
	equipmentSlot, err := strconv.Atoi(parsedBodyValue["equipmentSlot"][0]) // The slot of the desired equipment in the users inventory
	//cost := parsedBodyValue["cost"][0]                                      // cost of the demand
	demandJson := parsedBodyValue["demands"][0]
	demands, err := gamedata.DecodeDemandJson(demandJson)
	if err != nil {
		// Error occurred. Param was not an integer
		returnServerError(w, "body value")
		return
	}
	fmt.Printf("demands %v\n", demands)
	// try to create the trade request
	out, success := tryCreateRequest(sellerId, equipmentSlot, demands)
	if !success {
		returnServerError(w, "trade req")
		return
	}
	// Print out trade request json to client
	fmt.Fprint(w, out)
}

func tryCreateRequest(sellerId int32, equipmentSlot int, demands gamedata.Demands) (string, bool) {
	seller, sellerExists := UserList.Users[sellerId] // actual user of the seller
	fmt.Printf("seller %v\n", seller)
	if !sellerExists || equipmentSlot > len(seller.Inventory.Equipment) || equipmentSlot < 0 {
		// Seller doesn't exist || equipment slot doesnt exist
		return "Unable to create request.", false
	}
	equipmentToSell := seller.Inventory.Equipment[equipmentSlot] // the equipment from the sellers inventory
	if gamedata.EquipmentIsEmpty(equipmentToSell) {
		// equipment is empty/doesnt exist
		return "Unable to create request.", false
	}
	//fmt.Printf("seller id: %v, equipment slot: %v, cost: %v, seller: %v, equipment: %v\n", sellerId, equipmentSlot,
	//	demands.Currency, seller, equipmentToSell)
	newTradeRequest := gamedata.TradeRequest{Id: nextTradeRequestId, Seller: sellerId, Item: equipmentToSell, Demands: demands}
	TradeRequests.AddToRequestCache(newTradeRequest)
	fmt.Printf("new request: %v\n", newTradeRequest)
	nextTradeRequestId++
	req, err := newTradeRequest.EncodeRequestToJson()
	if err != nil {
		// Unable to enccode trade request
		return "Unable to create request.", false
	}
	return req, true
}

func ViewRequest(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/") // split url paths
	reqId, err := strconv.Atoi(p[3])
	requestId := int32(reqId)
	if err != nil {
		// Error occurred. Param was not an integer
		returnServerError(w, "strconv")
		return
	}
	fmt.Printf("Trade request ID: %v", requestId)
	// take request id
	// check the db for the request id
	// get the request json
	// send that request json to the client
	tradeRequestJson, exists := TradeRequests.TradeRequests[requestId]
	if !exists {
		// Error occurred. Param was not an integer
		returnServerError(w, "trade req")
		return
	}
	fmt.Printf("tradeRequestJson: %v", tradeRequestJson)
	fmt.Fprint(w, tradeRequestJson) // Send json to the client
}

func FulfillRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: get potential json (total currency (real app would have an actual service to take care of these checks)),
	//		 from POST, grab request id from params and find the req by id
	p := strings.Split(r.URL.Path, "/") // split url paths
	id, err := strconv.Atoi(p[3])
	if err != nil {
		// Error occurred. Param was not an integer
		returnServerError(w, "strconv")
		return
	}
	fmt.Printf("Trade request ID: %v", id)
	tradeReqId := int32(id)
	tradeReq, exists := TradeRequests.TradeRequests[tradeReqId]
	if !exists {
		// Trade request doesn't exist
		returnServerError(w, "trade req")
		return
	}
	fmt.Printf("trade req: %v\n", tradeReq)
	// Read the POST body from the client
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// Error occurred. Param was not an integer
		returnServerError(w, "reading body")
		return
	}
	parsedBodyValue, err := url.ParseQuery(string(body)) // Parse request body into a Value
	if err != nil {
		// Error occurred.
		returnServerError(w, "query parsing")
		return
	}
	//fmt.Printf("body value: %v\n", parsedBodyValue)
	bId, err := strconv.Atoi(parsedBodyValue["buyer"][0]) // Get first index
	if err != nil {
		// Error occurred.
		returnServerError(w, "string conversion")
		return
	}
	buyerId := int32(bId)
	fmt.Printf("buyer id %v\n", buyerId)
	// Try to create the new fulfillment
	newFulfillment, validReqs := tryCreateFulfillment(buyerId, tradeReq)
	if !validReqs {
		// Error occurred.
		returnServerError(w, "try create fulfillment")
		return
	}
	fmt.Printf("ful: %v\n", newFulfillment)
	// Add fulfillment to pending
	PendingTradeFulfillments.Fulfillments[newFulfillment.Id] = newFulfillment
	fulJson, err := newFulfillment.EncodeFulfillmentToJson()
	fmt.Printf("all pending: %v\n", PendingTradeFulfillments.Fulfillments)
	fmt.Fprint(w, fulJson)
}

func returnServerError(w http.ResponseWriter, str string) {
	fmt.Printf("%v - error: %v - %v\n", str, http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError))
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, fmt.Sprintf("%d - %s",
		http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
}

// Check if the buyer exists, check if they have enough currency
func tryCreateFulfillment(buyerId int32, tradeReq gamedata.TradeRequest) (gamedata.TradeFulfillment, bool) {
	buyerExists := UserList.UserExists(buyerId)
	if !buyerExists {
		return gamedata.TradeFulfillment{}, false
	}
	bounty := tradeReq.Demands.Currency
	hasEnoughCurrency := UserList.HasEnoughCurrency(buyerId, bounty)
	if !hasEnoughCurrency {
		return gamedata.TradeFulfillment{}, false
	}
	// Create the fulfillment
	minerYield := calculateMinerYield(bounty)
	newFulfillment := gamedata.NewFulfillment(nextTradeFulfillmentId, tradeReq.Id, tradeReq.Seller, buyerId, tradeReq.Item,
		bounty-minerYield, -bounty, minerYield)
	return newFulfillment, true
}

// 2% yield
func calculateMinerYield(bounty float64) float64 {
	return bounty * 0.02
}

// Miners get an array of pending transactions to put into their mpt to mine
// how to insert these transactions into the mpt?
func GetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	//PendingTradeFulfillments
	pendingJson, err := PendingTradeFulfillments.EncodeFulfillmentsToJson()
	if err != nil {
		// Error occurred.
		returnServerError(w, "fulfillment json")
		return
	}
	fmt.Printf("pending json: %v\n", pendingJson)
	fmt.Fprintf(w, pendingJson)
}

// Called by miner when successfully mines a block, tho miner could fake this data
func UpdatePendingTransactions(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		// Error occurred. Param was not an integer
		fmt.Printf("reading body: %v - %v\n", http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	parsedBodyValue, err := url.ParseQuery(string(body)) // Parse request body into a Value
	if err != nil {
		// Error occurred. Param was not an integer
		fmt.Printf("query parsing - error: %v | %v - %v\n", err, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	parsedData := parsedBodyValue["data"][0] // Get first index
	fmt.Printf("parsedData body value: %v\n", parsedData)
	//var fulfillmentsToUpdate gamedata.TradeFulfillments
	fulfillmentsToUpdate, err := gamedata.DecodeFulfillmentJsonArrayToInterface(parsedData)
	//err = json.Unmarshal([]byte(parsedData), &fulfillmentsToUpdate)
	if err != nil {
		// Error occurred. Param was not an integer
		fmt.Printf("query parsing - error: %v | %v - %v\n", err, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}

	fmt.Printf("fff %v\n", fulfillmentsToUpdate)
	for id := range fulfillmentsToUpdate.Fulfillments {
		//delete(PendingTradeFulfillments.Fulfillments, id)
		updateTradeDatabase(id)
	}
}

// Update the trade databases
func updateTradeDatabase(id int32) {
	transferFulfillmentData(id)
	fmt.Printf("req: before remove: %v\n", TradeRequests.TradeRequests)
	fmt.Printf("ful: before remove: %v\n", PendingTradeFulfillments.Fulfillments)
	ful := PendingTradeFulfillments.Fulfillments[id]
	TradeRequests.RemoveFromRequestCache(ful.RequestId)
	PendingTradeFulfillments.RemoveFulfillment(id)
	fmt.Printf("req: after remove: %v\n", TradeRequests.TradeRequests)
	fmt.Printf("ful: after remove: %v\n", PendingTradeFulfillments.Fulfillments)
}

func transferFulfillmentData(id int32) {
	fulfillment := PendingTradeFulfillments.GetFulfillment(id)
	fmt.Printf("fuasfl %v\n", fulfillment)
	UserList.TradeItem(fulfillment.Seller, fulfillment.Buyer, fulfillment.Item, fulfillment.SellerYield, fulfillment.BuyerYield)
	fmt.Printf("users %v\n", UserList.Users)
}
