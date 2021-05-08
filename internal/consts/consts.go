package consts

import (
	"database/sql"
	"errors"
	"time"
)

type ctxParam string

var (
	ErrInternal             = errors.New("внутренняя ошибка сервера") //Internal server error
	ErrBadRequest           = errors.New("некорректный запрос")       //Bad request received
	ErrEmailAlreadyExist    = errors.New("еmail уже существует")      //Email already exists
	ErrEmailNotValid        = errors.New("еmail невалидный")
	ErrNicknameAlreadyExist = errors.New("имя пользователя уже существует") //Nickname already exists
	ErrNicknameTooShort     = errors.New("имя пользователя должно быть более 3 символов")
	ErrNicknameTooLong      = errors.New("имя пользователя должно быть менее 15 символов")
	ErrIncorrectNickname    = errors.New("неверный логин")   //Incorrect nickname
	ErrIncorrectPassword    = errors.New("неверный  пароль") //Incorrect password
	ErrNotAuthorized        = errors.New("не авторизован")   //Not authorized
	ErrUserNotExist         = errors.New("пользователь не найден")
	ErrPermissionDenied     = errors.New("в доступе отказано")
	ErrOnlyPOST             = errors.New("разрешены только POST запросы")
	ErrOnlyGet              = errors.New("разрешены только GET запросы")
	ErrHashPassword         = errors.New("ошибка в хеширование пароля")
	ErrCSRF                 = errors.New("invalid csrf token received")
	ErrInvalidSessionToken  = errors.New("invalid session token received")
	ErrSessionTokenNotFound = errors.New("invalid session token not found")
	ErrNoData               = sql.ErrNoRows
	RegistrationSuccess     = "You have registered successfully"
	ProfileSuccess          = "User's profile fetched successfully"
	ErrOnlyDelete           = errors.New("only delete requests allowed")
	LogoutSuccess           = "You have logged out"
)

const (
	SessionName                    = "real_time_forum_session_id"
	SessionExpireDuration          = 1 * time.Hour
	ConstAuthedUserParam  ctxParam = "authorized_user"
)
