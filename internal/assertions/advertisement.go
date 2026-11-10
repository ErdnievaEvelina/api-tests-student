package assertions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, expectedTitle string) {
	for _, item := range items {
		if item.Title == expectedTitle {
			return
		}
	}
	assert.Failf(t, "Объявление не найдено", "Ожидали найти объявление с title='%s' среди %d элементов", expectedTitle, len(items))
}
