package dto

import (
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (c CreateUserRequest) Validate() *errs.AppError {
	if stringUtils.IsBlank(c.Email) {
		return errs.NewValidationError("email cannot be empty")
	}
	if stringUtils.IsBlank(c.Password) {
		return errs.NewValidationError("password cannot be empty")
	} else {
		if ok := utility.IsStrongPassword(c.Password); !ok {
			return errs.NewValidationError("password should contain at least one uppercase character,one digit, one special character and length must be 8 above")
		}
	}
	return nil
}

type GetUserBookListResponse struct {
	BookName        string `json:"book_name"`
	Author          string `json:"author"`
	PublicationYear string `json:"publication_year"`
}

type AddBookRequest struct {
	BookName        string `json:"book_name"`
	Author          string `json:"author"`
	PublicationYear int `json:"publication_year"`
}

func (req AddBookRequest) Validate() *errs.AppError {
	if stringUtils.IsBlank(req.BookName) {
		return errs.NewValidationError("Book Name can't blank")
	}
	if stringUtils.IsBlank(req.Author) {
		return errs.NewValidationError("Author can't blank")
	}
	if req.PublicationYear == 0{
		return errs.NewValidationError("PublicationYear can't Zero")
	}
	return nil
}