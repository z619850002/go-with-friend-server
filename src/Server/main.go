package main

import (
	"log"
	"fmt"
	"net/http"
	"go-with-friend-server/src/Server/Utils"
						)

func main() {

	hub := Utils.NewHub()
	go hub.Run()
	fmt.Println("start to listrn port : 8005, waiting... ")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Println("a connection from ：" + r.RemoteAddr)
		Utils.ServeWs(hub, w, r)
	})


	if err := http.ListenAndServe(":8005" , nil); err != nil {
		log.Fatal("ListenAndServe:" , err)
	}



}


