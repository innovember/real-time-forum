package delivery

import (
	"net/http"

	"github.com/innovember/real-time-forum/internal/category"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/mwares"
	"github.com/innovember/real-time-forum/pkg/response"
)

type CategoryHandler struct {
	categoryUcase category.CategoryUsecase
}

func NewCategoryHandler(categoryUcase category.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUcase: categoryUcase,
	}
}

func (ch *CategoryHandler) Configure(mux *http.ServeMux, mm *mwares.MiddlewareManager) {
	mux.HandleFunc("/api/v1/categories", mm.CORSConfig(ch.HandlerGetAllCategories))
}

func (ch *CategoryHandler) HandlerGetAllCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories, err := ch.categoryUcase.GetAllCategories()
		if err != nil {
			response.JSON(w, true, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		response.JSON(w, true, http.StatusOK, consts.AllCategories, categories)
		return
	default:
		response.JSON(w, false, http.StatusMethodNotAllowed, consts.ErrOnlyGet.Error(), nil)
		return
	}
}
