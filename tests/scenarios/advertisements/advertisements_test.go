package advertisements

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	adClient "api-tests-template/internal/client/http/advertisements"
	adManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
	scenario "api-tests-template/tests/scenarios"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователя")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Создание объявления и проверка через GET /advertisement, поиск и фото
func (s *TestSuite) TestCreateAdvertisementPositive() {
	var adID string
	title := "Test Ad " + utils.RandomString(10)
	desc := scenario.TestDesc
	price := scenario.DefaultPrice
	photoPath := scenario.PhotoPath
	defer func() {
		if adID != "" {
			utils.LogWithLabelAndTimeStamp("TearDown", "Удаляем тестовое объявление "+adID)
			adManager.DeleteAdvertisement(s.T(), s.loginData.Token, adID)
		}
	}()

	s.Run("Создаем объявление", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Создание объявления с title="+title)
		resp := adManager.CreateAdvertisement(s.T(), s.loginData.Token, title, desc, price, photoPath)
		adID = resp.ID
		utils.LogWithLabelAndTimeStamp("Check", "Объявление создано с ID="+resp.ID)
		require.Equal(s.T(), title, resp.Title, "Title полученного объявления не совпадает")
		require.Equal(s.T(), desc, resp.Description, "Description полученного объявления не совпадает")
		require.Equal(s.T(), price, resp.Price, "Price полученного объявления не совпадает")
	})

	s.Run("Проверяем через GET /advertisement?id={id}", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Запрос объявления по ID="+adID)
		ad := adManager.GetAdvertisementByID(s.T(), s.loginData.Token, adID)
		utils.LogWithLabelAndTimeStamp("Check", "Проверка параметров полученного объявления")
		require.Equal(s.T(), adID, ad.ID, "ID полученного объявления не совпадает")
		require.Equal(s.T(), title, ad.Title, "Title полученного объявления не совпадает")
		require.Equal(s.T(), desc, ad.Description, "Description полученного объявления не совпадает")
		require.Equal(s.T(), price, ad.Price, "Price полученного объявления не совпадает")
	})

	s.Run("Проверяем, что объявление есть в поиске по названию", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Проверка, что объявление есть в поиске по названию")
		list := adManager.GetAdvertisementsList(s.T(), s.loginData.Token, title)

		var found bool
		for _, item := range list.Items {
			if item.ID == adID {
				found = true
				break
			}
		}
		require.True(s.T(), found, "Созданное объявление не найдено в поиске")
	})

	s.Run("Проверяем доступность фото через GET /advertisements/{id}/photos", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Проверка доступности фото")
		photos := adManager.GetAdvertisementPhotos(s.T(), s.loginData.Token, adID)
		require.NotEmpty(s.T(), photos, "Список фото не должен быть пустым, так как мы загрузили файл")
	})
}

// Negative: Создание без обязательных параметров и с некорректными
func (s *TestSuite) TestCreateAdvertisementNegative() {
	photoPath := scenario.PhotoPath

	s.Run("Пустой title", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), s.loginData.Token, "", "desc", 100.0, photoPath)
		require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode)

		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.NotEmpty(s.T(), errResp.Message)
	})

	s.Run("Отрицательная цена", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), s.loginData.Token, "Valid "+utils.RandomString(5), "desc", -10.0, photoPath)
		require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode)

		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.NotEmpty(s.T(), errResp.Message)
	})
}

// Negative: Запрос без авторизации
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	photoPath := scenario.PhotoPath

	s.Run("Без токена", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), "", "Test", "Desc", 100.0, photoPath)
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)

		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.Equal(s.T(), "unauthorized", errResp.Error)
	})

	s.Run("С неверным токеном", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), "invalid_token", "Test", "Desc", 100.0, photoPath)
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode)
	})
}
