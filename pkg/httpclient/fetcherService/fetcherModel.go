package fetcherService

import "github.com/dibyendu/Authentication-Authorization/pkg/dto"

type GetUserBookListResponse struct{
	BookName string `json:"book_name"`
	Author string `json:"author"`
	PublicationYear string `json:"publication_year"`
}
func ToDtoSlice(fetcherData []*GetUserBookListResponse) []*dto.GetUserBookListResponse {
	var dtoData []*dto.GetUserBookListResponse

	for _, item := range fetcherData {
		dtoItem := &dto.GetUserBookListResponse{
			BookName:        item.BookName,
			Author:          item.Author,
			PublicationYear: item.PublicationYear,
		}
		dtoData = append(dtoData, dtoItem)
	}
	return dtoData
}
