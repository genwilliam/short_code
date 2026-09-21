package short_link

import "time"

const (
	StatusDisabled uint8 = 0
	StatusEnabled  uint8 = 1
)

type CreateShortLinkParams struct {
	ShortCode   string
	OriginalURL string
	Status      uint8
	ExpiresAt   *time.Time
}

type ShortLink struct {
	ID          uint64
	ShortCode   string
	OriginalURL string
	Status      uint8
	ExpiresAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
