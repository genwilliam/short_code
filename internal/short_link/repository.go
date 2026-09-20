package short_link

import (
	"context"
	"database/sql"
	"errors"
)

var ErrShortLinkNotFound = errors.New("short link not found")

type ShortLinkRepository struct {
	db *sql.DB
}

func NewShortLinkRepository(db *sql.DB) *ShortLinkRepository {
	return &ShortLinkRepository{
		db: db,
	}
}
func (r *ShortLinkRepository) Insert(
	ctx context.Context,
	link *ShortLink,
) error {
	const query = `
		INSERT INTO short_links (
			short_code,
			original_url,
			status,
			expires_at
		)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		link.ShortCode,
		link.OriginalURL,
		link.Status,
		link.ExpiresAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	link.ID = uint64(id)

	return nil
}
