package domain

import (
	"context"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/pkg/dto"
)

type UserRepository interface {
	IsEmailExists(ctx context.Context, email string) (*CreateUserResponse, *errs.AppError)
	GetBookList(ctx context.Context, userType string) ([]*GetUserBookListResponse, *errs.AppError)
	AddBook(ctx context.Context, req AddBookRequest) ([]*GetUserBookListResponse, *errs.AppError)
	DeleteBook(ctx context.Context, bookName string) ([]*GetUserBookListResponse, *errs.AppError)
}

type CreateUserRequest struct {
	Name     string `bson:"name"`
	Role     string `bson:"role"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type CreateUserResponse struct {
	Name     string `bson:"name"`
	Role     string `bson:"role"`
	Email    string `bson:"email"`
	Password string `bson:"password,omitempty"`
}

type GetUserBookListResponse struct {
	BookName        string `json:"book_name"`
	Author          string `json:"author"`
	PublicationYear string `json:"publication_year"`
}

func (r GetUserBookListResponse) ToDto() *dto.GetUserBookListResponse {
	return &dto.GetUserBookListResponse{
		BookName:        r.BookName,
		Author:          r.Author,
		PublicationYear: r.PublicationYear,
	}
}

func ToDtoSlice(r []*GetUserBookListResponse) []*dto.GetUserBookListResponse {
	dtoSlice := make([]*dto.GetUserBookListResponse, len(r))
	for i, v := range r {
		dtoSlice[i] = v.ToDto()
	}
	return dtoSlice
}

type AddBookRequest struct {
	BookName        string `json:"book_name"`
	Author          string `json:"author"`
	PublicationYear int    `json:"publication_year"`
}
