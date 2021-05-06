package usecases

import (
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/user"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (uu *UserUsecase) Create(user *models.User) error {
	err := helpers.Validate(user)
	if err != nil {
		return err
	}
	name, err := uu.userRepo.SelectByNickname(user.Nickname)
	if err != nil {
		return err
	}
	if name != nil {
		return consts.ErrNicknameAlreadyExist
	}
	email, err := uu.userRepo.SelectByEmail(user.Email)
	if err != nil {
		return err
	}
	if email != nil {
		return consts.ErrEmailAlreadyExist
	}
	hashedPassword, err := helpers.Hash(user.Password)
	if err != nil {
		return consts.ErrHashPassword
	}
	user.Password = hashedPassword
	if err = uu.userRepo.Insert(user); err != nil {
		return err
	}
	return nil
}
