package server

import (
	"context"
	"net"
	"smol/smolhabittracker/api"
	"smol/smolhabittracker/internal/repository"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	grpcServer := newServer(t)
	listener, err := net.Listen("tcp", "")
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = grpcServer.Serve(listener)
		require.NoError(t, err)
	}()
	defer func() {
		grpcServer.Stop()
		wg.Wait()
	}()

	habitsClient, err := newClient(t, listener.Addr().String())
	require.NoError(t, err)

	idWalk := addHabit(t, habitsClient, nil, "walk in the forest")

	isRead := addHabit(t, habitsClient, ptr(3), "read a few pages")

	addHabit(t, habitsClient, nil, "walk in the forest")

	addHabit(t, habitsClient, ptr(3), "read a few pages")

	addHabitWithError(t, habitsClient, 5, "        ", codes.InvalidArgument)

	listHabitsMatches(t, habitsClient, []*api.Habit{
		{
			Name:            "walk in the forest",
			WeeklyFrequency: 1,
		},
		{
			Name:            "read a few pages",
			WeeklyFrequency: 3,
		},
	})

	tickHabit(t, habitsClient, idWalk)
	tickHabit(t, habitsClient, idWalk)

	tickHabit(t, habitsClient, isRead)

	getHabitStatusMatches(t, habitsClient, idWalk, &api.GetHabitStatusResponse{
		Habit: &api.Habit{
			Id:              idWalk,
			Name:            "walk in the forest",
			WeeklyFrequency: 1,
		},
		TicksCount: 2,
	})

	getHabitStatusMatches(t, habitsClient, isRead, &api.GetHabitStatusResponse{
		Habit: &api.Habit{
			Id:              isRead,
			Name:            "read a few pages",
			WeeklyFrequency: 3,
		},
		TicksCount: 1,
	})
}

func newServer(t *testing.T) *grpc.Server {
	t.Helper()
	s := New(repository.New(t), t)

	return s.registerGRPCServer()
}

func newClient(t *testing.T, serverAddress string) (api.HabitsClient, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(serverAddress, creds)
	require.NoError(t, err)

	return api.NewHabitsClient(conn), nil
}

func addHabit(t *testing.T, habitsClient api.HabitsClient, freq *int32, name string) string {
	resp, err := habitsClient.CreateHabit(context.Background(), &api.CreateHabitRequest{
		Name:            name,
		WeeklyFrequency: freq,
	})
	require.NoError(t, err)

	return resp.Habit.Id
}

func ptr(i int32) *int32 {
	return &i
}

func addHabitWithError(t *testing.T, habitsClient api.HabitsClient, freq int32, name string, statusCode codes.Code) {
	_, err := habitsClient.CreateHabit(context.Background(), &api.CreateHabitRequest{
		Name:            name,
		WeeklyFrequency: &freq,
	})
	statusErr, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, statusCode, statusErr.Code())
}

func listHabitsMatches(t *testing.T, habitsClient api.HabitsClient, expected []*api.Habit) {
	list, err := habitsClient.ListHabits(context.Background(), &api.ListHabitsRequest{})
	require.NoError(t, err)

	for i := range list.Habits {
		assert.NotEqual(t, "", list.Habits[i].Id)
		list.Habits[i].Id = ""
	}
	assert.Equal(t, list.Habits, expected)
}

func tickHabit(t *testing.T, habitsClient api.HabitsClient, id string) {
	_, err := habitsClient.TickHabit(context.Background(), &api.TickHabitRequest{
		HabitId: id,
	})
	require.NoError(t, err)
}

func getHabitStatusMatches(t *testing.T, habitsCliient api.HabitsClient, id string, expected *api.GetHabitStatusResponse) {
	h, err := habitsCliient.GetHabitStatus(context.Background(), &api.GetHabitStatusRequest{HabitId: id})
	require.NoError(t, err)

	assert.Equal(t, expected.Habit, h.Habit)
	assert.Equal(t, expected.TicksCount, h.TicksCount)
}
