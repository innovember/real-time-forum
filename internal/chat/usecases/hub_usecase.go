package usecases

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
)

type HubUsecase struct {
	hubRepo  chat.HubRepository
	roomRepo chat.RoomRepository
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHubUsecase(hubRepo chat.HubRepository,
	roomRepo chat.RoomRepository) *HubUsecase {
	return &HubUsecase{
		hubRepo:  hubRepo,
		roomRepo: roomRepo,
	}
}

func (hu *HubUsecase) NewHub() *models.Hub {
	hub := hu.hubRepo.NewHub()
	return hub
}

func (hu *HubUsecase) GetHub(roomID int64) (*models.Hub, error) {
	hub, ok := hu.hubRepo.GetHub(roomID)
	if !ok {
		return nil, consts.ErrHubNotFound
	}
	return hub, nil
}

func (hu *HubUsecase) DeleteHub(roomID int64) {
	hu.hubRepo.DeleteHub(roomID)
}

func (hu *HubUsecase) Register(roomID int64, hub *models.Hub) {
	hu.hubRepo.Register(roomID, hub)
}

func (hu *HubUsecase) NewClient(userID int64, hub *models.Hub, conn *websocket.Conn, send chan *models.Message) *models.Client {
	return &models.Client{
		UserID: userID,
		Hub:    hub,
		Conn:   conn,
		Send:   send,
		Mu:     sync.Mutex{},
	}
}

func (hu *HubUsecase) ServeWS(w http.ResponseWriter, r *http.Request, hub *models.Hub, userID int64) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := hu.NewClient(userID, hub, wsConn, make(chan *models.Message))
	client.Hub.Register <- client

	go hu.WritePump(client)
	go hu.ReadPump(client)
}

func (hu *HubUsecase) writeJSON(c *models.Client, data interface{}) error {
	c.Mu.Lock()
	err := c.Conn.WriteJSON(data)
	c.Mu.Unlock()
	return err
}

func (hu *HubUsecase) WritePump(c *models.Client) {
	go func() {
		ticker := time.NewTicker(pingPeriod)

		defer func() {
			ticker.Stop()
			c.Conn.Close()
		}()
		for {
			select {
			case message, ok := <-c.Send:
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					hu.writeJSON(c, &models.Message{
						HTTPCode: websocket.CloseMessage,
						State:    false,
					})
					return
				}
				hu.writeJSON(c, message)
			case <-ticker.C:
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := hu.writeJSON(c, &models.Message{
					HTTPCode: websocket.PingMessage,
					State:    false,
				}); err != nil {
					return
				}

			}
		}
	}()
}

func (hu *HubUsecase) ReadPump(c *models.Client) {
	go func() {
		defer func() {
			c.Hub.Unregister <- c
			c.Conn.Close()
		}()

		c.Conn.SetReadLimit(maxMessageSize)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.Conn.SetPongHandler(func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		for {
			_, messageBytes, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}
			inputMsg := models.Message{}
			json.Unmarshal(messageBytes, &inputMsg)
			user := &models.User{ID: c.UserID}
			inputMsg.User = user
			inputMsg.MessageDate = helpers.GetCurrentUnixTime()
			outputMessage, err := hu.roomRepo.InsertMessage(&inputMsg)
			if err != nil {
				log.Println("insert message err ,error: ", err)
				continue
			}
			outputMessage.HTTPCode = 200
			outputMessage.State = true
			if outputMessage.User.ID == c.UserID {
				outputMessage.IsYourMessage = true
			}
			hu.writeJSON(c, outputMessage)
			c.Hub.Broadcast <- outputMessage
		}
	}()
}
