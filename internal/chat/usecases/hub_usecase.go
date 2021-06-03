package usecases

import (
	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
)

type HubUsecase struct {
	hubRepo chat.HubRepository
}

func NewHubUsecase(hubRepo chat.HubRepository) *HubUsecase {
	return &HubUsecase{
		hubRepo: hubRepo,
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

func (hu *HubUsecase) ServeWS() {

}
