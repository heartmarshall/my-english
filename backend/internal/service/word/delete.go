package word

import (
	"context"
	"errors"

	"github.com/heartmarshall/my-english/internal/database"
	"github.com/heartmarshall/my-english/internal/service"
)

// Delete удаляет слово и все связанные данные.
// Связанные данные удаляются каскадно на уровне БД.
func (s *Service) Delete(ctx context.Context, id int64) error {
	err := s.words.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return service.ErrWordNotFound
		}
		return err
	}

	return nil
}
