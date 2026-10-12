package internal

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/genwilliam/short_code/internal/short_link"
	"github.com/go-sql-driver/mysql"
)

type shortLinkStore interface {
	Create(context.Context, short_link.CreateShortLinkParams) (*short_link.ShortLink, error)
	ResolveByCode(context.Context, string) (*short_link.ShortLink, error)
}

var shortCodePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{4,32}$`)

type createShortLinkRequest struct {
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type createShortLinkResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

func NewHandler(store shortLinkStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok\n")
	})
	mux.HandleFunc("POST /api/short-links", func(w http.ResponseWriter, r *http.Request) {
		createShortLink(w, r, store)
	})
	mux.HandleFunc("GET /{code}", func(w http.ResponseWriter, r *http.Request) {
		resolveShortLink(w, r, store)
	})
	return mux
}

func createShortLink(w http.ResponseWriter, r *http.Request, store shortLinkStore) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request createShortLinkRequest
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request must contain one JSON object")
		return
	}

	request.OriginalURL = strings.TrimSpace(request.OriginalURL)
	parsedURL, err := url.ParseRequestURI(request.OriginalURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		writeError(w, http.StatusBadRequest, "original_url must be an absolute HTTP or HTTPS URL")
		return
	}
	if request.ExpiresAt != nil && !request.ExpiresAt.After(time.Now()) {
		writeError(w, http.StatusBadRequest, "expires_at must be in the future")
		return
	}

	code := strings.TrimSpace(request.ShortCode)
	if code == "" {
		code, err = generateShortCode()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not generate short code")
			return
		}
	} else if !shortCodePattern.MatchString(code) {
		writeError(w, http.StatusBadRequest, "short_code must be 4-32 letters, digits, hyphens, or underscores")
		return
	}

	link, err := store.Create(r.Context(), short_link.CreateShortLinkParams{
		ShortCode:   code,
		OriginalURL: request.OriginalURL,
		Status:      short_link.StatusEnabled,
		ExpiresAt:   request.ExpiresAt,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			writeError(w, http.StatusConflict, "short_code already exists")
			return
		}
		if errors.Is(err, short_link.ErrInvalidShortLink) {
			writeError(w, http.StatusBadRequest, "invalid short link")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create short link")
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	writeJSON(w, http.StatusCreated, createShortLinkResponse{
		ShortCode: link.ShortCode,
		ShortURL:  scheme + "://" + r.Host + "/" + url.PathEscape(link.ShortCode),
	})
}

func resolveShortLink(w http.ResponseWriter, r *http.Request, store shortLinkStore) {
	code := r.PathValue("code")
	if !shortCodePattern.MatchString(code) {
		writeError(w, http.StatusNotFound, "short link not found")
		return
	}

	link, err := store.ResolveByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, short_link.ErrShortLinkNotFound) {
			writeError(w, http.StatusNotFound, "short link not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not resolve short link")
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}

func generateShortCode() (string, error) {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 8)
	limit := big.NewInt(int64(len(alphabet)))
	for i := range code {
		index, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", fmt.Errorf("generate short code: %w", err)
		}
		code[i] = alphabet[index.Int64()]
	}
	return string(code), nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
