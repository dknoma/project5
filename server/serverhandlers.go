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

var nextUserId = 0
var ifStarted bool

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
	//OutID = MyPort
	//SELF_ADDR = fmt.Sprintf("%v%v", SELF_ADDR, MyPort)
	//fmt.Printf("INIT: %v, %v\n", SELF_ADDR, TA_SERVER)
	//SBC = data.NewBlockChain() // Init synch blockchain here
	UserList.InitUserList()
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
	Register()
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
	// Create a new use
	newUser := gamedata.User{Id: clientId, Equipment: []gamedata.Equipment{}, Currency: 10000}
	newUser.GenerateEquipment()
	// Add new user to user list
	UserList.Users[clientId] = newUser
	nextUserId++
	fmt.Fprint(w, clientId)
}

// POST body
// 	"item": json
// 	"sellerId": id
// 	"demands": json
func CreateRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: get POST body of item json, seller id, (if actual app would have database w/ user ids, etc...)
	//		 as well as the demand json (desired currency)
	//		 How to store tx in MPT? Function in main bc nodes that update the MPT to use for a miner's block
	//		 Make it so MPT isn't randomly generated, but instead contains the gamedata from requests and fulfillments
	//		 	This might make it so whenever the MPT is changed, the nonce must start over. Ensures that they are in a block
	//				if NEED to ensure that all trade requests are put up in the marketplace
	//			OR functions that just updates the mpt to use
	//				may not be as reliable
	//fmt.Printf("Trade request ID: %v", id)
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
	fmt.Printf("body value: %v\n", parsedBodyValue) // Print out the parsed body value
	// TODO: Store this trade request in a trade request database: this allows a non-existent theoretical frontend
	//		 to actuall display this trades for players to actually see and interact with. Unless there is an efficient
	//		 way to store

	// Convert to id from the map to an int, to int32
	id, err := strconv.Atoi(parsedBodyValue["seller"][0])
	sellerId := int32(id)
	equipmentSlot, err := strconv.Atoi(parsedBodyValue["equipmentSlot"][0]) // The slot of the desired equipment in the users inventory
	cost := parsedBodyValue["cost"][0]                                      // cost of the demand

	seller, sellerExists := UserList.Users[sellerId] // actual user of the seller
	if !sellerExists {
		// Seller doesn't exist
		fmt.Printf("%v - %v\n", http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	equipmentToSell := seller.Equipment[equipmentSlot] // the equipment from the sellers inventory

	fmt.Printf("seller id: %v, equipment slot: %v, cost: %v, seller: %v, equipment: %v\n", sellerId, equipmentSlot,
		cost, seller, equipmentToSell)

	//var newTradeRequest gamedata.TradeRequest

	// seller id
	// equipment slot (slot of the equipment in the users account)
	//		get the equipment json from this slot in the users inventory
	// demand cost
	// 		verify if seller id exists
	//		verify if equipment actually exists in the player's inventory
	//		verify valid demand
	//		create request and store into db
	// Might make sense to store the requests OFF chain, and ONLY fulfillments ON chain
	//		A game based off of blockchain where miners are also players and whatnot, having everything on chain would make more sense
}

func ViewRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: IF STORING REQUESTS IN CHAIN: To view a request it MUST be in the canonical chain. Must make a call to GetCanonical and check if the tx
	//			exists in that chain. Probably have some sort of cache to store tx requestId to height (private bc)
	//		 ELSE
	//			Storing in off chain db that just stores requests in order to show them in the front end
	//
	p := strings.Split(r.URL.Path, "/") // split url paths
	reqId, err := strconv.Atoi(p[2])
	requestId := int32(reqId)
	if err != nil {
		// Error occurred. Param was not an integer
		fmt.Printf("%v - %v\n", http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
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
		fmt.Printf("%v - %v\n", http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	fmt.Printf("tradeRequestJson: %v", tradeRequestJson)
	fmt.Fprint(w, tradeRequestJson) // Send json to the client
}

func FulfillRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: get potential json (total currency (real app would have an actual service to take care of these checks)),
	//		 from POST, grab request id from params and find the req by id
	p := strings.Split(r.URL.Path, "/") // split url paths
	id, err := strconv.Atoi(p[2])
	if err != nil {
		// Error occurred. Param was not an integer
		fmt.Printf("%v - %v\n", http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	fmt.Printf("Trade request ID: %v", id)

	tradeReq, exists := TradeRequests.TradeRequests[int32(id)]
	if !exists {
		// Trade request doesn't exist
		fmt.Printf("%v - %v\n", http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}

	fmt.Printf("trade req: %v\n", tradeReq)

	// Read the POST body from the client
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
		// Error occurred.
		fmt.Printf("query parsing - error: %v | %v - %v\n", err, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, fmt.Sprintf("%d - %s",
			http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	//fmt.Printf("body value: %v\n", parsedBodyValue)
	parsedData := parsedBodyValue["data"][0] // Get first index
	fmt.Printf("stuff %v\n", parsedData)

	// get buyer id
	//		verify that the buyer id is valid
	//
	//
}

//TODO: In order to check the blockchain for trades, must access the blockchain somehow. In order to do this
//		must check the chain from the peerlist. Maybe this server is a listener on the blockchain network BUT
//		does NOT mine.

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
