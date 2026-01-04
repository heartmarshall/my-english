package study

import "github.com/heartmarshall/my-english/internal/model"

// Filter используется для фильтрации очереди изучения
type Filter struct {
	Tags     []string
	Statuses []model.LearningStatus
}
