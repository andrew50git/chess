package engine

import (
	"chess/game"
	"chess/util"
	"fmt"
	"time"
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
	depth := 5
	moves := state.GetEngineMoves(player)
	var best *game.Move
	var moveI int
	start := time.Now()
	for time.Since(start) < time.Second {
		best, moveI, _ = getBestMove(state, moves, player, depth, BigNum)
		moves = util.RemoveIndex(moves, moveI)
		moves = append([]game.Move{*best}, moves...) //TODO: multiple move priority
		depth++
		fmt.Println(depth, best)
	}
	ch <- best
}

func getBestMove(state *game.State, moves []game.Move, player game.Player, depth int, max float32) (*game.Move, int, float32) {
	state.Turn = player
	bestI := -1
	bestEval := -BigNum
	for i, m := range moves {
		captureType := game.NilPiece
		convertType := game.NilPiece
		if m.Capture != nil {
			captureType = state.Board[m.Capture.X][m.Capture.Y].Type
		}
		if m.ConvertType != captureType {
			convertType = m.ConvertType
		}

		if captureType == game.King {
			return &moves[i], i, BigNum
		}
		state.RunMove(m)
		var ev float32
		if depth == 1 {
			ev = evalState(state, player)
		} else {
			_, _, ev = getBestMove(state, state.GetEngineMoves((player+1)%2), (player+1)%2, depth-1, -bestEval)
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
		return nil, -1, -BigNum
	}
	return &moves[bestI], bestI, bestEval
}

func evalState(state *game.State, pov game.Player) float32 {
	var res float32 = 0
	//res += 0.1 * float32(len(state.GetMoves(pov)))
	//res -= 0.1 * float32(len(state.GetMoves((pov+1)%2)))
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
