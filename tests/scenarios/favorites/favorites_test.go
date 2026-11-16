package favorites

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	favClient "api-tests-template/internal/client/http/favorites"
	adManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	scenario "api-tests-template/tests/scenarios"
)

type TestSuite struct {
	suite.Suite
	loginData   models.LoginOkResponse
	loginData2  models.LoginOkResponse
	foreignAdID string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя 1 (владелец избранного)")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Авторизация пользователя 2 (владелец объявления)")
	s.loginData2 = auth.Login(s.T(), os.Getenv("TEST_LOGIN_2"), os.Getenv("TEST_PASSWORD_2"))

	base.Precondition("Создание объявления от имени пользователя 2")
	photoPath := scenario.PhotoPath
	title := "Foreign Ad " + utils.RandomString(8)
	resp := adManager.CreateAdvertisement(s.T(), s.loginData2.Token, title, "desc", 1000.0, photoPath)
	s.foreignAdID = resp.ID
	require.NotEmpty(s.T(), s.foreignAdID, "Не удалось создать чужое объявление")

	base.Precondition("Очистка избранного от тестового объявления")
	favorites.DeleteFavorite(s.T(), s.loginData.Token, s.foreignAdID)
}

func (s *TestSuite) TearDownSuite() {
	if s.foreignAdID != "" {
		utils.LogWithLabelAndTimeStamp("TearDown", "Удаляем тестовое объявление "+s.foreignAdID)
		adManager.DeleteAdvertisement(s.T(), s.loginData2.Token, s.foreignAdID)
	}
	base.TearDownSuite()
}

// Positive: Полный цикл работы с избранным
func (s *TestSuite) TestFavoritesPositiveFlow() {
	favorites.DeleteFavorite(s.T(), s.loginData.Token, s.foreignAdID)
	s.Run("Добавляем объявление в избранное", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Добавление в избранное ad_id="+s.foreignAdID)
		favorites.AddFavorite(s.T(), s.loginData.Token, s.foreignAdID)
	})

	s.Run("Объявление появляется в списке избранного", func() {
		list := favorites.GetFavorites(s.T(), s.loginData.Token)
		var found bool
		for _, item := range list.Items {
			if item.Advertisement.ID == s.foreignAdID {
				found = true
				break
			}
		}
		require.True(s.T(), found, "Добавленное объявление отсутствует в избранном")
	})

	s.Run("Удаляем объявление из избранного — оно пропадает из списка", func() {
		favorites.DeleteFavorite(s.T(), s.loginData.Token, s.foreignAdID)
		list := favorites.GetFavorites(s.T(), s.loginData.Token)
		for _, item := range list.Items {
			require.NotEqual(s.T(), s.foreignAdID, item.Advertisement.ID, "Удалённое объявление осталось в избранном")
		}
	})

	s.Run("При удалении объявления через DELETE /advertisement/{id} оно исчезает из избранного", func() {
		favorites.AddFavorite(s.T(), s.loginData.Token, s.foreignAdID)
		adManager.DeleteAdvertisement(s.T(), s.loginData2.Token, s.foreignAdID)
		list := favorites.GetFavorites(s.T(), s.loginData.Token)
		for _, item := range list.Items {
			require.NotEqual(s.T(), s.foreignAdID, item.Advertisement.ID, "Удалённое объявление осталось в избранном")
		}
		s.foreignAdID = ""
	})
}

// Negative: Операции с несуществующим ad_id
func (s *TestSuite) TestFavoritesNonExistentAdID() {
	nonExistentID := "11111111-1111-1111-1111-111111111111"

	s.Run("Добавление в избранное с несуществующим advertisement_id", func() {
		raw := favClient.HttpPostFavorite(s.T(), s.loginData.Token, nonExistentID)
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusNotFound, raw.StatusCode, "Ожидался статус 404 для несуществующего объявления")
	})

	s.Run("Удаление из избранного с несуществующим ad_id", func() {
		raw := favClient.HttpDeleteFavorite(s.T(), s.loginData.Token, nonExistentID)
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusNotFound, raw.StatusCode, "Ожидался статус 404 для несуществующего избранного")
	})
}

// Negative: Запрос без авторизации
func (s *TestSuite) TestFavoritesUnauthorized() {
	s.Run("GET /favorites без токена", func() {
		raw := favClient.HttpGetFavorites(s.T(), "")
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
	})

	s.Run("POST /favorites без токена", func() {
		raw := favClient.HttpPostFavorite(s.T(), "", s.foreignAdID)
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
	})

	s.Run("DELETE /favorites/{ad_id} без токена", func() {
		raw := favClient.HttpDeleteFavorite(s.T(), "", "any-id")
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
	})
}

// Negative: Бизнес-ограничения POST /favorites
func (s *TestSuite) TestFavoritesBusinessConstraints() {
	s.Run("Повторное добавление того же объявления в избранное", func() {
		favorites.DeleteFavorite(s.T(), s.loginData.Token, s.foreignAdID)

		favorites.AddFavorite(s.T(), s.loginData.Token, s.foreignAdID)

		raw := favClient.HttpPostFavorite(s.T(), s.loginData.Token, s.foreignAdID)
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusConflict, raw.StatusCode, "Ожидался статус 409 при повторном добавлении")

		favorites.DeleteFavorite(s.T(), s.loginData.Token, s.foreignAdID)
	})
	s.Run("Попытка добавить в избранное своё собственное объявление", func() {
		photoPath := scenario.PhotoPath
		title := "My Own Ad " + utils.RandomString(8)
		resp := adManager.CreateAdvertisement(s.T(), s.loginData.Token, title, "desc", 1000.0, photoPath)
		ownAdID := resp.ID
		require.NotEmpty(s.T(), ownAdID)

		raw := favClient.HttpPostFavorite(s.T(), s.loginData.Token, ownAdID)
		defer raw.Body.Close()
		require.Equal(s.T(), http.StatusForbidden, raw.StatusCode, "Ожидался статус 403 при добавлении своего объявления")

		adManager.DeleteAdvertisement(s.T(), s.loginData.Token, ownAdID)
	})
}
