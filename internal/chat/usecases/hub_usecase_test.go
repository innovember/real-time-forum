package usecases_test

import (
	"fmt"
	"testing"
)

func TestNewHub(t *testing.T) {
	_, hubUsecase := setupChatUsecases()
	hub := hubUsecase.NewHub()
	hubUsecase.Register(1, hub)
	fmt.Println(hub)
}

func TestGetHub(t *testing.T) {
	_, hubUsecase := setupChatUsecases()
	hub := hubUsecase.NewHub()
	hubUsecase.Register(1, hub)
	h, err := hubUsecase.GetHub(1)
	if err != nil {
		t.Error("get hub err ", err)
	}
	fmt.Println(h)
}
