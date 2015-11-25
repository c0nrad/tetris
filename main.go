package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

type Player struct {
} // GetMove

type Color struct {
	Color termbox.Attribute
}

type DisplayEngine struct {
	F termbox.Attribute
	B termbox.Attribute
}

func NewDisplayEngine() *DisplayEngine {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.HideCursor()
	return &DisplayEngine{F: termbox.ColorDefault, B: termbox.ColorDefault}
}

func (engine *DisplayEngine) Close() {
	termbox.Close()
}

func (engine *DisplayEngine) DrawBoard(x, y, width, height int) {
	top := BoardCorner + strings.Repeat(BoardHorizontal, width*BoardHorizontalScale) + BoardCorner

	engine.DrawLine(x, y, top)

	for i := 0; i < height; i++ {
		engine.DrawCell(x, y+i, BoardVertical)
		engine.DrawCell(x+width*BoardHorizontalScale+1, y+i, BoardVertical)
	}

	engine.DrawLine(x, y+height, top)
}

func (engine *DisplayEngine) DrawCell(x, y int, c string) {
	termbox.SetCell(x, y, rune(c[0]), engine.F, engine.B)
}

func (engine *DisplayEngine) DrawLine(x, y int, line string) {
	for i, c := range line {
		termbox.SetCell(x+i, y, rune(c), engine.F, engine.B)
	}
}

func (engine *DisplayEngine) DrawComponent(c *Component) {
	engine.F = c.Color.Color
	engine.B = c.Color.Color
	for _, block := range c.Blocks {
		engine.DrawLine(block.X*BoardHorizontalScale+1, block.Y+1, strings.Repeat(" ", BoardHorizontalScale))
	}
	engine.F = termbox.ColorDefault
	engine.B = termbox.ColorDefault
}

func (engine *DisplayEngine) EraseComponent(c *Component) {
	for _, block := range c.Blocks {
		engine.DrawLine(block.X*BoardHorizontalScale+1, block.Y+1, strings.Repeat(" ", BoardHorizontalScale))
	}
	engine.F = termbox.ColorDefault
}

type Board struct {
	Width  int
	Height int

	Components       []*Component
	CurrentComponent *Component
}

func NewBoard(width, height int) *Board {
	return &Board{width, height, nil, nil}
}

func (b *Board) AddComponent(c *Component) {
	b.Components = append(b.Components, c)
	b.CurrentComponent = c
}

func (b *Board) BoundsCheck(c *Component) {
	for _, block := range c.Blocks {
		if block.X < 0 {
			c.Move(1, 0)
		} else if block.X+c.Width >= b.Width/BoardHorizontalScale {
			c.Move(-1, 0)
		}

		if block.Y < 0 {
			c.Move(0, 1)
		} else if block.Y+c.Height > b.Height+1 {
			c.Move(0, -1)
		}
	}
}

type Component struct {
	Blocks        []*Block
	X, Y          int
	Height, Width int
	Color         Color
} // rotate, move

func (c *Component) Move(dx, dy int) {
	c.X += dx
	c.Y += dy

	for _, block := range c.Blocks {
		block.X += dx
		block.Y += dy
	}
}

func (c *Component) Dimensionality() int {
	if c.Height > c.Width {
		return c.Height
	}
	return c.Width
}

func (c *Component) Rotate() {
	size := c.Dimensionality()
	grid := make([][]*Block, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]*Block, size)
	}

	for _, block := range c.Blocks {
		x, y := block.Normalize(c.X, c.Y)
		grid[x][y] = block
	}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if grid[x][y] == nil {
				continue
			}

			fmt.Println(grid, grid[x][y])
			move := Rotate3Matrix[x][y]
			grid[x][y].X += move.X
			grid[x][y].Y += move.Y
			fmt.Println(move, x, y)

			termbox.Sync()
			termbox.PollEvent()
		}
	}

}

type RotateGrid [][]*Block

type RotateMove struct {
	X, Y int
}

var Rotate3Matrix = [][]RotateMove{
	[]RotateMove{RotateMove{2, 0}, {1, 1}, {0, 2}},
	[]RotateMove{RotateMove{1, -1}, {0, 0}, {-1, 1}},
	[]RotateMove{RotateMove{0, -2}, {-1, -1}, {-2, 0}},
}

func NewJComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorRed}
	c.Blocks = append(c.Blocks, &Block{1, 0, true}) // 1
	c.Blocks = append(c.Blocks, &Block{1, 1, true}) // 2
	c.Blocks = append(c.Blocks, &Block{1, 2, true}) //43
	c.Blocks = append(c.Blocks, &Block{0, 2, true})
	c.Height = 3
	c.Width = 2
	return c
}

type Block struct {
	X, Y   int
	Active bool
}

func (b *Block) Normalize(x, y int) (int, int) {
	return b.X - x, b.Y - y
}

const (
	BoardHorizontalScale = 2
	BoardWidthUnits      = 10
	BoardWidth           = BoardHorizontalScale * (BoardWidthUnits + 2)

	BoardVerticallScale = 1
	BoardHeightUnits    = 25
	BoardHeight         = 25

	BoardHorizontal = "-"
	BoardVertical   = "|"
	BoardCorner     = "+"
)

func EventHandler(b *Board, engine *DisplayEngine) {
	for {
		event := termbox.PollEvent()
		engine.EraseComponent(b.CurrentComponent)

		if event.Ch == 'a' {
			b.CurrentComponent.Move(-1, 0)
		} else if event.Ch == 'd' {
			b.CurrentComponent.Move(1, 0)
		} else if event.Ch == 'w' {
			b.CurrentComponent.Move(0, -1)
		} else if event.Ch == 's' {
			b.CurrentComponent.Move(0, 1)
		} else if event.Ch == 'r' {
			b.CurrentComponent.Rotate()
		} else if event.Ch == 'q' {
			os.Exit(1)
		}

		b.BoundsCheck(b.CurrentComponent)

		engine.EraseComponent(b.CurrentComponent)
		engine.DrawComponent(b.CurrentComponent)

		termbox.Sync()
	}

}

func main() {
	engine := NewDisplayEngine()
	defer engine.Close()

	board := NewBoard(BoardWidth, BoardHeight)

	c := NewJComponent()
	c.Move(5, 5)
	board.AddComponent(c)

	engine.DrawBoard(0, 0, BoardWidthUnits, BoardHeightUnits)
	engine.DrawComponent(c)
	termbox.Sync()

	EventHandler(board, engine)
}
