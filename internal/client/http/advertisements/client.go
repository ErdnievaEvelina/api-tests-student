package advertisements

import (
	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Query("limit", strconv.Itoa(50)).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpGetAdvertisementByID(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementGetPath).
		Query("id", id).
		Expect(t).
		End().Response
}
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	url := fmt.Sprintf(path.AdvertisementPhotosPath, id)
	return apiRunner.GetRunner().Auth(token).Create().Get(url).
		Expect(t).
		End().Response
}

func HttpPostAdvertisement(t *testing.T, token string, title, description string, price float64, photoPath string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementCreatePath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.FormatFloat(price, 'f', -1, 64)).
		MultipartFormData("quantity", "1").
		MultipartFile("photos", photoPath).
		Expect(t).
		End().Response
}
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementDeletePath).
		Query("id", id).
		Expect(t).
		End().Response
}
