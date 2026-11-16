package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"net/http"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	rawResponse := advertisements.HttpGetAdvertisements(t, token, text)
	defer rawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(rawResponse.Body)
	require.NoError(t, errResponse)

	var list models.AdvertisementsList
	err := json.Unmarshal(respArray, &list)
	require.NoError(t, err)

	return list
}
func CreateAdvertisement(t *testing.T, token, title, desc string, price float64, photoPath string) models.CreateAdvertisementResponse {
	raw := advertisements.HttpPostAdvertisement(t, token, title, desc, price, photoPath)
	defer raw.Body.Close()

	require.Equal(t, http.StatusCreated, raw.StatusCode, "Ожидался статус 201 при создании объявления")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.CreateAdvertisementResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	require.NotEmpty(t, resp.ID, "ID созданного объявления пуст")

	return resp
}
func GetAdvertisementByID(t *testing.T, token, id string) models.Advertisement {
	raw := advertisements.HttpGetAdvertisementByID(t, token, id)
	defer raw.Body.Close()

	require.Equal(t, http.StatusOK, raw.StatusCode, "Ожидался статус 200 при получении объявления")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.Advertisement
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	require.Equal(t, http.StatusOK, raw.StatusCode, "Ожидался статус 200 при получении фото")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp []models.Photo
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func ParseErrorResponse(t *testing.T, raw *http.Response) models.ErrorResponse {
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var errResp models.ErrorResponse
	require.NoError(t, json.Unmarshal(body, &errResp))
	return errResp
}
func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()
	require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, raw.StatusCode,
		"Ожидался успешный статус 204 при удалении объявления. Тело ответа: %s")
}
