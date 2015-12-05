package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

func DrawScreen(b *Board, engine *DisplayEngine) {
	for _, c := range b.Components {
		engine.DrawComponent(c)
	}

	engine.DrawComponentStats(13, 1, b.CurrentComponent)
	engine.DrawBottomMap(0, 26, b.CalculateBottomMap())
	engine.DrawSideMap(11, 1, b.CalculateSideMap())
	engine.DrawScore(13, 23, b.CompletedLines*10)
	termbox.Sync()

}

func DropHandler(b *Board, engine *DisplayEngine) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for _ = range ticker.C {
		b.HandlerMutex.Lock()

		fmt.Println("Game over checks", b.IsPlaced, b.CurrentComponent.HasPrevious())
		if b.IsGameOver || (b.IsPlaced && !b.CurrentComponent.HasPrevious()) {
			b.IsGameOver = true
			b.GameOver(engine)
			os.Exit(2)
		}

		b.CurrentComponent.SavePrevious()

		engine.EraseComponent(b.CurrentComponent)
		b.CurrentComponent.Move(0, 1)

		wasPlaced := b.IsPlaced

		// Is move valid? or game over? if not revert
		if b.IsOutBounds(b.CurrentComponent) || b.IsCollide(b.CurrentComponent) {
			fmt.Println("Current component", b.CurrentComponent)
			if !b.CurrentComponent.Revert() {
				b.IsGameOver = true
				b.GameOver(engine)
				os.Exit(3)
			}

			if wasPlaced {
				b.IsPlaced = false
				b.DropComponent()
			}
		}

		DrawScreen(b, engine)

		b.HandlerMutex.Unlock()

	}
}

func EventHandler(b *Board, engine *DisplayEngine) {

	for {
		event := termbox.PollEvent()
		// fmt.Printf("%+v\n", event)

		b.HandlerMutex.Lock()

		b.CurrentComponent.SavePrevious()
		engine.EraseComponent(b.CurrentComponent)

		if event.Ch == 'a' {
			b.CurrentComponent.Move(-1, 0)
		} else if event.Ch == 'd' {
			b.CurrentComponent.Move(1, 0)
			// } else if event.Ch == 'w' {
			// b.CurrentComponent.Move(0, -1)
		} else if event.Ch == 's' {
			b.CurrentComponent.Move(0, 1)
		} else if event.Ch == 'r' {
			b.CurrentComponent.Rotate()
		} else if event.Ch == 'q' {
			os.Exit(1)
		} else if event.Key == 32 {
			b.DropComponent()
		}

		b.IsPlaced = false

		// Is move valid? or game over? if not revert
		if b.IsOutBounds(b.CurrentComponent) || b.IsCollide(b.CurrentComponent) {
			if !b.CurrentComponent.Revert() {
				b.IsGameOver = true
				b.GameOver(engine)
			}
		}

		// Remvoe completed lines
		completedLines := b.RemoveCompletedLines()
		if completedLines > 0 {
			engine.EraseBoard(0, 0, BoardWidthUnits, BoardHeightUnits)
			b.CompletedLines += completedLines
		}

		DrawScreen(b, engine)

		b.HandlerMutex.Unlock()
	}
}

func main() {
	engine := NewDisplayEngine()
	defer engine.Close()

	c := RandomComponent()

	board := NewBoard(BoardWidthUnits, BoardHeightUnits)
	board.AddComponent(c)

	engine.DrawBoard(0, 0, BoardWidthUnits, BoardHeightUnits)
	engine.DrawRules(13, 10)
	DrawScreen(board, engine)

	go DropHandler(board, engine)
	EventHandler(board, engine)
}
