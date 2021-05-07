package consts

import (
	"errors"
)

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
	ErrHashPassword         = errors.New("hash password error")
)
