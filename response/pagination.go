package response

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cobanhub/lib/errorx"
	"github.com/cobanhub/lib/httpx"
)

const (
	DefaultPage   = 1
	DefaultLimit  = 10
	MaxLimit      = 100
	DefaultSortBy = "created_at"
	DefaultOrder  = "desc"
)

type (
	Pagination struct {
		Page       int      `json:"page"`
		Limit      int      `json:"limit"`
		Offset     int      `json:"offset"`
		TotalData  int64    `json:"total_data"`
		TotalPages int      `json:"total_pages"`
		SortBy     string   `json:"sort_by"`
		Order      string   `json:"order"`
		Search     string   `json:"search"`
		Fields     []string `json:"fields"`
	}

	CursorPagination struct {
		Limit      int    `json:"limit"`
		Cursor     string `json:"cursor,omitempty"`
		NextCursor string `json:"next_cursor,omitempty"`
		HasNext    bool   `json:"has_next"`
		Refresh    bool   `json:"refresh,omitempty"`

		engine *CursorEngine
	}

	cursorPayload struct {
		Version   int       `json:"v"`
		ID        string    `json:"id"`
		SortValue any       `json:"sv"`
		IssuedAt  time.Time `json:"iat"`
	}

	CursorConfig struct {
		Secret []byte
		TTL    time.Duration
	}

	CursorEngine struct {
		aead cipher.AEAD
		ttl  time.Duration
	}
)

func NewCursorEngine(cfg CursorConfig) (*CursorEngine, error) {
	if len(cfg.Secret) != 16 && len(cfg.Secret) != 24 && len(cfg.Secret) != 32 {
		return nil, errors.New("cursor secret must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(cfg.Secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &CursorEngine{
		aead: gcm,
		ttl:  cfg.TTL,
	}, nil
}

func ParsePagination(req httpx.Request) *Pagination {
	p := &Pagination{
		Page:   DefaultPage,
		Limit:  DefaultLimit,
		Offset: 0,
		SortBy: DefaultSortBy,
		Order:  DefaultOrder,
		Search: "",
		Fields: []string{},
	}

	// page
	if v := req.Query("page"); v != "" {
		if page, err := strconv.Atoi(v); err == nil && page > 0 {
			p.Page = page
		}
	}

	// limit
	if v := req.Query("limit"); v != "" {
		if limit, err := strconv.Atoi(v); err == nil && limit > 0 {
			if limit > MaxLimit {
				limit = MaxLimit
			}
			p.Limit = limit
		}
	}

	// offset
	if v := req.Query("offset"); v != "" {
		if offset, err := strconv.Atoi(v); err == nil && offset >= 0 {
			p.Offset = offset
		}
	}

	// sort_by
	if v := req.Query("sort_by"); v != "" {
		p.SortBy = v
	}

	// order
	if v := strings.ToLower(req.Query("order")); v == "asc" || v == "desc" {
		p.Order = v
	}

	// search
	if v := req.Query("search"); v != "" {
		p.Search = v
	}

	// fields
	if v := req.Query("fields"); v != "" {
		for _, f := range strings.Split(v, ",") {
			if trimmed := strings.TrimSpace(f); trimmed != "" {
				p.Fields = append(p.Fields, trimmed)
			}
		}
	}

	// calculate offset if page is provided
	if p.Offset == 0 && p.Page > 1 {
		p.Offset = (p.Page - 1) * p.Limit
	}

	return p
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetLimit() int {
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return p.Offset
}

func (p *Pagination) GetSortBy() string {
	return p.SortBy
}

func (p *Pagination) GetOrder() string {
	return p.Order
}

func (p *Pagination) SetTotal(total int64) {
	p.TotalData = total
	if total == 0 {
		p.TotalPages = 0
		return
	}

	p.TotalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
}

func (e *CursorEngine) Parse(limit int, cursor string) *CursorPagination {
	if limit <= 0 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return &CursorPagination{
		Limit:  limit,
		Cursor: cursor,
		engine: e,
	}
}
func (p *CursorPagination) SetNext(id string, sortValue any) error {
	payload := cursorPayload{
		Version:   1,
		ID:        id,
		SortValue: sortValue,
		IssuedAt:  time.Now(),
	}

	enc, err := p.engine.encrypt(payload)
	if err != nil {
		return err
	}

	p.NextCursor = enc
	p.HasNext = true
	return nil
}

func (p *CursorPagination) Decode() (id string, sortValue any, ok bool, err error) {
	if p.Cursor == "" {
		return "", nil, false, nil
	}

	payload, err := p.engine.decrypt(p.Cursor)
	if err != nil {
		if errors.Is(err, errorx.ErrCursorExpired) {
			p.Refresh = true
			return "", nil, false, nil
		}
		return "", nil, false, err
	}

	return payload.ID, payload.SortValue, true, nil
}
func (e *CursorEngine) encrypt(payload cursorPayload) (string, error) {
	nonce := make([]byte, e.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	ciphertext := e.aead.Seal(nil, nonce, raw, nil)
	return base64.RawURLEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func (e *CursorEngine) decrypt(cursor string) (*cursorPayload, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, errorx.ErrCursorInvalid
	}

	if len(raw) < e.aead.NonceSize() {
		return nil, errorx.ErrCursorInvalid
	}

	nonce := raw[:e.aead.NonceSize()]
	data := raw[e.aead.NonceSize():]

	plain, err := e.aead.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, errorx.ErrCursorInvalid
	}

	var payload cursorPayload
	if err := json.Unmarshal(plain, &payload); err != nil {
		return nil, errorx.ErrCursorInvalid
	}

	if e.ttl > 0 && time.Since(payload.IssuedAt) > e.ttl {
		return nil, errorx.ErrCursorExpired
	}

	return &payload, nil
}
