package server

import (
	"encoding/json"
	"fmt"
	"github.com/dknoma/cs686-blockchain-p3-dknoma/p1"
	"github.com/dknoma/cs686-blockchain-p3-dknoma/p3/data"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//import (
//	"../p1"
//	"../p2"
//	"./data"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"golang.org/x/crypto/sha3"
//	"io/ioutil"
//	"math"
//	"math/rand"
//	"net/http"
//	"net/url"
//	"strconv"
//	"strings"
//	"time"
//)

var MyPort int32
var MyID int32 = 0

var BLOCKCHAIN_SERVER = "http://localhost:6686"

//var REGISTER_SERVER = BLOCKCHAIN_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = BLOCKCHAIN_SERVER + "/upload"
var SELF_ADDR = "http://localhost:"

var SBC data.SyncBlockChain
var BlockchainPeers data.PeerList

var nextUserId = 0
var ifStarted bool

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
	//OutID = MyPort
	//SELF_ADDR = fmt.Sprintf("%v%v", SELF_ADDR, MyPort)
	//fmt.Printf("INIT: %v, %v\n", SELF_ADDR, TA_SERVER)
	SBC = data.NewBlockChain() // Init synch blockchain here
	mpt := p1.MerklePatriciaTrie{}
	mpt.NewTree()
	block := SBC.GenBlock(mpt, "")
	SBC.Insert(block)
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {
	// After register, send heartbeat every 5-10 seconds
	if ifStarted {
		return
	}
	ifStarted = true
	SELF_ADDR = fmt.Sprintf("%v%v", SELF_ADDR, MyPort)
	Register()
	StartHeartBeat()
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

	if SELF_ADDR != BLOCKCHAIN_SERVER {
		fmt.Printf("can download\n")
		Download()
	}
}

// Download blockchain from TA server
func Download() {
	fmt.Printf("downloading...\n")
	//resp, err := http.Get(BC_DOWNLOAD_SERVER) // GET to server
	newHeartBeat := data.NewHeartBeatData(false, BlockchainPeers.GetSelfId(), "", "", SELF_ADDR)
	jsonHBBytes, err := json.Marshal(newHeartBeat)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	form := url.Values{}
	form.Set("heartBeat", string(jsonHBBytes))
	resp, err := http.PostForm(BC_DOWNLOAD_SERVER, form) // POST to server
	//req, err := http.NewRequest("POST", BC_DOWNLOAD_SERVER, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body := resp.Body
	if err != nil {
		fmt.Printf("decode resp err: %v\n", err)
		return
	}

	respData, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Printf("decode resp err: %v\n", err)
		return
	}
	respString := string(respData)
	fmt.Printf("DOWNLOAD: decode resp: %v\n", respString)
	SBC.UpdateEntireBlockChain(respString)
}

func GiveClientId(w http.ResponseWriter, r *http.Request) {
	clientId := nextUserId
	nextUserId++
	fmt.Fprint(w, clientId)
}

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: get POST body of item json, seller id, (if actual app would have database w/ user ids, etc...)
	//		 as well as the demand json (desired currency)
	//fmt.Printf("Trade request ID: %v", id)
}

func ViewRequest(w http.ResponseWriter, r *http.Request) {
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

}

//TODO: In order to check the blockchain for trades, must access the blockchain somehow. In order to do this
//		must check the chain from the peerlist. Maybe this server is a listener on the blockchain network BUT
//		does NOT mine.

// NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string)
func StartHeartBeat() {
	randomRange := rand.Intn(10-5) + 5
	ticker := time.NewTicker(time.Duration(randomRange) * time.Second)
	go func() {
		for t := range ticker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			fmt.Println("Sending heartbeat...")

			pmJson, err := BlockchainPeers.PeerMapToJson()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			preparedData := data.PrepareHeartBeatData(&SBC, BlockchainPeers.GetSelfId(), pmJson, SELF_ADDR)
			preparedJsonBytes, err := json.Marshal(preparedData)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			preparedJson := string(preparedJsonBytes)
			BlockchainPeers.Rebalance() // Call rebalance to check if peerlist needs to be rebalanced before sending heartbeat
			for addr := range BlockchainPeers.PeerMap {
				if addr == SELF_ADDR {
					fmt.Println("\t\tFound own address.")
					continue
				} // Dont send to self
				postData := url.Values{}
				postData.Set("data", preparedJson)
				resp, err := http.PostForm(addr+"/heartbeat/receive", postData) // POST to server
				if err != nil {
					fmt.Printf("Heartbeat send rror: %v\n", err)
					return
				}
				resp.Body.Close()
			}

			randomRange = rand.Intn(10-5) + 5
			ticker = time.NewTicker(time.Duration(randomRange) * time.Second)
		}
	}()
}
