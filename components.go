package main

import (
	"fmt"
	"math/rand"

	"github.com/nsf/termbox-go"
)

type Component struct {
	Blocks        []*Block
	X, Y          int
	Height, Width int
	Color         Color
	Name          string
	IsPlaced      bool

	Previous *Component
}

type Block struct {
	X, Y   int
	Active bool
}

func (b *Block) Normalize(x, y int) (int, int) {
	return b.X - x, b.Y - y
}

func RandomComponent() *Component {
	choices := 7
	choice := rand.Intn(choices)

	switch choice {
	case 0:
		return NewJComponent()
	case 1:
		return NewLComponent()
	case 2:
		return NewOComponent()
	case 3:
		return NewTComponent()
	case 4:
		return NewZComponent()
	case 5:
		return NewSComponent()
	case 6:
		return NewIComponent()
	}
	return nil
}

func NewJComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorGreen}
	c.Blocks = append(c.Blocks, &Block{1, 0, true}) // 1
	c.Blocks = append(c.Blocks, &Block{1, 1, true}) // 2
	c.Blocks = append(c.Blocks, &Block{1, 2, true}) //43
	c.Blocks = append(c.Blocks, &Block{0, 2, true})
	c.Height = 3
	c.Width = 2
	c.Name = "J-Block"
	c.Move(4, 0)
	return c
}

func NewLComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorBlue}
	c.Blocks = append(c.Blocks, &Block{0, 0, true}) // 1
	c.Blocks = append(c.Blocks, &Block{0, 1, true}) // 2
	c.Blocks = append(c.Blocks, &Block{0, 2, true}) // 34
	c.Blocks = append(c.Blocks, &Block{1, 2, true})
	c.Height = 3
	c.Width = 2
	c.Name = "L-Block"
	c.Move(4, 0)

	return c
}

func NewOComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorYellow}
	c.Blocks = append(c.Blocks, &Block{0, 0, true}) //12
	c.Blocks = append(c.Blocks, &Block{1, 0, true}) //34
	c.Blocks = append(c.Blocks, &Block{0, 1, true})
	c.Blocks = append(c.Blocks, &Block{1, 1, true})
	c.Height = 2
	c.Width = 2
	c.Name = "O-Block"
	c.Move(4, 0)

	return c
}

func NewTComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorRed}
	c.Blocks = append(c.Blocks, &Block{1, 0, true}) // 1
	c.Blocks = append(c.Blocks, &Block{0, 1, true}) //234
	c.Blocks = append(c.Blocks, &Block{1, 1, true})
	c.Blocks = append(c.Blocks, &Block{2, 1, true})
	c.Height = 2
	c.Width = 3
	c.Name = "T-Block"
	c.Move(4, 0)

	return c
}

func NewZComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorCyan}
	c.Blocks = append(c.Blocks, &Block{0, 0, true}) //12
	c.Blocks = append(c.Blocks, &Block{1, 0, true}) // 34
	c.Blocks = append(c.Blocks, &Block{1, 1, true})
	c.Blocks = append(c.Blocks, &Block{2, 1, true})
	c.Height = 2
	c.Width = 3
	c.Name = "Z-Block"
	c.Move(4, 0)

	return c
}

func NewSComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorMagenta}
	c.Blocks = append(c.Blocks, &Block{1, 0, true}) // 12
	c.Blocks = append(c.Blocks, &Block{2, 0, true}) //34
	c.Blocks = append(c.Blocks, &Block{0, 1, true})
	c.Blocks = append(c.Blocks, &Block{1, 1, true})
	c.Height = 2
	c.Width = 3
	c.Name = "S-Block"
	c.Move(4, 0)

	return c
}

func NewIComponent() *Component {
	c := new(Component)
	c.X, c.Y = 0, 0
	c.Color = Color{termbox.ColorWhite}
	c.Blocks = append(c.Blocks, &Block{0, 0, true}) //1
	c.Blocks = append(c.Blocks, &Block{0, 1, true}) //2
	c.Blocks = append(c.Blocks, &Block{0, 2, true}) //3
	c.Blocks = append(c.Blocks, &Block{0, 3, true}) //4
	c.Height = 4
	c.Width = 1
	c.Name = "I-Block"
	c.Move(4, 0)

	return c
}

func (c *Component) HasPrevious() bool {
	return c.Previous != nil
}

func (c *Component) Revert() bool {
	if c.Previous == nil {
		return false
	}

	fmt.Println("SWAG MAX")
	fmt.Printf("%+v\n", c.Previous)

	c.Blocks = c.Previous.Blocks
	fmt.Println("SWAG1")
	c.Width, c.Height = c.Previous.Width, c.Previous.Height
	fmt.Println("SWAG2")

	c.X, c.Y = c.Previous.X, c.Previous.Y
	fmt.Println("SWAG3")

	c.IsPlaced = c.Previous.IsPlaced
	fmt.Println("SWAG4")

	c.Previous = c.Previous.Previous
	fmt.Println("SWAG4")

	fmt.Println("SWAG5")

	return true
}

func (c *Component) SavePrevious() {
	previous := *c
	previous.Blocks = make([]*Block, len(c.Blocks))
	for i := range c.Blocks {
		tmp := *c.Blocks[i]
		previous.Blocks[i] = &tmp
	}
	c.Previous = &previous
}

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

// recalculate component x/y/width/height
func (c *Component) Normalize() {
	minX, minY, maxX, maxY := 100, 100, 0, 0
	for _, block := range c.Blocks {
		if block.X < minX {
			minX = block.X
		}
		if block.Y < minY {
			minY = block.Y
		}

		if block.X > maxX {
			maxX = block.X
		}
		if block.Y > maxY {
			maxY = block.Y
		}
	}
	originalX, originalY := c.X, c.Y
	c.X, c.Y = minX, minY
	c.Width = maxX - minX + 1
	c.Height = maxY - minY + 1
	c.Move(originalX-c.X, originalY-c.Y)
}

func GetRotationMatrix(size int) [][]RotateMove {
	if size == 3 || size == 2 {
		return Rotate3Matrix
	} else if size == 4 {
		return Rotate4Matrix
	} else {
		panic("no rotation matrix found")
	}
	return nil
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

	rotationMatrix := GetRotationMatrix(size)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if grid[x][y] == nil {
				continue
			}

			move := rotationMatrix[y][x]
			grid[x][y].X += move.X
			grid[x][y].Y += move.Y
		}
	}

	c.Normalize()
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

var Rotate4Matrix = [][]RotateMove{
	[]RotateMove{RotateMove{3, 0}, {2, 1}, {1, 2}, {0, 3}},
	[]RotateMove{RotateMove{2, -1}, {1, 0}, {0, 1}, {-1, 2}},
	[]RotateMove{RotateMove{1, -2}, {0, -1}, {-1, 0}, {-2, 1}},
	[]RotateMove{RotateMove{0, -3}, {-1, -2}, {-2, -1}, {-3, 0}},
}
