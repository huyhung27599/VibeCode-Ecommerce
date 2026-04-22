package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type txKey struct{}

// WithTx runs fn in a transaction. The transactional *gorm.DB is attached to
// the context; repositories can pick it up via FromContext.
func WithTx(ctx context.Context, db *gorm.DB, fn func(ctx context.Context) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

// FromContext returns the transactional *gorm.DB if present, otherwise the
// provided fallback db.
func FromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return fallback
}

// IsNotFound reports whether err is gorm.ErrRecordNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
