package main

import (
	"./server"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var MyPort = int32(6666)

func main() {
	appServerRouter := server.NewAppServerRouter()
	port := strconv.Itoa(int(MyPort))
	fmt.Printf("port: %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, appServerRouter))
	//if len(os.Args) > 1 {
	//port, err := strconv.Atoi(os.Args[1])
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}
	//p3.MyPort = MyPort
	//}
}
