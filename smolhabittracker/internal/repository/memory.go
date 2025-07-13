package repository

import (
	"context"
	"fmt"
	"smol/smolhabittracker/internal/habit"
	"smol/smolhabittracker/internal/isoweek"
	"sort"
	"sync"
	"time"
)

type ticksPerWeek map[isoweek.ISO8601][]time.Time

type HabitRepository struct {
	mutex sync.Mutex
	lgr   Logger

	habits map[habit.ID]habit.Habit
	ticks  map[habit.ID]ticksPerWeek
}

type Repository interface {
	Add(ctx context.Context, habit habit.Habit)
	ListAll(ctx context.Context) ([]habit.Habit, error)
}

func New(lgr Logger) *HabitRepository {
	return &HabitRepository{
		lgr:    lgr,
		habits: make(map[habit.ID]habit.Habit),
		ticks:  make(map[habit.ID]ticksPerWeek),
	}
}

func (hr *HabitRepository) Add(_ context.Context, habit habit.Habit) error {
	hr.lgr.Logf("Adding a habit...")

	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.habits[habit.ID] = habit

	return nil
}

func (hr *HabitRepository) List(_ context.Context, id habit.ID) (habit.Habit, error) {
	hr.lgr.Logf("Finding a habit...")
	h, found := hr.habits[id]
	if !found {
		return habit.Habit{}, fmt.Errorf("habit %q not registered: %w", id, ErrNotFound)
	}

	return h, nil
}

func (hr *HabitRepository) ListAll(_ context.Context) ([]habit.Habit, error) {
	hr.lgr.Logf("Listing habits, sorted by creation time...")

	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	habits := make([]habit.Habit, 0)
	for _, h := range hr.habits {
		habits = append(habits, h)
	}

	sort.Slice(habits, func(i, j int) bool {
		return habits[i].CreationTime.Before(habits[j].CreationTime)
	})

	return habits, nil
}

func (hr *HabitRepository) AddTick(_ context.Context, id habit.ID, t time.Time) error {
	hr.lgr.Logf("Adding a tcik...")

	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	_, found := hr.ticks[id]
	if !found {
		hr.ticks[id] = make(ticksPerWeek)
	}

	w := isoweek.At(t)

	ticks, found := hr.ticks[id][w]
	if !found {
		ticks = make([]time.Time, 0, 1)
	}

	hr.ticks[id][w] = append(ticks, t)

	return nil
}

func (hr *HabitRepository) ListAllTicks(_ context.Context, id habit.ID) ([]time.Time, error) {
	hr.lgr.Logf("Listing ticks for a while...")

	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	ticks := make([]time.Time, 0)
	for _, weeklyTicks := range hr.ticks[id] {
		ticks = append(ticks, weeklyTicks...)
	}

	return ticks, nil
}

func (hr *HabitRepository) ListWeeklyTicks(_ context.Context, id habit.ID, t time.Time) ([]time.Time, error) {
	hr.lgr.Logf("Listing weekly ticks for a habit...")

	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	loggedWeeks, found := hr.ticks[id]
	if !found {
		if _, ok := hr.habits[id]; ok {
			return []time.Time{}, nil
		}
		return nil, fmt.Errorf("id %q not registered: %w", id, ErrNotFound)
	}

	w := isoweek.At(t)
	if loggedWeeks[w] == nil {
		return []time.Time{}, nil
	}

	return loggedWeeks[w], nil
}

type Logger interface {
	Logf(format string, args ...any)
}
