package advertisements

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	adClient "api-tests-template/internal/client/http/advertisements"
	adManager "api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
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
	var created adModels.CreateAdvertisementResponse
	title := "Test Ad " + utils.RandomString(10)
	desc := scenario.TestDesc
	price := scenario.DefaultPrice
	photoPath := scenario.PhotoPath
	qty := scenario.Qty
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
		created = resp
		utils.LogWithLabelAndTimeStamp("Check", "Объявление создано с ID="+resp.ID)
		require.Equal(s.T(), title, resp.Title, "Title полученного объявления не совпадает")
		require.Equal(s.T(), desc, resp.Description, "Description полученного объявления не совпадает")
		require.Equal(s.T(), price, resp.Price, "Price полученного объявления не совпадает")
		require.Equal(s.T(), qty, resp.Quantity, "Quantity полученного объявления не совпадает")
		require.NotEmpty(s.T(), resp.Photos, "В ответе создания должно быть фото")
	})

	s.Run("Проверяем через GET /advertisement?id={id}", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Запрос объявления по ID="+adID)
		ad := adManager.GetAdvertisementByID(s.T(), s.loginData.Token, adID)
		utils.LogWithLabelAndTimeStamp("Check", "Проверка параметров полученного объявления")
		require.Equal(s.T(), adID, ad.ID, "ID полученного объявления не совпадает")
		require.Equal(s.T(), title, ad.Title, "Title полученного объявления не совпадает")
		require.Equal(s.T(), desc, ad.Description, "Description полученного объявления не совпадает")
		require.Equal(s.T(), price, ad.Price, "Price полученного объявления не совпадает")
		require.Equal(s.T(), qty, ad.Quantity, "Quantity полученного объявления не совпадает")
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
		require.NotEmpty(s.T(), photos, "Список фото не должен быть пустым")
		ad := adManager.GetAdvertisementByID(s.T(), s.loginData.Token, adID)
		require.Equal(s.T(), len(created.Photos), len(photos),
			"Количество фото в ответе создания и GET /photos не совпадает")

		createdIDs := map[string]bool{}
		for _, p := range created.Photos {
			createdIDs[p.ID] = true
		}
		for _, p := range photos {
			require.True(s.T(), createdIDs[p.ID],
				"ID фото %s отсутствует в ответе создания", p.ID)
		}
		require.Equal(s.T(), len(ad.Photos), len(photos), "Количество фото в GET /advertisement и GET /photos не совпадает")
	})
}

// Negative: Создание без обязательных параметров и с некорректными
func (s *TestSuite) TestCreateAdvertisementNegative() {
	photoPath := scenario.PhotoPath

	s.Run("Пустой title", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), s.loginData.Token, "", "desc", 100.0, photoPath)
		require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode, "Ожидалась ошибка валидации (400)")

		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.Equal(s.T(), "title is required", errResp.Message, "Ожидается title is required")
		require.Equal(s.T(), "invalid_request", errResp.Error, "Ожидается invalid_request")
	})

	s.Run("Отрицательная цена", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), s.loginData.Token, "Valid "+utils.RandomString(5), "desc", -10.0, photoPath)
		require.Equal(s.T(), http.StatusBadRequest, raw.StatusCode, "Ожидалась ошибка валидации (400)")

		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.Equal(s.T(), "price must be between 0 and 1000000 inclusive", errResp.Message, "Ожидается price must be between 0 and 1000000 inclusive")
		require.Equal(s.T(), "invalid_request", errResp.Error, "Ожидается invalid_request")
	})
}

// Negative: Запрос без авторизации
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	photoPath := scenario.PhotoPath

	s.Run("Без токена", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), "", "Test", "Desc", 100.0, photoPath)
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode, "Ожидалась ошибка авторизации (401)")

		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.Equal(s.T(), "unauthorized", errResp.Error, "Ожидался unauthorized")
		require.Contains(s.T(), errResp.Message, "Authorization header required")
	})

	s.Run("С неверным токеном", func() {
		raw := adClient.HttpPostAdvertisement(s.T(), "invalid_token", "Test", "Desc", 100.0, photoPath)
		errResp := adManager.ParseErrorResponse(s.T(), raw)
		require.Equal(s.T(), http.StatusUnauthorized, raw.StatusCode, "Ожидалась ошибка авторизации (401)")
		require.Equal(s.T(), "unauthorized", errResp.Error, "Ожидался unauthorized")
		require.Contains(s.T(), errResp.Message, "Invalid or expired token", "Ожидался Invalid or expired token")
	})
}
