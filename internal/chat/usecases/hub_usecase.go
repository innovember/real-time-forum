package usecases

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	True  = true
	False = false
)

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

func (hu *HubUsecase) NewClient(userID int64, hub *models.Hub, conn *websocket.Conn, send chan *models.WsEvent) *models.Client {
	return &models.Client{
		UserID: userID,
		Hub:    hub,
		Conn:   conn,
		Send:   send,
		Mu:     sync.Mutex{},
	}
}

func (hu *HubUsecase) ServeWS(w http.ResponseWriter, r *http.Request, hub *models.Hub, roomID, userID int64) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := hu.NewClient(userID, hub, wsConn, make(chan *models.WsEvent))
	client.Hub.Register <- client

	go hu.WritePump(client, roomID)
	go hu.ReadPump(client, roomID)
}

func (hu *HubUsecase) writeJSON(c *models.Client, data interface{}) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return c.Conn.WriteJSON(data)
}

func (hu *HubUsecase) WritePump(c *models.Client, roomID int64) {
	go func() {
		ticker := time.NewTicker(pingPeriod)

		defer func() {
			ticker.Stop()
			c.Conn.Close()
		}()
		for {
			select {
			case event, ok := <-c.Send:
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					hu.writeJSON(c, &models.WsEvent{
						EventType:   models.WsEventType.WsClosed,
						RecipientID: c.UserID,
						RoomID:      roomID,
						State:       false,
					})
					return
				}
				hu.writeJSON(c, event)
			case <-ticker.C:
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := hu.writeJSON(c, &models.WsEvent{
					EventType:   models.WsEventType.PingMessage,
					RecipientID: c.UserID,
					RoomID:      roomID,
					State:       false,
				}); err != nil {
					return
				}

			}
		}
	}()
}

func (hu *HubUsecase) ReadPump(c *models.Client, roomID int64) {
	go func() {
		defer func() {
			c.Hub.Unregister <- c
			hu.writeJSON(c, &models.WsEvent{
				EventType:   models.WsEventType.WsClosed,
				RecipientID: c.UserID,
				RoomID:      roomID,
				State:       false,
			})
			c.Conn.Close()
		}()

		c.Conn.SetReadLimit(maxMessageSize)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		c.Conn.SetPongHandler(func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		for {
			event, err := GetEvent(c)
			if err != nil {
				hu.writeJSON(c, &models.WsEvent{
					EventType:   models.WsEventType.WsError,
					Content:     err.Error(),
					RecipientID: c.UserID,
					RoomID:      roomID,
					State:       false,
				})
				return
			}
			switch event.EventType {
			case models.WsEventType.Message:
				err = hu.CreateMessage(c, roomID, &event)

			case models.WsEventType.PongMessage:
				c.Conn.SetReadDeadline(time.Now().Add(pongWait))

			case models.WsEventType.TypingInReceiver:
				err = hu.SendTypingInResponse(c, roomID, &event)

			default:
				err = consts.ErrEventType
			}

			if err != nil {
				hu.writeJSON(c, &models.WsEvent{
					EventType:   models.WsEventType.WsError,
					Content:     err.Error(),
					RecipientID: c.UserID,
					RoomID:      roomID,
					State:       false,
				})
				return
			}
		}
	}()
}

func GetEvent(c *models.Client) (models.WsEvent, error) {
	var event models.WsEvent
	_, messageBytes, err := c.Conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("ws error: %v", err)
		}
		return event, err
	}
	err = json.Unmarshal(messageBytes, &event)
	if err != nil {
		return event, err
	}
	return event, nil
}

func (hu *HubUsecase) CreateMessage(c *models.Client, roomID int64, event *models.WsEvent) error {
	inputMsg := models.Message{Content: event.Content}
	eventMsg := models.WsEvent{EventType: models.WsEventType.Message}
	if inputMsg.Content != "" && strings.TrimSpace(inputMsg.Content) != "" {
		inputMsg.RoomID = roomID
		user := &models.User{ID: c.UserID}
		inputMsg.User = user
		inputMsg.MessageDate = helpers.GetCurrentUnixTime()
		outputMessage, err := hu.roomRepo.InsertMessage(&inputMsg)
		if err != nil {
			return err
		}
		eventMsg.Message = outputMessage
		c.Hub.Broadcast <- &eventMsg
	}
	return nil
}

func (hu *HubUsecase) SendTypingInResponse(c *models.Client, roomID int64, event *models.WsEvent) error {
	users, err := hu.roomRepo.SelectUsersByRoom(roomID)
	if err != nil {
		return err
	}
	receiverID := helpers.SelectSecondUser(users, c.UserID)
	if event.TypingInReceiverID == c.UserID {
		return consts.ErrTypingInSameUser
	}
	if event.TypingInReceiverID != receiverID {
		return consts.ErrTypingIn
	}
	eventMsg := models.WsEvent{
		EventType:      models.WsEventType.TypingInSender,
		RecipientID:    event.TypingInReceiverID,
		TypingInUserID: c.UserID,
	}
	c.Hub.Broadcast <- &eventMsg
	return nil
}
