package server

import (
	"context"
	"smol/smolhabittracker/api"
	"smol/smolhabittracker/internal/habit"
)

func (s *Server) ListHabits(ctx context.Context, _ *api.ListHabitsRequest) (*api.ListHabitsResponse, error) {
	s.lgr.Logf("ListHabists request received")
	habits, err := habit.ListHabits(ctx, s.db)
	if err != nil {
		return nil, err
	}

	return convertHabitsToAPI(habits), nil
}

func convertHabitsToAPI(habits []habit.Habit) *api.ListHabitsResponse {
	hts := make([]*api.Habit, len(habits))

	for i := range habits {
		hts[i] = &api.Habit{
			Id:              string(habits[i].ID),
			Name:            string(habits[i].Name),
			WeeklyFrequency: int32(habits[i].WeeklyFrequency),
		}
	}

	return &api.ListHabitsResponse{
		Habits: hts,
	}
}
