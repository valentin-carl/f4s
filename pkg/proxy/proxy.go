package proxy

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
)

const (
	bufferSize = 1024 // TODO large enough ???
)

func Proxy(port int) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(color.RedString("proxy error while creating tcp listener: %s", err.Error()))
	}

	defer func() {
		log.Print("closing TCP listener")
		_ = listener.Close()
	}()

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Print(color.RedString("error while accepting client connection"))
			break // TODO reconsider this break statement
		}
		go HandleConnection(clientConnection)
	}

	log.Print(color.YellowString("proxy done; stopped accepting connections"))
}

func HandleConnection(cc net.Conn) {

	// TODO
	//  quadruple-check that ONLY the correct connections are closed and other streams are not affected

	defer func() {
		log.Print(color.CyanString("proxy closing client connection %s <--> %s", cc.LocalAddr(), cc.RemoteAddr()))
		err := cc.Close()
		if err != nil {
			log.Print(color.YellowString("proxy could not close client connection %s <--> %s: %s", cc.LocalAddr(), cc.RemoteAddr(), err.Error()))
		}
	}()

	serverAddr, request, functionName := GetFunctionAddress(cc)
	sc, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal(color.RedString("proxy error while trying to connect to server %s: %s", serverAddr, err.Error()))
	}

	proxyID := ProxyID(cc, sc, functionName)

	defer func() {
		log.Print(color.CyanString("proxy closing server connection %s <--> %s", sc.LocalAddr(), sc.RemoteAddr()))
		err := sc.Close()
		if err != nil {
			log.Print(color.YellowString("proxy could not close server connection %s <--> %s: %s", sc.LocalAddr(), sc.RemoteAddr(), err.Error()))
		}
	}()

	_, err = sc.Write(request)
	if err != nil {
		log.Fatal(color.RedString("proxy %s error while trying to forward initial request from client to server: %s", proxyID, err.Error()))
	}

	go func() {
		// client -> server
		_, err := io.Copy(sc, cc)
		if err != nil {
			log.Fatal(color.RedString("proxy %s error while forwarding client->server: ", proxyID, err.Error()))
		}
		log.Print(color.YellowString("proxy %s io copy server->client done", proxyID))
	}()

	// server -> client
	_, err = io.Copy(cc, sc)
	if err != nil {
		log.Fatal(color.RedString("proxy %s error while forwarding server->client: ", proxyID, err.Error()))
	}
	log.Print(color.YellowString("proxy %s io copy client->server done", proxyID))

	log.Print(color.GreenString("proxy %s done", proxyID))

	// TODO
	//  tell the controller that this stream is done
	//  -> do it here instead of in the function instance
}

func GetFunctionAddress(cc net.Conn) (string, []byte, string) {
	// returns: function url, initial client request, function name

	// TODO
	//  the stuff the proxy reads here will also have to be sent to the server
	//  make absolutely sure that the buffer is large enough
	//  if we loose data here, the connection will likely have problems

	// get function name
	buffer := make([]byte, bufferSize)
	n, err := cc.Read(buffer)
	if err != nil {
		log.Fatal(color.RedString("proxy error while reading client request to extract function name: %s", err.Error()))
	}
	request := string(buffer[:n])

	log.Printf("initial request: %s", request)

	r := regexp.MustCompile(` /(.*) HTTP`)
	results := r.FindStringSubmatch(string(buffer))
	functionName := results[1]
	log.Printf("proxy read function name: %s", functionName)

	// request function address
	// TODO
	//  gRPC to controller
	addr := "ws://localhost:12345/"
	return addr, buffer[:n], functionName
}

func ProxyID(cc net.Conn, sc net.Conn, functionName string) string {

	s := strings.Split(cc.RemoteAddr().String(), ":")
	clientPortRemote := s[len(s)-1]
	s = strings.Split(cc.LocalAddr().String(), ":")
	clientPortLocal := s[len(s)-1]
	s = strings.Split(sc.RemoteAddr().String(), ":")
	serverPortRemote := s[len(s)-1]
	s = strings.Split(sc.LocalAddr().String(), ":")
	serverPortLocal := s[len(s)-1]

	proxyID := fmt.Sprintf("%s[%s<->%s|%s<->%s]", functionName, clientPortRemote, clientPortLocal, serverPortLocal, serverPortRemote)
	return proxyID
}
