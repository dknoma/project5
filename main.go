package main

import (
	"./server"
	"fmt"
	"log"
	"net/http"
	"os"
)

var MyPort = int32(8000)

func main() {
	appServerRouter := server.NewAppServerRouter()
	if len(os.Args) > 1 {
		//port, err := strconv.Atoi(os.Args[1])
		//if err != nil {
		//	fmt.Printf("Error: %v\n", err)
		//	return
		//}
		//p3.MyPort = MyPort
		fmt.Printf("port: %v\n", os.Args[1])
		log.Fatal(http.ListenAndServe(":"+os.Args[1], appServerRouter))
	}
}
