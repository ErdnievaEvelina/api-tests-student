package advertisements

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// GetAdvertisementsList возвращает список объявлений по текстовому поиску.
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)

	var list models.AdvertisementsList
	utils.DecodeJSONResponse(t, raw, &list)
	return list
}

// CreateAdvertisement создает объявление и возвращает ответ с его данными.
func CreateAdvertisement(t *testing.T, token, title, desc string, price float64, photoPath string) models.CreateAdvertisementResponse {
	raw := advertisements.HttpPostAdvertisement(t, token, title, desc, price, photoPath)

	require.Equal(t, http.StatusCreated, raw.StatusCode, "Ожидался статус 201 при создании объявления")

	var resp models.CreateAdvertisementResponse
	utils.DecodeJSONResponse(t, raw, &resp)
	require.NotEmpty(t, resp.ID, "ID созданного объявления пуст")

	return resp
}

// GetAdvertisementByID возвращает детальную информацию об объявлении по его ID.
func GetAdvertisementByID(t *testing.T, token, id string) models.Advertisement {
	raw := advertisements.HttpGetAdvertisementByID(t, token, id)

	require.Equal(t, http.StatusOK, raw.StatusCode, "Ожидался статус 200 при получении объявления")

	var resp models.Advertisement
	utils.DecodeJSONResponse(t, raw, &resp)
	return resp
}

// GetAdvertisementPhotos возвращает список фотографий для указанного объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)

	require.Equal(t, http.StatusOK, raw.StatusCode, "Ожидался статус 200 при получении фото")

	var resp []models.Photo
	utils.DecodeJSONResponse(t, raw, &resp)
	return resp
}

// ParseErrorResponse читает и парсит тело ответа в структуру ErrorResponse.
func ParseErrorResponse(t *testing.T, raw *http.Response) models.ErrorResponse {
	var errResp models.ErrorResponse
	utils.DecodeJSONResponse(t, raw, &errResp)
	return errResp
}

// DeleteAdvertisement удаляет объявление по ID.
func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()

	require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, raw.StatusCode,
		"Ожидался успешный статус 204 при удалении объявления")
}
