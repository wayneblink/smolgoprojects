package solver

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolver_explore(t *testing.T) {
	tests := map[string]struct {
		inputImage string
		wantSize   int
	}{
		"cross": {
			inputImage: "testdata/explore_cross.png",
			wantSize:   2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			maze, err := openMaze(tt.inputImage)
			require.NoError(t, err)

			s := &Solver{
				maze:           maze,
				palette:        defaultPalette(),
				pathsToExplore: make(chan *path, 3),
			}

			s.explore(&path{at: image.Point{0, 2}})

			assert.Equal(t, tt.wantSize, len(s.pathsToExplore))
		})
	}
}
