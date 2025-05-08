package database

import (
	"errors"
	"fmt"

	"github.com/Ilya-c4talyst/go_calculator/service/models"
	"gorm.io/gorm"
)

// Создание выражения
func AddExpression(expression *models.Expression) error {
	result := DB.Create(&expression)
	if result.Error != nil {
		return errors.New("не удалось сохранить выражение")
	}
	return nil
}

// Получение списка выражений
func GetExpressions(userID uint32) ([]*models.Expression, error) {
	var expressions []*models.Expression

	// Ищем выражения только для указанного пользователя
	result := DB.Where("user_id = ?", userID).Find(&expressions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Возвращаем пустой слайс, если записей нет
			return []*models.Expression{}, nil
		}
		return nil, errors.New("не удалось получить список выражений")
	}

	return expressions, nil
}

// Получение выражения по ID
func GetExpressionbyID(userID uint32, id int) (*models.Expression, error) {
	var expression models.Expression

	// Ищем выражения только для указанного пользователя
	result := DB.Where("user_id = ? AND Id = ?", userID, id).Find(&expression)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &models.Expression{}, nil
		}
		return nil, errors.New("не удалось получить выражение")
	}

	return &expression, nil
}

// Получение выражения по ID
func UpdateExpression(updatedExpr *models.Expression) error {
	// Сначала проверяем существование выражения у этого пользователя
	var existingExpr models.Expression
	result := DB.Where("user_id = ? AND id = ?", updatedExpr.User_id, updatedExpr.Id).First(&existingExpr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("выражение не найдено или доступ запрещен")
		}
		return fmt.Errorf("ошибка при поиске выражения: %w", result.Error)
	}

	// Обновляем только разрешенные поля
	result = DB.Model(&existingExpr).Updates(models.Expression{
		Status: updatedExpr.Status,
		Result: updatedExpr.Result,
	})

	if result.Error != nil {
		return fmt.Errorf("ошибка при обновлении выражения: %w", result.Error)
	}

	return nil
}
