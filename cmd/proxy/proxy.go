package main

import (
	"gitcom.com/valentin-carl/f4s/pkg/proxy"
	"github.com/fatih/color"
	"log"
	"os"
	"strconv"
)

const (
	PROXY_PORT_DEFAULT = 1234
)

func main() {
	log.Print(color.GreenString("starting proxy"))

	var proxyPort int
	if pp := os.Getenv("PROXY_PORT"); pp == "" {
		proxyPort = PROXY_PORT_DEFAULT
	} else {
		ppi, err := strconv.Atoi(pp)
		if err != nil {
			log.Print(color.YellowString("env var PROXY_PORT is not a number (got '%s'), using default instead: %s", pp, err.Error()))
			proxyPort = PROXY_PORT_DEFAULT
		} else {
			proxyPort = ppi
		}
	}

	proxy.Proxy(proxyPort)

	log.Print(color.GreenString("proxy done"))

	// TODO
	//  check this out: https://github.com/fatih/pool & https://github.com/silenceper/pool
	//  maybe create a connection pool for each function instead of handling each connection individually
}
