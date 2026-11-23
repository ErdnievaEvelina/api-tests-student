package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// DecodeJSONResponse читает тело ответа и десериализует его в переданный объект.
// Функция автоматически закрывает Body ответа.
func DecodeJSONResponse(t *testing.T, raw *http.Response, target interface{}) {
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	err = json.Unmarshal(body, target)
	require.NoError(t, err, "Не удалось распарсить JSON из тела ответа")
}
