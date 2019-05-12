package main

import (
	"fmt"
	"github.com/dknoma/project5/server/gamedata"
)

var MyPort = int32(6666)

func main() {
	//appServerRouter := server.NewAppServerRouter()
	//port := strconv.Itoa(int(MyPort))
	//fmt.Printf("port: %v\n", port)
	//log.Fatal(http.ListenAndServe(":"+port, appServerRouter))

	//testJson := "{\"tradeFulfillment\": {\"id\": 1,\"sellerId\": 6690,\"buyerId\": 6684,\"item\": {\"name\": \"sword\",\"id\": 2,\"owner\": 1,\"description\": \"This is a sword I got from a slime.\",\"stats\" : {\"level\": 1,\"atk\": 5,\"def\": 5}},\"sellerYield\": {\"currency\": 1000},\"buyerYield\": {\"currency\": -1000},\"minerYield\": {\"currency\": 10}}}"
	testJson := "{\"id\": 1,\"sellerId\": 6690,\"buyerId\": 6684,\"item\": {\"name\": \"sword\",\"id\": 2,\"owner\": 1,\"description\": \"This is a sword I got from a slime.\",\"stats\" : {\"level\": 1,\"atk\": 5,\"def\": 5}},\"sellerYield\": 1000,\"buyerYield\": -1000,\"minerYield\": 10}"
	f, err := gamedata.DecodeFulfillmentFromJSON(testJson)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("f: %v\n", f)

	testsJson := "[{\"id\": 1,\"sellerId\": 6690,\"buyerId\": 6684,\"item\": {\"name\": \"sword\",\"id\": 2,\"owner\": 1,\"description\": \"This is a sword I got from a slime.\",\"stats\" : {\"level\": 1,\"atk\": 5,\"def\": 5}},\"sellerYield\": 1000,\"buyerYield\": -1000,\"minerYield\": 10},{\"id\": 2,\"sellerId\": 6690,\"buyerId\": 6684,\"item\": {\"name\": \"sword\",\"id\": 2,\"owner\": 1,\"description\": \"This is a sword I got from a slime.\",\"stats\" : {\"level\": 1,\"atk\": 5,\"def\": 5}},\"sellerYield\": 1000,\"buyerYield\": -1000,\"minerYield\": 10}]"

	//fuls := gamedata.TradeFulfillments{}
	fuls := gamedata.TradeFulfillments{}
	works := fuls.RemoveFulfillments(testsJson)
	fmt.Printf("ladhgkadg: %v\n", works)

	//if len(os.Args) > 1 {
	//port, err := strconv.Atoi(os.Args[1])
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//p3.MyPort = MyPort
	//}
}
