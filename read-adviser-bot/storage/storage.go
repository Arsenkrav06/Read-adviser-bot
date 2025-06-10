package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"read-adviser-bot/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("No saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := h.Write([]byte(p.URL)); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := h.Write([]byte(p.UserName)); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
