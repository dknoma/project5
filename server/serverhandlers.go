package server

import (
	"fmt"
	"net/http"
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

var nextUserId = 0

func GiveClientId(w http.ResponseWriter, r *http.Request) {
	clientId := nextUserId
	nextUserId++
	fmt.Fprint(w, clientId)
}

func TryToBuyItem(w http.ResponseWriter, r *http.Request) {

}
