package short_link

import "time"

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
