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
	if err != consts.ErrNoData && err != nil {
		return err
	}
	if name != nil {
		return consts.ErrNicknameAlreadyExist
	}
	email, err := uu.userRepo.SelectByEmail(user.Email)
	if err != consts.ErrNoData && err != nil {
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

func (uu *UserUsecase) GetByNickname(nickname string) (*models.User, error) {
	user, err := uu.userRepo.SelectByNickname(nickname)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUsecase) GetByEmail(email string) (*models.User, error) {
	user, err := uu.userRepo.SelectByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUsecase) GetByID(userID int64) (*models.User, error) {
	user, err := uu.userRepo.SelectByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUsecase) UpdateActivity(userID int64) error {
	err := uu.userRepo.UpdateActivity(userID)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) CheckPassword(input *models.InputUserSignIn) error {
	user, err := uu.GetByEmailOrNickname(input.Nickname)
	if err != nil {
		return err
	}
	if err = helpers.VerifyPassword(user.Password, input.Password); err != nil {
		return consts.ErrIncorrectPassword
	}
	return nil
}

func (uu *UserUsecase) GetByEmailOrNickname(login string) (*models.User, error) {
	name, err := uu.userRepo.SelectByNickname(login)
	if err != nil && err != consts.ErrNoData {
		return nil, err
	}
	if name != nil {
		return name, nil
	}
	email, err := uu.userRepo.SelectByEmail(login)
	if err != nil && err != consts.ErrNoData {
		return nil, err
	}
	if email == nil {
		return nil, consts.ErrUserNotExist
	}
	return email, nil
}

func (uu *UserUsecase) GetAllUsers() ([]models.User, error) {
	users, err := uu.userRepo.SelectAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}
