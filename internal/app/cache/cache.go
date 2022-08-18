package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/pechenegi/backend/internal/pkg/models"
)

var (
	errOldStats = errors.New("cannot add old stats to cache")
)

type Cache interface {
	AddOrReplaceDebtStats(ctx context.Context, userID string, stats models.DebtStats) error
	GetDebtStatsIfExists(ctx context.Context, userID string) (models.DebtStats, bool)
}

type cache struct {
	sync.RWMutex
	store map[string]models.DebtStats
}

func InitCache(ctx context.Context) Cache {
	return &cache{
		store: make(map[string]models.DebtStats),
	}
}

func (c *cache) AddOrReplaceDebtStats(ctx context.Context, userID string, stats models.DebtStats) error {
	if err := validateStats(ctx, stats); err != nil {
		return err
	}
	c.Lock()
	c.store[userID] = stats
	c.Unlock()
	return nil
}

func (c *cache) GetDebtStatsIfExists(ctx context.Context, userID string) (models.DebtStats, bool) {
	c.Lock()
	defer c.Unlock()
	v, ok := c.store[userID]
	if !ok {
		return models.DebtStats{}, false
	}

	if err := validateStats(ctx, v); err != nil {
		c.deleteEntry(ctx, userID)
		return models.DebtStats{}, false
	}
	return v, ok
}

func (c *cache) deleteEntry(ctx context.Context, userID string) {
	delete(c.store, userID)
}

func validateStats(ctx context.Context, stats models.DebtStats) error {
	nowY, nowM, nowD := time.Now().Date()
	if calcY, calcM, calcD := stats.CalculatedAt.Date(); calcD != nowD || calcM != nowM || calcY != nowY {
		return errOldStats
	}
	return nil
}
