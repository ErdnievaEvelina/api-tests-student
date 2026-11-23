package advertisementSearch

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

const (
	existingAdTitle = "Моё тестовое объявление про велосипед"

	searchBeginning = "Моё"
	searchEnd       = "велосипед"
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

	base.Precondition("Проверить, что константы заданы")
	require.NotEmpty(s.T(), existingAdTitle, "задать existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "задать searchBeginning")
	require.NotEmpty(s.T(), searchEnd, "задать searchEnd")

	base.Precondition("Авторизация пользователя")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, searchBeginning)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)

		s.Run("Среди найденных есть наше объявление", func() {
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"конец", searchEnd},
	}

	for _, tc := range cases {
		tc := tc
		s.Run(tc.name, func() {
			list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
