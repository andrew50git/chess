package bench

import (
	"chess/engine"
	"chess/game"
	"fmt"
)

func Bench() {
	state := game.NewStartState(game.White)
	game.PrintBoard(state.Board)
	gameDone := false
	turn := game.White
	var ch chan *game.Move = make(chan *game.Move, 1)
	for !gameDone {
		fmt.Printf("%v to move\n", game.PlayerToString[turn])
		engine.GetBestMove(state, turn, ch)
		m := <-ch
		if m.Capture != nil && state.Board[m.Capture.X][m.Capture.Y].Type == game.King {
			gameDone = true
		}
		state.RunMove(*m)
		game.PrintBoard(state.Board)
		turn = (turn + 1) % 2
	}
}
