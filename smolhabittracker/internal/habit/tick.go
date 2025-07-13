package habit

import (
	"context"
	"fmt"
	"time"
)

type habitFinder interface {
	List(ctx context.Context, id ID) (Habit, error)
}

type tickAdder interface {
	AddTick(ctx context.Context, id ID, t time.Time) error
}

func Tick(ctx context.Context, habitDB habitFinder, tickDB tickAdder, id ID, t time.Time) error {
	_, err := habitDB.List(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot find habit %q: %w", id, err)
	}

	err = tickDB.AddTick(ctx, id, t)
	if err != nil {
		return fmt.Errorf("cannot insert tick for habit %q: %w", id, err)
	}

	return nil
}
