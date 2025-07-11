package solver

import (
	"image"
	plt "image/color/palette"

	"golang.org/x/image/draw"
)

func (s *Solver) countExplorablePixels() int {
	explorablePixels := 0
	for row := s.maze.Bounds().Min.Y; row < s.maze.Bounds().Max.Y; row++ {
		for col := s.maze.Bounds().Min.X; col < s.maze.Bounds().Max.X; col++ {
			if s.maze.RGBAAt(col, row) != s.palette.wall {
				explorablePixels++
			}
		}
	}

	return explorablePixels
}

func (s *Solver) registerExplorablePixels() {
	const totalExpectedFrames = 30

	explorablePixels := s.countExplorablePixels()
	pixelsExplored := 0

	for {
		select {
		case <-s.quit:
			return
		case pos := <-s.exploredPixels:
			s.maze.Set(pos.X, pos.Y, s.palette.explored)
			pixelsExplored++
			if pixelsExplored%(explorablePixels/totalExpectedFrames) == 0 {
				s.drawCurrentFrameToGIF()
			}
		}
	}
}

func (s *Solver) drawCurrentFrameToGIF() {
	const (
		gifWidth      = 500
		frameDuration = 20
	)

	frame := image.NewPaletted(image.Rect(0, 0, gifWidth, gifWidth*s.maze.Bounds().Dy()/s.maze.Bounds().Dx()), plt.Plan9)

	draw.NearestNeighbor.Scale(frame, frame.Rect, s.maze, s.maze.Bounds(), draw.Over, nil)

	s.animation.Image = append(s.animation.Image, frame)
	s.animation.Delay = append(s.animation.Delay, frameDuration)
}

func (s *Solver) writeLastFrame() {
	stepsFromTreasure := s.solution

	for stepsFromTreasure != nil {
		s.maze.Set(stepsFromTreasure.at.X, stepsFromTreasure.at.Y, s.palette.solution)
		stepsFromTreasure = stepsFromTreasure.previousStep
	}

	const solutionFrameDuration = 300
	s.drawCurrentFrameToGIF()
	s.animation.Delay[len(s.animation.Delay)-1] = solutionFrameDuration
}
