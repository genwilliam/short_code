package short_link

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrShortLinkNotFound = errors.New("short link not found")
	ErrInvalidShortLink  = errors.New("invalid short link")
)

type ShortLinkRepository struct {
	db *sql.DB
}

func NewShortLinkRepository(db *sql.DB) *ShortLinkRepository {
	return &ShortLinkRepository{
		db: db,
	}
}

func (r *ShortLinkRepository) Create(
	ctx context.Context,
	params CreateShortLinkParams,
) (*ShortLink, error) {
	if strings.TrimSpace(params.ShortCode) == "" ||
		strings.TrimSpace(params.OriginalURL) == "" {
		return nil, ErrInvalidShortLink
	}

	if params.Status != StatusDisabled && params.Status != StatusEnabled {
		return nil, ErrInvalidShortLink
	}

	link := &ShortLink{
		ShortCode:   params.ShortCode,
		OriginalURL: params.OriginalURL,
		Status:      params.Status,
		ExpiresAt:   params.ExpiresAt,
	}

	if err := r.Insert(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (r *ShortLinkRepository) Insert(
	ctx context.Context,
	link *ShortLink,
) error {
	if link == nil || strings.TrimSpace(link.ShortCode) == "" ||
		strings.TrimSpace(link.OriginalURL) == "" ||
		(link.Status != StatusDisabled && link.Status != StatusEnabled) {
		return ErrInvalidShortLink
	}

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

func (r *ShortLinkRepository) UpdateStatus(
	ctx context.Context,
	shortCode string,
	status uint8,
) error {
	if strings.TrimSpace(shortCode) == "" ||
		(status != StatusDisabled && status != StatusEnabled) {
		return ErrInvalidShortLink
	}

	const query = `
		UPDATE short_links
		SET status = ?
		WHERE short_code = ?
		  AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, status, shortCode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrShortLinkNotFound
	}

	return nil
}

func (r *ShortLinkRepository) Delete(
	ctx context.Context,
	shortCode string,
) error {
	if strings.TrimSpace(shortCode) == "" {
		return ErrInvalidShortLink
	}

	const query = `
		UPDATE short_links
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE short_code = ?
		  AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrShortLinkNotFound
	}

	return nil
}

func (r *ShortLinkRepository) FindByCode(
	ctx context.Context,
	shortCode string,
) (*ShortLink, error) {

	const query = `
		SELECT
			id,
			short_code,
			original_url,
			status,
			expires_at,
			created_at,
			updated_at,
			deleted_at
		FROM short_links
		WHERE short_code = ?
		  AND deleted_at IS NULL
	`

	var link ShortLink

	err := r.db.QueryRowContext(
		ctx,
		query,
		shortCode,
	).Scan(
		&link.ID,
		&link.ShortCode,
		&link.OriginalURL,
		&link.Status,
		&link.ExpiresAt,
		&link.CreatedAt,
		&link.UpdatedAt,
		&link.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrShortLinkNotFound
		}

		return nil, err
	}

	return &link, nil
}
