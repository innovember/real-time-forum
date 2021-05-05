package helpers

import (
	"regexp"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
)

func Validate(user *models.User) error {
	if !CheckValidEmail(user.Email) {
		return consts.ErrEmailNotValid
	}
	err := CheckNickname(user.Nickname)
	if err != nil {
		return err
	}
	return nil
}

func CheckValidEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	emailLength := len(email)

	if emailLength < 3 || emailLength > 254 {
		return false
	}

	emailMatchResult := emailRegex.MatchString(email)

	return emailMatchResult
}

func CheckNickname(nickname string) error {

	if len(nickname) == 0 || len(nickname) < 3 {
		return consts.ErrNicknameTooShort
	}

	if len(nickname) > 15 {
		return consts.ErrNicknameTooLong
	}
	return nil
}
