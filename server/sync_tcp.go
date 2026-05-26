package server

import (
	"io"
	"key-value_store/config"
	"log"
	"net"
	"strconv"
)

func readCommand(c net.Conn) (string, error){
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil{
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error{
	if _, err := c.Write([]byte(cmd)); err != nil{
		return err
	}
	return nil
}

func RunSyncTCPServer(){
	log.Println("Starting a sync TCP Server on", config.Host, config.Port)
	
	var con_clients int = 0

  	// listening to the configured host:port
	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil{
		panic(err)
	}
	
	for{
		// blocking call : waiting for the new client to connect
		c, err := lsnr.Accept()
		if err != nil{
			panic(err)
		}

		// increment the number of concurrent clients
		con_clients += 1
		log.Println("Client connected with address:", c.RemoteAddr(), "Concurrent Clients", con_clients)

		for{
			// over the socket, continuously read the command and print it out
			cmd, err := readCommand(c)

			if err != nil{
				c.Close()
				con_clients -= 1
				log.Println("Client disconnected", c.RemoteAddr(), "Concurrent Clients", con_clients)

				if err == io.EOF{
					break
				}
				log.Println("err", err)
			}
			log.Println("Command", cmd)
			if err = respond(cmd, c); err != nil{
				log.Println("err writes:", err)
			}
		}

	}
}