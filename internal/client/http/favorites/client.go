package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

func HttpPostFavorite(t *testing.T, token, adID string) *http.Response {
	req := models.AddFavoriteRequest{AdvertisementID: adID}
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(req).
		Expect(t).
		End().Response
}

func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token, adID string) *http.Response {
	url := fmt.Sprintf(path.FavoriteDeletePath, adID)
	return apiRunner.GetRunner().Auth(token).Create().Delete(url).
		Expect(t).
		End().Response
}
