package usecases

import (
	"log"
	"time"

	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/session"
)

type SessionUsecase struct {
	sessionRepo session.SessionRepository
}

func NewSessionUsecase(sessionRepo session.SessionRepository) session.SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepo,
	}
}

func (sUc *SessionUsecase) CreateSession(session *models.Session) error {
	err := sUc.sessionRepo.Insert(session)
	if err != nil {
		return err
	}
	return nil
}

func (sUc *SessionUsecase) DeleteSession(token string) error {
	err := sUc.sessionRepo.Delete(token)
	if err != nil {
		return err
	}
	return nil
}

func (sUc *SessionUsecase) GetByToken(token string) (*models.Session, error) {
	session, err := sUc.sessionRepo.SelectByToken(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (sUc *SessionUsecase) DeleteExpiredSessions() {
	for {
		err := sUc.sessionRepo.DeleteTokens()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (sUc *SessionUsecase) CreateOnlineUser(userID int64) error {
	err := sUc.sessionRepo.InsertOnlineUser(userID)
	if err != nil {
		return err
	}
	return nil
}

func (sUc *SessionUsecase) DeleteOnlineUser(userID int64) error {
	err := sUc.sessionRepo.DeleteOnlineUser(userID)
	if err != nil {
		return err
	}
	return nil
}

func (sUc *SessionUsecase) DeleteOnlineUsers() {
	for {
		err := sUc.sessionRepo.DeleteOnlineUsers()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}
