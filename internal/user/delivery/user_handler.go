package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/internal/user"
	"github.com/innovember/real-time-forum/pkg/response"
)

type UserHandler struct {
	userUcase user.UserUsecase
}

func NewUserHandler(userUcase user.UserUsecase) *UserHandler {
	return &UserHandler{
		userUcase: userUcase,
	}
}

func (uh *UserHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {
	mux.HandleFunc("/api/v1/user/signup", mm.CORSConfig(uh.HandlerRegisterUser))
}

func (uh *UserHandler) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var (
			input models.InputUserSignUp
			err   error
		)
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
			return
		}
		user := models.User{
			Nickname:  input.Nickname,
			Email:     input.Email,
			Password:  input.Password,
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Age:       input.Age,
			Gender:    input.Gender,
		}
		err = uh.userUcase.Create(&user)
		if err != nil {
			switch err {
			case consts.ErrEmailNotValid, consts.ErrNicknameTooShort, consts.ErrNicknameTooLong,
				consts.ErrNicknameAlreadyExist, consts.ErrEmailAlreadyExist,
				consts.ErrHashPassword:
				response.JSON(w, false, http.StatusBadRequest, err.Error(), nil)
				return
			default:
				response.JSON(w, false, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}
		response.JSON(w, true, http.StatusCreated, "new user has been created", nil)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyPOST.Error(), nil)
		return
	}
}
