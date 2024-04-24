package domain

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/logger"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
	"os"
	"strings"
)

type UserRepositoryDb struct {
}

func NewUserRepositoryDb() UserRepositoryDb {
	return UserRepositoryDb{
	}
}

func (d UserRepositoryDb) IsEmailExists(ctx context.Context, email string) (*CreateUserResponse, *errs.AppError) {
	var (
		userDetail CreateUserResponse
	)

	var users = []CreateUserRequest{
		{"admin", "admin", "admin@gmail.com", "Admin@1234"},
		{"user1", "user", "user1@gmail.com", "User1@1234"},
	}

	for _, u := range users {
		if u.Email == email {
			password, err := utility.HashPassword(u.Password)
			if err != nil {
				logger.Error("password hashing failed"+ err.Error())
				return nil, errs.NewValidationError("password hashing failed"+ err.Error())
			}
			userDetail.Name = u.Name
			userDetail.Role = u.Role
			userDetail.Email = u.Email
			userDetail.Password = password
			return &userDetail, nil
		}
	}
	return nil, errs.NewNotFoundError("User not found")
}

func (d UserRepositoryDb) GetBookList(ctx context.Context, userType string) ([]*GetUserBookListResponse, *errs.AppError) {
	var books []*GetUserBookListResponse

	if userType == "admin" {
		books1, err := readBooksFromFile("regularUser.csv")
		if err != nil {
			return nil, err
		}
		books2, err := readBooksFromFile("adminUser.csv")
		if err != nil {
			return nil, err
		}
		books = append(books, books1...)
		books = append(books, books2...)
	} else {
		books1, err := readBooksFromFile("regularUser.csv")
		if err != nil {
			return nil, err
		}
		books = append(books, books1...)
	}

	return books, nil
}

func readBooksFromFile(filename string) ([]*GetUserBookListResponse, *errs.AppError) {
	var books []*GetUserBookListResponse

	file, err := os.Open(filename)
	if err != nil {
		logger.Error("Failed to open file:"+ err.Error())
		return nil, errs.NewUnexpectedError("Failed to open file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		logger.Error("Failed to read file:"+ err.Error())
		return nil, errs.NewUnexpectedError("Failed to read file")
	}

	for i, line := range lines {
		if i == 0 {
			continue // Skip header
		}
		if len(line) >= 3 {
			book := &GetUserBookListResponse{
				BookName:        line[0],
				Author:          line[1],
				PublicationYear: line[2],
			}
			books = append(books, book)
		}
	}
	return books, nil
}

func(n UserRepositoryDb) AddBook(ctx context.Context, req AddBookRequest) ([]*GetUserBookListResponse, *errs.AppError){
	if err := addBookToCSV("regularUser.csv", strings.TrimSpace(req.BookName), req.Author, req.PublicationYear); err != nil {
		return nil, errs.NewValidationError("Failed to add book to the file"+ err.Error())
	}
	return nil, nil
}

func addBookToCSV(filename, bookName, author string, publicationYear int) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Move the file pointer to the end of the file
	if _, err := file.Seek(0, os.SEEK_END); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing the new record to the CSV file
	if err := writer.Write([]string{bookName, author, fmt.Sprintf("%d", publicationYear)}); err != nil {
		return err
	}

	return nil
}

func (d UserRepositoryDb) DeleteBook(ctx context.Context, bookName string) ([]*GetUserBookListResponse, *errs.AppError) {
	err := deleteRowByBookName("regularUser.csv", bookName)
	if err != nil {
		return nil, errs.NewValidationError("error occurred while deleting")
	}
	return nil, nil
}

func deleteRowByBookName(filename string, bookName string) error {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV file here
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// case insensitive
	var rowIndexToRemove = -1
	bookNameLower := strings.ToLower(bookName)
	for i, record := range records {
		if len(record) >= 1 && strings.ToLower(record[0]) == bookNameLower {
			rowIndexToRemove = i
			break
		}
	}

	if rowIndexToRemove == -1 {
		return fmt.Errorf("book name '%s' not found", bookName)
	}

	// Remove the row
	records = append(records[:rowIndexToRemove], records[rowIndexToRemove+1:]...)

	// Write the updated records back to the CSV file
	file, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}