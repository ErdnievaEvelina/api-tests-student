package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

// AddFavorite добавляет объявление в избранное пользователя.
func AddFavorite(t *testing.T, token, adID string) {
	raw := favorites.HttpPostFavorite(t, token, adID)
	defer raw.Body.Close()

	require.Equal(t, http.StatusCreated, raw.StatusCode,
		"Ожидался статус 201 при добавлении в избранное")
}

// GetFavorites возвращает список избранных объявлений пользователя.
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	raw := favorites.HttpGetFavorites(t, token)

	require.Equal(t, http.StatusOK, raw.StatusCode,
		"Ожидался статус 200 при получении избранного")

	var resp models.FavoritesList
	utils.DecodeJSONResponse(t, raw, &resp)
	return resp
}

// DeleteFavorite удаляет объявление из избранного.
func DeleteFavorite(t *testing.T, token, adID string) {
	raw := favorites.HttpDeleteFavorite(t, token, adID)
	defer raw.Body.Close()

	if raw.StatusCode == http.StatusNotFound {
		utils.LogWithLabelAndTimeStamp("TearDown", "Избранное для объявления "+adID+" уже удалено (404)")
		return
	}

	require.Equal(t, http.StatusNoContent, raw.StatusCode,
		"Ожидался статус 204 при удалении из избранного")
}
