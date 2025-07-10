package guess

import (
	"net/http"
	"net/http/httptest"
	"smol/smolwebwordle/internal/api"
	"smol/smolwebwordle/internal/session"
	"smol/smolwebwordle/internal/smolwordle"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	game, _ := smolwordle.New([]string{"pocket"})
	handle := Handler(sucessGameGuesserStub{session.Game{
		ID:           "123456",
		SmolWordle:   *game,
		AttemptsLeft: 5,
		Guesses:      nil,
		Status:       session.StatusPlaying,
	}})

	req, err := http.NewRequest(http.MethodPut, "/games/123456", strings.NewReader(`{"guess":"pocket"}`))
	require.NoError(t, err)

	req.SetPathValue(api.GameID, "123456")

	recorder := httptest.NewRecorder()

	handle(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"id":"123456","attempts_left":0,"guesses":null,"word_length":0,"status":""}`, recorder.Body.String())
}

type sucessGameGuesserStub struct {
	game session.Game
}

func (g sucessGameGuesserStub) Find(id session.GameID) (session.Game, error) {
	return g.game, nil
}

func (g sucessGameGuesserStub) Update(game session.Game) error {
	return nil
}
