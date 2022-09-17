package engine

import (
	"chess/game"
	"chess/util"
	"fmt"
)

var (
	pieceTypeToValue map[game.PieceType]float32 = map[game.PieceType]float32{
		game.Pawn:   1,
		game.Knight: 3,
		game.Bishop: 3,
		game.Rook:   5,
		game.Queen:  9,
		game.King:   1000,
	}
)

var (
	BigNum float32 = 100000
)

func GetBestMove(state *game.State, player game.Player, ch chan *game.Move) {
	piecesLeft := 0
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				piecesLeft++
			}
		}
	}
	depth := 5
	if piecesLeft <= 7 {
		depth += (7 - piecesLeft) / 2
	}
	best, _ := getBestMove(state, player, depth, BigNum)
	ch <- best
}

func getBestMove(state *game.State, player game.Player, depth int, max float32) (*game.Move, float32) {
	state.Turn = player
	moves := state.GetEngineMoves(player)
	bestI := -1
	bestEval := -BigNum
	if len(moves) == 0 {
		fmt.Println(depth)
	}
	for i, m := range moves {
		captureType := game.NilPiece
		convertType := game.NilPiece
		if m.Capture != nil {
			captureType = state.Board[m.Capture.X][m.Capture.Y].Type
		}
		if m.ConvertType != captureType {
			convertType = m.ConvertType
		}
		state.RunMove(m)
		var ev float32
		if depth == 1 {
			ev = evalState(state, player)
		} else {
			_, ev = getBestMove(state, (player+1)%2, depth-1, -bestEval)
			ev = -ev
			state.Turn = player
		}
		if ev > bestEval {
			bestEval = ev
			bestI = i
		}
		state.ReverseMove(m, captureType, convertType)

		if ev >= max {
			break
		}
	}
	if bestI == -1 {
		return nil, -BigNum
	}
	return &moves[bestI], bestEval
}

func evalState(state *game.State, pov game.Player) float32 {
	var res float32 = 0
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				if state.Board[i][j].Owner == pov {
					res += pieceTypeToValue[state.Board[i][j].Type]
				} else {
					res -= pieceTypeToValue[state.Board[i][j].Type]
				}
				if state.Board[i][j].Type == game.Knight || state.Board[i][j].Type == game.Pawn {
					rowValue := pieceTypeToValue[state.Board[i][j].Type] * ((3.5 - util.Abs(float32(i)-3.5)) / 7)
					colValue := pieceTypeToValue[state.Board[i][j].Type] * ((3.5 - util.Abs(float32(j)-3.5)) / 7)
					if state.Board[i][j].Owner == pov {
						res += rowValue
						res += colValue
					} else {
						res -= rowValue
						res -= colValue
					}
				}
			}
		}
	}
	return res
}

//TODO: move uistate.winner to state
