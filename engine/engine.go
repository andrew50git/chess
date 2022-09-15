package engine

import (
	"chess/game"
	"math/rand"
)

func GetBestMove(state *game.State, moves []game.Move) int {
	move := rand.Intn(len(moves) - 1)
	return move
}
