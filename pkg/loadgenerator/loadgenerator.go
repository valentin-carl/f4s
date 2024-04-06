package loadgenerator

import (
	"context"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

const (
	a = 97
	z = 122
)

type Message []byte

type LGContextKey int

const (
	Key0 LGContextKey = iota
)

func GenerateLoad(ctx context.Context, t <-chan time.Time, dst *string) {

	values := ctx.Value(Key0).(struct {
		id   int
		seed int64
		size int
	})

	var (
		id   = values.id
		seed = values.seed
		size = values.size
	)

	done := ctx.Done()

	log.Printf("load generator created with id %d", id)

	c, response, err := websocket.DefaultDialer.Dial(*dst, nil)
	if err != nil {
		log.Fatal(color.RedString("lg %d unable to open WebSocket connection: %s", id, err.Error()))
	}
	log.Printf("lg %d got reponse %s", id, response.Status)

	s := rand.NewSource(seed)
	r := rand.New(s)

	for {
		select {
		case <-t:
			{
				msg := GenerateMessage(r, size)
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Print(color.YellowString("lg %d unable to send message: %s", id, err.Error()))
					//return
				}
			}
		case <-done:
			{
				log.Print(color.YellowString("lg %d context cancelled", id))
				goto TheEnd
			}
		}
	}

TheEnd:
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Print(color.RedString("lg %d error while trying to send close message: %s", id, err.Error()))
	}
	log.Print(color.GreenString("lg %d done", id))
}

func GenerateMessage(r *rand.Rand, size int) Message {
	msg := make(Message, size)
	for i := 0; i < size; i++ {
		msg[i] = byte(r.Intn(z+1-a) + a)
	}
	return msg
}
