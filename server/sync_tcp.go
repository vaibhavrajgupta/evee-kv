package server

import (
	"fmt"
	"io"
	"key-value_store/config"
	"key-value_store/core"
	"log"
	"net"
	"strconv"
	"strings"
)

func readCommand(c net.Conn) (*core.RedisCmd, error){
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil{
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])

	if err != nil{
		return nil, err
	}

	return &core.RedisCmd{
		Cmd: strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func respondError(err error, c net.Conn){
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c net.Conn){
	err := core.EvalAndRespond(cmd, c)
	if err != nil{
		respondError(err, c)
	}
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
			respond(cmd, c)
		}

	}
}