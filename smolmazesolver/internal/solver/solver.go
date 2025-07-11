package solver

import (
	"fmt"
	"image"
	"image/gif"
	"log"
	"sync"
)

type Solver struct {
	maze           *image.RGBA
	palette        palette
	pathsToExplore chan *path
	quit           chan struct{}

	solution *path
	mutex    sync.Mutex

	exploredPixels chan image.Point
	animation      *gif.GIF
}

func New(imagePath string) (*Solver, error) {
	img, err := openMaze(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open maze image: %w", err)
	}

	return &Solver{
		maze:           img,
		palette:        defaultPalette(),
		pathsToExplore: make(chan *path, 1),
		quit:           make(chan struct{}),
		exploredPixels: make(chan image.Point),
		animation:      &gif.GIF{},
	}, nil
}

func (s *Solver) Solve() error {
	entrance, err := s.findEntrance()
	if err != nil {
		return fmt.Errorf("unable to find entrance: %w", err)
	}

	log.Printf("starting at %v", entrance)
	s.pathsToExplore <- &path{previousStep: nil, at: entrance}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.registerExplorablePixels()
	}()

	go func() {
		defer wg.Done()
		s.listenToBranches()
	}()

	wg.Wait()
	s.writeLastFrame()

	return nil
}

func (s *Solver) findEntrance() (image.Point, error) {
	for row := s.maze.Bounds().Min.Y; row < s.maze.Bounds().Max.Y; row++ {
		for col := s.maze.Bounds().Min.X; col < s.maze.Bounds().Max.X; col++ {
			if s.maze.RGBAAt(col, row) == s.palette.entrance {
				return image.Point{X: col, Y: row}, nil
			}
		}
	}

	return image.Point{}, fmt.Errorf("entrance position not found")
}
