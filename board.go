package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/nsf/termbox-go"
)

const (
	BoardHorizontalScale = 2
	BoardWidthUnits      = 10
	BoardWidth           = BoardHorizontalScale * (BoardWidthUnits + 2)

	BoardVerticalScale = 1
	BoardHeightUnits   = 25
	BoardHeight        = 25
)

type Board struct {
	Width  int
	Height int

	Components       []*Component
	CurrentComponent *Component

	CompletedLines int

	IsPlaced   bool
	IsGameOver bool

	HandlerMutex *sync.Mutex
}

func NewBoard(width, height int) *Board {
	return &Board{width, height, nil, nil, 0, false, false, new(sync.Mutex)}
}

func (b *Board) HasPlacedBlock(x, y int) bool {
	for _, c := range b.Components {
		for _, block := range c.Blocks {
			if x == block.X && y == block.Y && block.Active {
				return true
			}
		}
	}
	return false
}

func (b *Board) CalculateBottomMap() []int {
	bottom := make([]int, b.Width)
	for i := 0; i < len(bottom); i++ {
		bottom[i] = b.Height
	}

	for _, c := range b.Components {
		if !c.IsPlaced {
			continue
		}

		for _, block := range c.Blocks {
			if block.Y < bottom[block.X] {
				bottom[block.X] = block.Y + 1
			}
		}

	}
	for i := 0; i < len(bottom); i++ {
		bottom[i] = b.Height - bottom[i]
	}
	return bottom
}

func (b *Board) CalculateSideMap() []int {
	side := make([]int, b.Height)
	for _, c := range b.Components {
		if !c.IsPlaced {
			continue
		}

		for _, block := range c.Blocks {
			if block.Active {
				side[block.Y] += 1
			}
		}
	}
	return side
}

func (b *Board) CalculateCompletedLines() int {
	min := b.Height
	bottom := b.CalculateBottomMap()
	for _, v := range bottom {
		if v < min {
			min = v
		}
	}

	return min
}

func (b *Board) RemoveCompletedLines() int {
	sides := b.CalculateSideMap()
	completed := 0
	for i := 0; i < len(sides); i++ {
		if sides[i] == b.Width {
			b.RemoveRow(i)
			completed += 1
		}
	}
	return completed
}

func (b *Board) RemoveRow(row int) {
	for _, c := range b.Components {
		if !c.IsPlaced {
			continue
		}

		for _, block := range c.Blocks {
			if block.Y == row {
				block.Active = false
			}

			if block.Y < row {
				block.Y += 1
			}
		}
	}
}

func (b *Board) RemoveInactiveComponents() {
	out := []*Component{}
	for _, c := range b.Components {
		allActive := true
		for _, block := range c.Blocks {
			if !block.Active {
				allActive = false
			}
		}
		if allActive {
			b.Components = append(b.Components, c)
		}
	}
	b.Components = out
}

func (b *Board) GameOver(engine *DisplayEngine) {
	engine.DrawLine(7, 10, "Game Over")
	engine.DrawLine(9, 20, "idiot")
	termbox.Sync()
	termbox.PollEvent()
	os.Exit(1)
}

func (b *Board) AddComponent(c *Component) {
	b.Components = append(b.Components, c)
	b.CurrentComponent = c
}

func (b *Board) IsOutBounds(c *Component) bool {

	if c.X < 0 || c.X+c.Width > b.Width {
		return true
	}

	if c.Y < 0 || c.Y+c.Height >= b.Height {
		b.IsPlaced = true
		return true
	}
	return false
}

func (b *Board) IsCollide(component *Component) bool {
	for _, c := range b.Components {
		if c == component {
			continue
		}

		for _, b1 := range component.Blocks {
			if !b1.Active {
				continue
			}
			for _, b2 := range c.Blocks {
				if !b2.Active {
					continue
				}
				if b1.X == b2.X && b1.Y == b2.Y && b1.Active && b2.Active {
					b.IsPlaced = true
					return true
				}
			}
		}
		// }
	}
	return false
}

func (b *Board) DropComponent() {
	b.IsPlaced = false
	for !b.IsCollide(b.CurrentComponent) && !b.IsOutBounds(b.CurrentComponent) {
		b.CurrentComponent.SavePrevious()
		b.CurrentComponent.Move(0, 1)
	}
	fmt.Println("drop component", b.IsCollide(b.CurrentComponent), b.IsOutBounds(b.CurrentComponent))
	b.CurrentComponent.Revert()

	b.CurrentComponent.IsPlaced = true
	b.AddComponent(RandomComponent())
}
