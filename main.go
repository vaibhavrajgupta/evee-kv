package main

import (
	"flag"
	"key-value_store/config"
	"key-value_store/server"
	"log"
)


func setupFlags(){
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the key-value_store server")
	flag.IntVar(&config.Port, "port", 7379, "port for the key-value_store store")
	flag.Parse()
}

func main(){
	setupFlags()
	log.Println("Store is starting.....")
	server.RunSyncTCPServer()
}