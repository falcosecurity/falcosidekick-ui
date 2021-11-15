package broadcast

import (
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
	echo "github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Broadcast struct {
	Channel chan string
}

var broadcast *Broadcast

func CreateBroadcast() *Broadcast {
	broadcast = new(Broadcast)
	return broadcast
}

func GetBroadcast() *Broadcast {
	return broadcast
}

func (b *Broadcast) BroadcastMessage() {
	b.Channel <- "New Event"
}

func WebSocketBroadcast(c echo.Context) error {
	utils.WriteLog("info", "New WebSocket Connection", false)
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		b := GetBroadcast()
		for {
			<-b.Channel
			err := websocket.Message.Send(ws, b)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
