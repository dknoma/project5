package server

import (
	"fmt"
	"github.com/dknoma/cs686-blockchain-p3-dknoma/p3/data"
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

// TODO: Maybe instead of keeping a peer list, keep a list of listeners
//		 choose 32 random listeners, send tx to them which then forward those tx to other miners like heartbeats.
//		 	- an improvement on this system is to utilize a reputation system, send to random 32 of the top 100 reputation
//			  miners to discourage byzantine/selfish behaviour
//		 or just send to everyone in the list, tho there could be a LOT of miners, bloating the system (tho there could also
//			be a lot of players so it might not be an issue for this "game")
//		OR
//			miners probe the dApp for possible transactions, dApp provides the latest tx that have not been validated
//				when successfully create a block, broadcast to the dApp that these transactions have been collected,
//				dApp then removes them from their list

var BlockchainPeers data.PeerList
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
	//OutID = MyPort
	//SELF_ADDR = fmt.Sprintf("%v%v", SELF_ADDR, MyPort)
	//fmt.Printf("INIT: %v, %v\n", SELF_ADDR, TA_SERVER)
	//SBC = data.NewBlockChain() // Init synch blockchain here
	UserList.InitUserList()
	TradeRequests.InitTradeRequests()
	PendingTradeFulfillments.InitTradeFulfillments()
	//mpt := p1.MerklePatriciaTrie{}
	//mpt.NewTree()
	//block := SBC.GenBlock(mpt, "")
	//SBC.Insert(block)
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
	//Register()
	//StartHeartBeat()
	fmt.Fprint(w, BlockchainPeers.GetSelfId())
}

// Register to TA's server, get an ID
func Register() {
	fmt.Printf("registering...\n")
	//resp, err := http.Get(REGISTER_SERVER) // GET to server
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
	////var selfId int32
	//body := resp.Body // Get the response body
	//respData, err := ioutil.ReadAll(body)
	//if err != nil {
	//	fmt.Printf("decode resp err: %v\n", err)
	//	return
	//}
	//id, err := strconv.Atoi(string(respData))
	//if err != nil {
	//	fmt.Printf("decode resp err: %v\n", err)
	//	return
	//}
	//selfId := int32(id)
	fmt.Printf("req: %v\n", MyID)
	BlockchainPeers.Register(MyID)
	BlockchainPeers = data.NewPeerList(MyID, 32)

	//if SELF_ADDR != BLOCKCHAIN_SERVER {
	//	fmt.Printf("can download\n")
	//	Download()
	//}
}

//// Download blockchain from TA server
//func Download() {
//	fmt.Printf("downloading...\n")
//	//resp, err := http.Get(BC_DOWNLOAD_SERVER) // GET to server
//	newHeartBeat := data.NewHeartBeatData(false, BlockchainPeers.GetSelfId(), "", "", SELF_ADDR)
//	jsonHBBytes, err := json.Marshal(newHeartBeat)
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//		return
//	}
//	form := url.Values{}
//	form.Set("heartBeat", string(jsonHBBytes))
//	resp, err := http.PostForm(BC_DOWNLOAD_SERVER, form) // POST to server
//	//req, err := http.NewRequest("POST", BC_DOWNLOAD_SERVER, strings.NewReader(form.Encode()))
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//		return
//	}
//	defer resp.Body.Close()
//	body := resp.Body
//	if err != nil {
//		fmt.Printf("decode resp err: %v\n", err)
//		return
//	}
//
//	respData, err := ioutil.ReadAll(body)
//	if err != nil {
//		fmt.Printf("decode resp err: %v\n", err)
//		return
//	}
//	respString := string(respData)
//	fmt.Printf("DOWNLOAD: decode resp: %v\n", respString)
//	SBC.UpdateEntireBlockChain(respString)
//}

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
	// TODO: Still allows for multiple requests for the same item.
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

//// Verify the list of pending transactions
//func VerifyTransaction(txJsonArray string) {
//	//PendingTradeFulfillments
//}

// Called by miner when successfully mines a block, tho miner could fake this data
// TODO: takes whatever ids it can get from the json, remove those pending transactions; there is currently no authentication
//		 to make sure that miners aren't lying.
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

func updateTradeDatabase(id int32) {
	fmt.Printf("req: before remove: %v\n", TradeRequests.TradeRequests)
	fmt.Printf("ful: before remove: %v\n", PendingTradeFulfillments.Fulfillments)
	ful := PendingTradeFulfillments.Fulfillments[id]
	delete(TradeRequests.TradeRequests, ful.RequestId)
	delete(PendingTradeFulfillments.Fulfillments, id)
	fmt.Printf("req: after remove: %v\n", TradeRequests.TradeRequests)
	fmt.Printf("ful: after remove: %v\n", PendingTradeFulfillments.Fulfillments)
}

// Ask another server to return a block of certain height and hash
// Gets called by HeartBeatReceive
//func AskForBlock(height int32, hash string) (bool, error) {
//	for addr := range BlockchainPeers.Copy() {
//		if addr == SELF_ADDR {
//			continue
//		} // Dont send to self
//		// "/block/{height}/{hash}"
//		getBlockURL := fmt.Sprintf("%v/block/%v/%v", addr, height, hash)
//		resp, err := http.Get(getBlockURL)
//		if err != nil {
//			fmt.Printf("Heartbeat send error: %v\n", err)
//			return false, err
//		}
//		statusCode := resp.StatusCode
//		switch statusCode {
//		case http.StatusNoContent:
//			// No block, go on to next peer
//			continue
//		case http.StatusOK:
//			// Successfully got block
//			body, err := ioutil.ReadAll(resp.Body)
//			if err != nil {
//				fmt.Printf("Heartbeat send error: %v\n", err)
//				return false, err
//			}
//			incomingParentBlock, err := p2.DecodeFromJSON(string(body))
//
//			// Check if the parent of this parent block exists
//			parentExists := SBC.CheckParentHash(incomingParentBlock)
//			if parentExists {
//				SBC.Insert(incomingParentBlock)
//				return true, nil
//			} else {
//				// Grandparent doesnt exist, must try to grab that as well
//				exists, err := AskForBlock(incomingParentBlock.Header.Height, incomingParentBlock.Header.ParentHash)
//				if err != nil {
//					fmt.Printf("AskForBlock error: %v\n", err)
//					return false, err
//				}
//				if exists {
//					SBC.Insert(incomingParentBlock)
//					return true, nil
//				}
//				return false, err
//			}
//		case http.StatusInternalServerError:
//			panic(fmt.Sprintf("%v - %v", http.StatusInternalServerError,
//				http.StatusText(http.StatusInternalServerError)))
//			resp.Body.Close()
//			return false, err
//		}
//		resp.Body.Close()
//	}
//	//panic(fmt.Sprintf("%v - %v", http.StatusInternalServerError,
//	//	http.StatusText(http.StatusInternalServerError)))
//	return false, nil
//}

// Received a heartbeat
//func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
//	body, err := ioutil.ReadAll(r.Body)
//	defer r.Body.Close()
//	if err != nil {
//		// Error occurred. Param was not an integer
//		fmt.Printf("reading body: %v - %v\n", http.StatusInternalServerError,
//			http.StatusText(http.StatusInternalServerError))
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprint(w, fmt.Sprintf("%d - %s",
//			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
//		return
//	}
//	parsedBodyValue, err := url.ParseQuery(string(body)) // Parse request body into a Value
//	if err != nil {
//		// Error occurred. Param was not an integer
//		fmt.Printf("query parsing - error: %v | %v - %v\n", err, http.StatusInternalServerError,
//			http.StatusText(http.StatusInternalServerError))
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprint(w, fmt.Sprintf("%d - %s",
//			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
//		return
//	}
//	//fmt.Printf("body value: %v\n", parsedBodyValue)
//	parsedData := parsedBodyValue["gamedata"][0] // Get first index
//	var newHeartBeatData data.HeartBeatData
//	err = json.Unmarshal([]byte(parsedData), &newHeartBeatData)
//	if err != nil {
//		// Error occurred. Param was not an integer
//		fmt.Printf("unmarshal heartbeat - error: %v | %v - %v\n", err, http.StatusInternalServerError,
//			http.StatusText(http.StatusInternalServerError))
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprint(w, fmt.Sprintf("%d - %s",
//			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
//		return
//	}
//	if newHeartBeatData.Addr == SELF_ADDR {
//		fmt.Printf("HB was from self. ignoring...\n")
//		return
//	}
//	// Add the remote nodes address and id to your peermap
//	//fmt.Printf("incoming hb gamedata: %v, %v\n", newHeartBeatData.Addr, newHeartBeatData.Id)
//	BlockchainPeers.Add(newHeartBeatData.Addr, newHeartBeatData.Id)
//	// Add this nodes peermap into your own peermap (excluding your own address)
//	newPeerMapJson := newHeartBeatData.PeerMapJson
//	BlockchainPeers.InjectPeerMapJson(newPeerMapJson, SELF_ADDR)
//
//	// TODO: Do we want the dApp to be able to check validity of block creation? No chain is stored in the dApp
//
//	// Check if the block in the heartbeat is a new block
//	//if newHeartBeatData.IfNewBlock {
//	//hbBlock, err := p2.DecodeFromJSON(newHeartBeatData.BlockJson)
//	//if err != nil {
//	//	// Error occurred. Param was not an integer
//	//	fmt.Printf("decodeing hb block - error: %v | %v - %v\n", err, http.StatusInternalServerError,
//	//		http.StatusText(http.StatusInternalServerError))
//	//	w.WriteHeader(http.StatusInternalServerError)
//	//	fmt.Fprint(w, fmt.Sprintf("%d - %s",
//	//		http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
//	//	return
//	//}
//
//	// Check if sender actually solved hash puzzle; then check parents
//	//hashStr := hbBlock.Header.ParentHash + hbBlock.Header.Nonce + hbBlock.Value.Root
//	//fmt.Printf("\thbBlock.Header.ParentHash: %v\n", hbBlock.Header.ParentHash)
//	//fmt.Printf("\tNonce: %v\n", hbBlock.Header.Nonce)
//	//fmt.Printf("\tRoot: %v\n", hbBlock.Value.Root)
//	//sum := sha3.Sum256([]byte(hashStr))
//	//encodedStr := hex.EncodeToString(sum[:])
//	////fmt.Printf("\tencoded str: %v\n", encodedStr)
//	//validChars := 0
//
//	//hasSolved := false
//	//for i := 0; i < len(encodedStr); i++ { // break out of loop when reach max number to check
//	//	if validChars >= MAX_NONCE {
//	//		fmt.Println("Found valid nonce count!")
//	//		hasSolved = true
//	//		break
//	//	}
//	//	if string(encodedStr[i]) == "0" {
//	//		validChars++
//	//	} else { // Found non-zero before reaching end
//	//		hasSolved = false
//	//		break
//	//	}
//	//}
//	//hasSolved = validChars == MAX_NONCE
//	//fmt.Printf("sender has solved: %v\n", hasSolved)
//	//if hasSolved { // If solved, check for parent blocks then insert into chain
//	//	parentExists := SBC.CheckParentHash(hbBlock)
//	//	// If parent doesnt exist, download it before inserting new block
//	//	if !parentExists {
//	//		AskForBlock(hbBlock.Header.Height, hbBlock.Header.Hash)
//	//	}
//	//	SBC.Insert(hbBlock)
//	//}
//	//}
//
//	// If heartbeat has hops, forward to peer list
//	newHeartBeatDataHops := newHeartBeatData.Hops
//	if newHeartBeatDataHops > 0 {
//		newHeartBeatData.Hops = newHeartBeatDataHops - 1
//		//newHeartBeatData.Addr = SELF_ADDR
//		newHeartBeatData.Id = BlockchainPeers.GetSelfId()
//		ForwardHeartBeat(newHeartBeatData)
//	}
//	gotem := "gotem"
//	fmt.Fprint(w, gotem)
//}
//
//// Forward the received heartbeat to everyone on peerlist
//func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
//	remainingHops := heartBeatData.Hops
//	fmt.Printf("remaining hops: %v\n", remainingHops)
//	// Return if no more hops
//	if remainingHops == 0 {
//		return
//	}
//	// heartBeatData.Hops = remainingHops - 1
//	jsonHBBytes, err := json.Marshal(heartBeatData)
//	if err != nil {
//		fmt.Printf("err: %v\n", err)
//		return
//	}
//	jsonHB := string(jsonHBBytes)
//	//fmt.Printf("json heartbeat string: %v\n", jsonHB)
//	for addr := range BlockchainPeers.PeerMap {
//		if addr == SELF_ADDR {
//			fmt.Println("Forwarding: found own address???")
//			continue
//		} // Dont send to self
//		postData := url.Values{}
//		postData.Set("gamedata", jsonHB)
//		resp, err := http.PostForm(addr+"/heartbeat/receive", postData) // POST to server
//		if err != nil {
//			fmt.Printf("Heartbeat send error: %v\n", err)
//			return
//		}
//		resp.Body.Close()
//	}
//}

// This is a special node that only listens in on blockchain changes from its peers/the chains miners. This
// node does NOT mine/solve nonce. It only needs to know what blocks there are and to check the transactions
//func StartHeartBeat() {
//	randomRange := rand.Intn(10-5) + 5
//	ticker := time.NewTicker(time.Duration(randomRange) * time.Second)
//	go func() {
//		for t := range ticker.C {
//			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
//			fmt.Println("Sending heartbeat...")
//
//			pmJson, err := BlockchainPeers.PeerMapToJson()
//			if err != nil {
//				fmt.Printf("Error: %v\n", err)
//				return
//			}
//			preparedData := data.PrepareHeartBeatData(&data.SyncBlockChain{}, BlockchainPeers.GetSelfId(), pmJson, SELF_ADDR)
//			preparedJsonBytes, err := json.Marshal(preparedData)
//			if err != nil {
//				fmt.Printf("Error: %v\n", err)
//				return
//			}
//			preparedJson := string(preparedJsonBytes)
//			BlockchainPeers.Rebalance() // Call rebalance to check if peerlist needs to be rebalanced before sending heartbeat
//			for addr := range BlockchainPeers.PeerMap {
//				if addr == SELF_ADDR {
//					fmt.Println("\t\tFound own address.")
//					continue
//				} // Dont send to self
//				postData := url.Values{}
//				postData.Set("gamedata", preparedJson)
//				resp, err := http.PostForm(addr+"/heartbeat/receive", postData) // POST to server
//				if err != nil {
//					fmt.Printf("Heartbeat send rror: %v\n", err)
//					return
//				}
//				resp.Body.Close()
//			}
//
//			randomRange = rand.Intn(10-5) + 5
//			ticker = time.NewTicker(time.Duration(randomRange) * time.Second)
//		}
//	}()
//}
