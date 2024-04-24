package fetcherService

import (
	"encoding/json"
	"github.com/dibyendu/Authentication-Authorization/lib/constants"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/logger"
	"net/http"
)

//func GetS3Details(userType string) ([]*GetUserBookListResponse, *errs.AppError) {
//	var (
//		url = "http://localhost:8080//user/home"
//		client = &http.Client{}
//		data   *GetUserBookListResponse
//	)
//
//	payloadBuf := new(bytes.Buffer)
//	err := json.NewEncoder(payloadBuf).Encode(nil)
//	if err != nil {
//		logger.Error("GetS3Details payload error: " + err.Error())
//		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
//	}
//
//	req1, err := http.NewRequest(http.MethodGet, url, payloadBuf)
//	if err != nil {
//		logger.Error("GetS3Details call err: " + err.Error())
//		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
//	}
//	req1.Header.Add("Content-Type", "application/json")
//	res, err := client.Do(req1)
//	if err != nil {
//		logger.Error("GetS3Details call err: " + err.Error())
//		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
//	}
//	defer res.Body.Close()
//
//	err = json.NewDecoder(res.Body).Decode(&data)
//	if err != nil {
//		logger.Error("GetS3Details call err: " + err.Error())
//		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
//	}
//	logger.Debug("GetS3Details http call status:" + res.Status)
//	if res.StatusCode != http.StatusOK {
//		logger.Error("GetS3Details scan call err: " + err.Error())
//		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
//	}
//	return data, nil
//}

func GetDetailFromHome(authToken string) ([]*GetUserBookListResponse, *errs.AppError) {
	var (
		url    = "http://localhost:8080/user/home" // Fix the URL format
		client = &http.Client{}
		data   []*GetUserBookListResponse // Correct the data type
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("GetDetailFromHome: error creating request: " + err.Error())
		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}
	req.Header.Add("Authorization", authToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("GetDetailFromHome: error making request: " + err.Error())
		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Error("GetDetailFromHome: non-200 status code: " + res.Status)
		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		logger.Error("GetDetailFromHome: error decoding response: " + err.Error())
		return nil, errs.NewUnexpectedError(constants.UNEXPECTED_ERROR)
	}

	return data, nil
}


