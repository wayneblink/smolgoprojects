package habit

import (
	"context"
	"fmt"
	"time"
)

type tickFinder interface {
	ListWeeklyTicks(ctx context.Context, id ID, t time.Time) ([]time.Time, error)
}

func GetStatus(ctx context.Context, habitDB habitFinder, tickDB tickFinder, id ID, t time.Time) (Habit, int, error) {
	h, err := habitDB.List(ctx, id)
	if err != nil {
		return Habit{}, 0, fmt.Errorf("cannot find habit %s: %w", id, err)
	}

	ticks, err := tickDB.ListWeeklyTicks(ctx, id, t)
	if err != nil {
		return Habit{}, 0, fmt.Errorf("cannot find weekly ticks for habit %q: %w", id, err)
	}

	return h, len(ticks), nil
}
