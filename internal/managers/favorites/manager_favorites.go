package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

func AddFavorite(t *testing.T, token, adID string) {
	raw := favorites.HttpPostFavorite(t, token, adID)
	defer raw.Body.Close()

	require.Equal(t, http.StatusCreated, raw.StatusCode,
		"Ожидался статус 201 при добавлении в избранное")
}

func GetFavorites(t *testing.T, token string) models.FavoritesList {
	raw := favorites.HttpGetFavorites(t, token)
	defer raw.Body.Close()

	require.Equal(t, http.StatusOK, raw.StatusCode,
		"Ожидался статус 200 при получении избранного")

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	var resp models.FavoritesList
	require.NoError(t, json.Unmarshal(body, &resp))
	return resp
}

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
