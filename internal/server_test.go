package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/genwilliam/short_code/internal/short_link"
)

type memoryShortLinkStore struct {
	links map[string]*short_link.ShortLink
}

func (s *memoryShortLinkStore) Create(_ context.Context, params short_link.CreateShortLinkParams) (*short_link.ShortLink, error) {
	if _, exists := s.links[params.ShortCode]; exists {
		return nil, errors.New("duplicate short code")
	}
	link := &short_link.ShortLink{
		ShortCode:   params.ShortCode,
		OriginalURL: params.OriginalURL,
		Status:      params.Status,
		ExpiresAt:   params.ExpiresAt,
	}
	s.links[link.ShortCode] = link
	return link, nil
}

func (s *memoryShortLinkStore) ResolveByCode(_ context.Context, code string) (*short_link.ShortLink, error) {
	link, exists := s.links[code]
	if !exists || link.Status != short_link.StatusEnabled {
		return nil, short_link.ErrShortLinkNotFound
	}
	return link, nil
}

func TestCreateAndResolveShortLink(t *testing.T) {
	store := &memoryShortLinkStore{links: make(map[string]*short_link.ShortLink)}
	handler := NewHandler(store)

	createRequest := httptest.NewRequest(http.MethodPost, "/api/short-links", strings.NewReader(
		`{"short_code":"test1234","original_url":"https://example.com/path"}`,
	))
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d; body: %s", createResponse.Code, http.StatusCreated, createResponse.Body)
	}

	var response createShortLinkResponse
	if err := json.Unmarshal(createResponse.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if response.ShortURL != "http://example.com/test1234" {
		t.Fatalf("short_url = %q, want %q", response.ShortURL, "http://example.com/test1234")
	}

	resolveRequest := httptest.NewRequest(http.MethodGet, "/test1234", nil)
	resolveResponse := httptest.NewRecorder()
	handler.ServeHTTP(resolveResponse, resolveRequest)
	if resolveResponse.Code != http.StatusFound {
		t.Fatalf("resolve status = %d, want %d", resolveResponse.Code, http.StatusFound)
	}
	if location := resolveResponse.Header().Get("Location"); location != "https://example.com/path" {
		t.Fatalf("Location = %q, want %q", location, "https://example.com/path")
	}
}

func TestCreateShortLinkRejectsNonHTTPURL(t *testing.T) {
	store := &memoryShortLinkStore{links: make(map[string]*short_link.ShortLink)}
	handler := NewHandler(store)
	request := httptest.NewRequest(http.MethodPost, "/api/short-links", strings.NewReader(
		`{"original_url":"javascript:alert(1)"}`,
	))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
