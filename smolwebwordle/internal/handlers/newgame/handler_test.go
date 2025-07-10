package newgame

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"smol/smolwebwordle/internal/api"
	"smol/smolwebwordle/internal/session"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	idFinderRegexp := regexp.MustCompile(`.+"id":"([0-9A-Z]+)".+`)

	tt := map[string]struct {
		wantStatusCode int
		wantBody       string
		creator        gameAdder
	}{
		"nominal": {
			wantStatusCode: http.StatusCreated,
			wantBody:       `{"id":"123456","attempts_left":5,"guesses":[],"word_length":5,"status":"Playing"}`,
			creator: gameAdderStub{
				err: nil,
			},
		},
	}

	for name, testCase := range tt {
		t.Run(name, func(t *testing.T) {
			f := Handler(testCase.creator)

			req, err := http.NewRequest(http.MethodPost, api.NewGameRoute, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			f.ServeHTTP(rr, req)

			assert.Equal(t, testCase.wantStatusCode, rr.Code)

			if testCase.wantBody == "" {
				return
			}

			body := rr.Body.String()
			id := idFinderRegexp.FindStringSubmatch(body)
			if len(id) != 2 {
				t.Fatal("cannot find one single id in the json output")
			}
			body = strings.Replace(body, id[1], "123456", 1)

			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			assert.JSONEq(t, testCase.wantBody, body)
		})
	}
}

func Test_createGame(t *testing.T) {
	corpusPath := "testdata/corpus.txt"

	g, err := createGame(gameAdderStub{nil}, corpusPath)
	require.NoError(t, err)

	assert.Equal(t, uint8(5), g.AttemptsLeft)
	assert.Equal(t, 0, len(g.Guesses))
	assert.Regexp(t, "[A-Z0-9]+", g.ID)
	assert.Equal(t, "GAMER", g.SmolWordle.ShowAnswer())
}

type gameAdderStub struct {
	err error
}

func (g gameAdderStub) Add(session.Game) error {
	return g.err
}
