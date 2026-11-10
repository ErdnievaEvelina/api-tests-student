package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
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
