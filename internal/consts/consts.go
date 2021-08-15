package consts

import (
	"database/sql"
	"errors"
	"time"
)

type ctxParam string

var (
	ErrInternal             = errors.New("internal server error") //Internal server error
	ErrBadRequest           = errors.New("bad request")           //Bad request received
	ErrEmailAlreadyExist    = errors.New("еmail already exist")   //Email already exists
	ErrEmailNotValid        = errors.New("invalid email")
	ErrNicknameAlreadyExist = errors.New("nickname already exist") //Nickname already exists
	ErrNicknameTooShort     = errors.New("nickname too short, at least 3 char required")
	ErrNicknameTooLong      = errors.New("nickname too long, at most 15 char required")
	ErrIncorrectNickname    = errors.New("invalid login")    //Incorrect nickname
	ErrIncorrectPassword    = errors.New("invalid password") //Incorrect password
	ErrNotAuthorized        = errors.New("unauthorized")     //Not authorized
	ErrUserNotExist         = errors.New("user not found")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrOnlyPOST             = errors.New("only POST method allowed")
	ErrOnlyGet              = errors.New("only GET method allowed")
	ErrOnlyDelete           = errors.New("only DELETE requests allowed")
	ErrOnlyPUT              = errors.New("only PUT requests allowed")
	ErrHashPassword         = errors.New("hash password error")
	ErrCSRF                 = errors.New("invalid csrf token received")
	ErrInvalidSessionToken  = errors.New("invalid session token received")
	ErrSessionTokenNotFound = errors.New("invalid session token not found")
	ErrNoData               = sql.ErrNoRows
	ErrRoomNotExist         = errors.New("there is no room with such id")
	ErrHubNotFound          = errors.New("hub not found")
	ErrSelfChat             = errors.New("change useriD, you cant use your ID")
	ErrEventType            = errors.New("invalid event type")
	ErrWsClosed             = errors.New("websocket conn closed")
	ErrTypingIn             = errors.New("given recipient not found in room")
	ErrTypingInSameUser     = errors.New("given recipient can not send typing request to himself")

	RegistrationSuccess   = "You have registered successfully"
	ProfileSuccess        = "User's profile fetched successfully"
	LogoutSuccess         = "You have logged out"
	AllUsers              = "list of all users"
	AllOnlineUsers        = "list of all online users"
	StatusOnline          = "online"
	StatusOffline         = "offline"
	PostCreated           = "Your post has been created"
	CommentCreated        = "Your comment has been created"
	AllCategories         = "List of all categories"
	UserData              = "User data fetched"
	LoginSuccess          = "You have logged in"
	PostByIDSuccess       = "Post data fetched"
	ErrPostNotExist       = errors.New("there is no post with such id")
	Posts                 = "Posts fetched"
	Comments              = "Comments fetched"
	AllChats              = "All chats fetched"
	RoomMessages          = "Messages fetched by roomID"
	Room                  = "Room ready for chat"
	MessageStatusUpdated  = "Message Status updated"
	MessagesStatusUpdated = "Status of messages for receiver updated"
)

const (
	SessionName                    = "real_time_forum_session_id"
	SessionExpireDuration          = 1 * time.Hour
	ConstAuthedUserParam  ctxParam = "authorized_user"
)
