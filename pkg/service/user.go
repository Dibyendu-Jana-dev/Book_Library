package service

import (
	"context"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dibyendu/Authentication-Authorization/lib/constants"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
	"github.com/dibyendu/Authentication-Authorization/pkg/domain"
	"github.com/dibyendu/Authentication-Authorization/pkg/dto"
	"github.com/dibyendu/Authentication-Authorization/pkg/httpclient/fetcherService"
	"github.com/dibyendu/Authentication-Authorization/pkg/middleware"
	"strings"
)

type DefaultUserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) DefaultUserService {
	return DefaultUserService{
		repo: repo,
	}
}

type UserService interface {
	SignIn(ctx context.Context, request dto.CreateUserRequest) (*dto.UserLoginResponse, *errs.AppError)
	GetBookList(ctx context.Context, userType string) ([]*dto.GetUserBookListResponse, *errs.AppError)
	AddBook(ctx context.Context, request dto.AddBookRequest, token string) ([]*dto.GetUserBookListResponse, *errs.AppError)
	DeleteBook(ctx context.Context, bookName, token string) ([]*dto.GetUserBookListResponse, *errs.AppError)
}

func (s DefaultUserService) SignIn(ctx context.Context, request dto.CreateUserRequest) (*dto.UserLoginResponse, *errs.AppError) {
	var (
		validateErr *errs.AppError
	)
	validateErr = request.Validate()
	if validateErr != nil {
		return nil, validateErr
	}
	req := domain.CreateUserRequest(request)
	data, err1 := s.repo.IsEmailExists(ctx, req.Email)
	if err1 != nil {
		return nil, err1
	}
	//matching the password
	err := utility.VerifyPassword(data.Password, request.Password)
	if err != nil {
		return nil, errs.NewValidationError(constants.PASSWORD_IS_NOT_VALID)
	}
	token, err := middleware.GenerateJWT(data.Email, data.Role)
	if err != nil {
		return nil, errs.NewValidationError(constants.TOKEN_IS_NOT_VALID)
	}
	return &dto.UserLoginResponse{
		Name:  data.Name,
		Role:  data.Role,
		Email: data.Email,
		Token: token,
	}, nil
}

func(s DefaultUserService) GetBookList(ctx context.Context, userType string) ([]*dto.GetUserBookListResponse, *errs.AppError){
	data, err := s.repo.GetBookList(ctx, userType)
	if err != nil {
		return nil, err
	}
	response := domain.ToDtoSlice(data)
	return response, nil
}

func (s DefaultUserService) AddBook(ctx context.Context, request dto.AddBookRequest, token string) ([]*dto.GetUserBookListResponse, *errs.AppError) {
	var (
		response    []*dto.GetUserBookListResponse
		validateErr *errs.AppError
	)
	//useAccess := middleware.GetUserInfo(ctx)
	//if ok := strings.EqualFold(strings.ToLower(useAccess.Role),"admin"); !ok {
	//	return nil, errs.NewValidationError("userr cannot be empty")
	//}
	validateErr = request.Validate()
	if validateErr != nil {
		return nil, validateErr
	}

	if request.PublicationYear < 0 || request.PublicationYear > 9999 {
		return nil, errs.NewValidationError(constants.INVALID_PUBLICATION_YEAR)
	}
	req := domain.AddBookRequest(request)
	_, err := s.repo.AddBook(ctx, req)

	if err != nil {
		return nil, err
	}
	data, err:= fetcherService.GetDetailFromHome(token)
	if err != nil {
		return nil, err
	}

	response = fetcherService.ToDtoSlice(data)
	return response, nil
}

func (s DefaultUserService) DeleteBook(ctx context.Context, bookName, token string) ([]*dto.GetUserBookListResponse, *errs.AppError) {
	var (
		response    []*dto.GetUserBookListResponse
	)
	if stringUtils.IsBlank(bookName) {
		return nil, errs.NewValidationError(constants.BOOK_NAME_CAN_NOT_BE_EMPTY)
	}

	_, err := s.repo.DeleteBook(ctx, strings.TrimSpace(bookName))
	if err != nil {
		return nil, err
	}
	data, err:= fetcherService.GetDetailFromHome(token)
	if err != nil {
		return nil, err
	}

	response = fetcherService.ToDtoSlice(data)
	return response, nil
}