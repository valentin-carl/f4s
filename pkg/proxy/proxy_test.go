package proxy

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {

	const proxyPort = 1234
	const nClients = 5

	t.Run("proxy test", func(t *testing.T) {

		log.Print("hi")

		const proxyUrl = "ws://localhost:1234/echo"
		const serverUrl = ":8080"

		client := func(id int, d time.Duration) {

			c, _, err := websocket.DefaultDialer.Dial(proxyUrl, nil)
			if err != nil {
				log.Printf("client %d couldn't connect to proxy: %s", id, err.Error())
				return
			}
			defer func() {
				log.Printf("client %d closing ws connection", id)
				err := c.Close()
				if err != nil {
					log.Fatalf("client %d error while closing ws connection: %s", id, err.Error())
				}
			}()

			t1 := time.NewTimer(d)
			t2 := time.NewTicker(time.Second)

			i := 0
			for {
				select {
				case <-t1.C:
					{
						log.Printf("client %d timer over", id)
						goto TheEnd
					}
				case <-t2.C:
					{
						log.Printf("client %d sending message", id)
						err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("message %d %d", id, i)))
						if err != nil {
							log.Printf("client %d error while sending message: %s", id, err.Error())
							goto TheEnd
						}
						i++
					}
				}
			}

		TheEnd:
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("client %d error while trying to send ws close message: %s", id, err.Error())
			}
			log.Printf("client %d done", id)
			time.Sleep(time.Second)
		}

		server := func() {
			upgrader := websocket.Upgrader{}
			http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
				c, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					log.Printf("server couldn't upgrade HTTP connection to WS: %s", err.Error())
				}
				defer func() {
					log.Printf("server closing connection to %s", c.RemoteAddr())
					err := c.Close()
					if err != nil {
						log.Fatalf("server error while closing WS connection: %s", err.Error())
					}
				}()
				for {
					_, m, err := c.ReadMessage()
					if err != nil {
						log.Printf("server error while reading message: %s", err.Error())
						break
					}
					log.Print("server got message: " + color.GreenString("%s", string(m)))
				}
				log.Print("server finished handling connection")
			})
			log.Fatal(http.ListenAndServe(serverUrl, nil))
		}

		go server()
		go Proxy(proxyPort)
		var wg sync.WaitGroup
		for i := 0; i < nClients; i++ {
			wg.Add(1)
			go func(x int) {
				log.Printf("starting client %d", x)
				client(x, time.Duration((x+1)*5)*time.Second)
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}
