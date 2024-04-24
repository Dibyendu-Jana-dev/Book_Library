package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/pkg/dto"
	"github.com/dibyendu/Authentication-Authorization/pkg/middleware"
	"github.com/dibyendu/Authentication-Authorization/pkg/service"
)

type UserHandler struct {
	Service service.UserService
}

func (h UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		request dto.CreateUserRequest
		ctx     = r.Context()
	)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		response, appError := h.Service.SignIn(ctx, request)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, response)
		}
	}
}

func (h UserHandler) GetBookList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	userInfo := middleware.GetUserInfo(ctx)
	if strings.EqualFold(strings.ToLower(userInfo.AuthToken), "invalid") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("unauthorized Access for get user").AsMessage())
		return
	}
	//if !strings.EqualFold(strings.ToLower(userInfo.Role), "admin") {
	//	writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("wrong user it's only for admin").AsMessage())
	//	return
	//}
	response, err := h.Service.GetBookList(ctx, userInfo.Role)
	if err != nil {
		if err.Code == http.StatusNoContent {
			writeResponseNoContent(w, err.Code)
		} else {
			writeResponse(w, err.Code, err.AsMessage())
		}
	} else {
		writeResponse(w, http.StatusOK, response)
	}
}

func (h UserHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	var (
		request dto.AddBookRequest
		ctx     = r.Context()
	)
	userInfo := middleware.GetUserInfo(ctx)
	if strings.EqualFold(strings.ToLower(userInfo.AuthToken), "invalid") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("unauthorized Access for get user").AsMessage())
		return
	}
	if !strings.EqualFold(strings.ToLower(userInfo.Role), "admin") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("only admin can access it").AsMessage())
		return
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		response, appError := h.Service.AddBook(ctx, request, userInfo.AuthToken)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, response)
		}
	}
}

func (h UserHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	var bookName string
	if !stringUtils.IsBlank(r.URL.Query().Get("book_name")) {
		bookName = r.URL.Query().Get("book_name")
	}
	ctx := r.Context()
	userInfo := middleware.GetUserInfo(ctx)
	if strings.EqualFold(strings.ToLower(userInfo.AuthToken), "invalid") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("unauthorized Access for get user").AsMessage())
		return
	}
	if !strings.EqualFold(strings.ToLower(userInfo.Role), "admin") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("wrong user it's only for admin access").AsMessage())
		return
	}
	user, appError := h.Service.DeleteBook(ctx, bookName, userInfo.AuthToken)
	if appError != nil {
		if appError.Code == http.StatusNoContent {
			writeResponseNoContent(w, appError.Code)
		} else {
			writeResponse(w, appError.Code, appError.AsMessage())
		}
	} else {
		writeResponse(w, http.StatusOK, user)
	}
}
