package main

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	BoardHorizontal = "-"
	BoardVertical   = "|"
	BoardCorner     = "+"
)

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

	for i := 1; i < height; i++ {
		engine.DrawCell(x, y+i, BoardVertical)
		engine.DrawCell(x+width*BoardHorizontalScale+1, y+i, BoardVertical)
	}

	engine.DrawLine(x, y+height, top)
}

func (engine *DisplayEngine) EraseBoard(x, y, width, height int) {
	for i := 1; i < height; i++ {
		engine.DrawLine(x+1, y+i, strings.Repeat(" ", width*BoardHorizontalScale))
	}
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
		if block.Active {
			engine.DrawLine(block.X*BoardHorizontalScale+1, block.Y+1, strings.Repeat(" ", BoardHorizontalScale))
		}
	}
	engine.ResetColors()
}

func (engine *DisplayEngine) DrawScore(x, y, score int) {
	xPos := x * BoardHorizontalScale
	yPos := y * BoardVerticalScale

	engine.DrawLine(xPos, yPos, "Score:")
	engine.DrawLine(xPos, yPos+1, fmt.Sprintf("%d", score))
}

func (engine *DisplayEngine) DrawRules(x, y int) {
	xPos := x * BoardHorizontalScale
	yPos := y * BoardVerticalScale
	engine.DrawLine(xPos, yPos, "Rules:")
	engine.DrawLine(xPos, yPos+1, "a     : move left   <--  ")
	engine.DrawLine(xPos, yPos+2, "d     : move right  -->  ")
	engine.DrawLine(xPos, yPos+3, "s     : move down    v   ")
	engine.DrawLine(xPos, yPos+4, "space : drop         V   ")
	engine.DrawLine(xPos, yPos+5, "r     : rotate       O   ")
}

func (engine *DisplayEngine) ResetColors() {
	engine.F = termbox.ColorDefault
	engine.B = termbox.ColorDefault

}

func (engine *DisplayEngine) DrawComponentStats(x, y int, c *Component) {
	xPos := x * BoardHorizontalScale
	yPos := y * BoardVerticalScale

	posLine := fmt.Sprintf("(%d, %d) -> (%d, %d)   ", c.X, c.Y, c.X+c.Width, c.Y+c.Height)
	engine.DrawLine(xPos, yPos, c.Name)
	engine.DrawLine(xPos, yPos+1, posLine)
}

func (engine *DisplayEngine) DrawBottomMap(x, y int, bottom []int) {
	out := ""
	for _, v := range bottom {
		out += fmt.Sprintf("%2d", v)
	}
	xPos := x*BoardHorizontalScale + 1
	yPos := y * BoardVerticalScale

	engine.DrawLine(xPos, yPos, out)
}

func (engine *DisplayEngine) DrawSideMap(x, y int, side []int) {
	xPos := x * BoardHorizontalScale
	yPos := y * BoardVerticalScale

	for i := 0; i < len(side); i++ {
		engine.DrawLine(xPos, yPos+i, fmt.Sprintf("%2d", side[i]))
	}
}

func (engine *DisplayEngine) EraseComponent(c *Component) {
	for _, block := range c.Blocks {
		engine.DrawLine(block.X*BoardHorizontalScale+1, block.Y+1, strings.Repeat(" ", BoardHorizontalScale))
	}
	engine.F = termbox.ColorDefault
}
