package server

import (
	"fmt"
	"github.com/dknoma/cs686-blockchain-p3-dknoma/p3/data"
	"net/http"
	"strconv"
	"strings"
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

var BLOCKCHAIN_SERVER = "http://localhost:6686"

var SBC data.SyncBlockChain
var BlockchainPeers data.PeerList

var nextUserId = 0

func GiveClientId(w http.ResponseWriter, r *http.Request) {
	clientId := nextUserId
	nextUserId++
	fmt.Fprint(w, clientId)
}

func TryToBuyItem(w http.ResponseWriter, r *http.Request) {
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
