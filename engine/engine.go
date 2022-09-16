package engine

import (
	"chess/game"
)

var (
	pieceTypeToValue map[game.PieceType]int = map[game.PieceType]int{
		game.Pawn:   1,
		game.Knight: 3,
		game.Bishop: 3,
		game.Rook:   5,
		game.Queen:  9,
		game.King:   1000,
	}
)

func GetBestMove(state *game.State, player game.Player) *game.Move {
	return getBestMove(state, player)
}

func getBestMove(state *game.State, player game.Player) *game.Move {
	moves := state.GetEngineMoves(player)
	bestI := -1
	bestEval := -9999
	for i, m := range moves {
		captureType := game.NilPiece
		if m.Captures != nil {
			captureType = state.Board[m.Captures.X][m.Captures.Y].Type
		}
		state.RunMove(m)
		ev := evalState(state, player)
		if evalState(state, player) > bestEval {
			bestEval = ev
			bestI = i
		}
		state.ReverseMove(m, captureType)
	}
	return &moves[bestI]
}

func evalState(state *game.State, pov game.Player) int {
	res := 0
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				if state.Board[i][j].Owner == pov {
					res += pieceTypeToValue[state.Board[i][j].Type]
				} else {
					res -= pieceTypeToValue[state.Board[i][j].Type]
				}
			}
		}
	}
	return res
}
