package cache

import (
	"context"
	"testing"
	"time"

	"github.com/pechenegi/backend/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestInitCache(t *testing.T) {
	t.Run("return initialized cache", func(t *testing.T) {
		_ = InitCache(context.Background())
	})
}

func TestGetDebtStatsIfExists(t *testing.T) {
	t.Run("return existing value", func(t *testing.T) {
		c := &cache{
			store: make(map[string]models.DebtStats),
		}
		userID := "1-2-3-4"
		stats := generateStats()
		c.store[userID] = stats

		actualStats, ok := c.GetDebtStatsIfExists(context.Background(), userID)
		assert.Equal(t, true, ok)
		assert.Equal(t, stats, actualStats)
	})

	t.Run("return not ok for non-existing value", func(t *testing.T) {
		c := &cache{
			store: make(map[string]models.DebtStats),
		}

		_, ok := c.GetDebtStatsIfExists(context.Background(), "1-2-3-4")
		assert.Equal(t, false, ok)
	})

	t.Run("return not ok for yesterday value and remove it", func(t *testing.T) {
		c := &cache{
			store: make(map[string]models.DebtStats),
		}
		userID := "1-2-3-4"
		stats := generateStats()
		stats.CalculatedAt = stats.CalculatedAt.Add(-1 * 24 * time.Hour)
		c.store[userID] = stats

		_, ok := c.GetDebtStatsIfExists(context.Background(), userID)
		assert.Equal(t, false, ok)

		_, ok = c.store[userID]
		assert.Equal(t, false, ok)
	})
}

func TestAddDebtStats(t *testing.T) {
	t.Run("add new stats", func(t *testing.T) {
		c := &cache{
			store: make(map[string]models.DebtStats),
		}
		userID := "1-2-3-4"
		stats := generateStats()
		err := c.AddOrReplaceDebtStats(context.Background(), userID, stats)
		assert.NoError(t, err)
	})

	t.Run("return err when trying to add old stats", func(t *testing.T) {
		c := &cache{
			store: make(map[string]models.DebtStats),
		}
		userID := "1-2-3-4"
		stats := generateStats()
		stats.CalculatedAt = stats.CalculatedAt.Add(-1 * 24 * time.Hour)
		err := c.AddOrReplaceDebtStats(context.Background(), userID, stats)
		assert.EqualError(t, err, errOldStats.Error())
	})
}

func generateStats() models.DebtStats {
	return models.DebtStats{
		StudyLoan: models.LoanStat{
			Loan:       690000.01,
			Paid:       154011.15,
			DaysTotal:  1095,
			DaysPassed: 286,
		},
		StartUpLoan: models.LoanStat{
			Loan:       100000.00,
			Paid:       23511.05,
			DaysTotal:  710,
			DaysPassed: 286,
		},
		CalculatedAt: time.Now(),
	}
}
