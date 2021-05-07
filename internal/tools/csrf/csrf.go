package csrf

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
)

func NewCSRFToken(session *models.Session) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(session.Token))

	if err != nil {
		return "", err
	}

	hashedSession := hex.EncodeToString(hasher.Sum(nil))

	token := fmt.Sprintf("%s:%d", hashedSession, session.ExpiresAt.Unix())
	return token, nil
}

func ValidateCSRFToken(session *models.Session, token string) error {
	tokenData := strings.Split(token, ":")

	if len(tokenData) != 2 {
		return consts.ErrCSRF
	}
	hasher := sha256.New()
	_, err := hasher.Write([]byte(session.Token))
	if err != nil {
		return err
	}
	hashedSession := hex.EncodeToString(hasher.Sum(nil))
	if hashedSession != tokenData[0] {
		return consts.ErrCSRF
	}
	expiresAt, err := strconv.Atoi(tokenData[1])
	if err != nil {
		return consts.ErrCSRF
	}
	if int64(expiresAt) < time.Now().Unix() {
		return consts.ErrCSRF
	}
	return nil
}
