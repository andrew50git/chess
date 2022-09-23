package engine

import (
	"chess/game"
	"chess/util"
	"fmt"
	"sort"
	"time"
)

var (
	pieceTypeToValue map[game.PieceType]float32 = map[game.PieceType]float32{
		game.Pawn:   1,
		game.Knight: 3.2,
		game.Bishop: 3.3,
		game.Rook:   5,
		game.Queen:  9,
		game.King:   1000,
	}
)

var (
	BigNum float32 = 10000000
)

type TranspositionState struct {
}

var Transpositions map[uint64]*TranspositionState

func Init() {
	InitZobrist()
}

func GetBestMove(state *game.State, player game.Player, ch chan *game.Move) {
	depth := 1
	moves := GetEngineMoves(state, player)
	var best *game.Move
	var moveI int
	start := time.Now()
	for time.Since(start) < time.Second*5 {
		best, moveI, _ = getBestMove(state, moves, player, depth, -BigNum, BigNum)
		moves = util.RemoveIndex(moves, moveI)
		moves = append([]game.Move{*best}, moves...) //TODO: multiple move priority
		fmt.Println(depth, best)
		depth++
	}
	ch <- best
}

func getBestMove(state *game.State, moves []game.Move, player game.Player, depth int, min, max float32) (*game.Move, int, float32) {
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
			return &moves[i], i, BigNum - 1
		}
		state.RunMove(m)
		var ev float32
		if depth == 1 {
			ev = evalState(state, player)
		} else {
			_, _, ev = getBestMove(state, GetEngineMoves(state, (player+1)%2), (player+1)%2, depth-1, -max, -min)
			ev = -ev
			state.Turn = player
		}
		if ev > bestEval {
			bestEval = ev
			bestI = i
		}
		min = util.Max(min, bestEval)
		state.ReverseMove(m, captureType, convertType)

		if min >= max {
			break
		}
	}
	if bestI == -1 {
		return nil, -1, -BigNum
	}
	return &moves[bestI], bestI, bestEval
}

var (
	PawnMap [][]float32 = [][]float32{{0, 0, 0, 0, 0, 0, 0, 0}, {0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}, {0.1, 0.1, 0.2, 0.3, 0.3, 0.2, 0.1, 0.1}, {0.05, 0.05, 0.1, 0.25, 0.25, 0.1, 0.05, 0.05}, {0, 0, 0, 0.2, 0.2, 0, 0, 0}, {0.05, -0.05, -0.1, 0, 0, -0.1, -0.05, 0.05}, {0.05, 0.1, 0.1, -0.2, -0.2, 0.1, 0.1, 0.05}, {0, 0, 0, 0, 0, 0, 0, 0}}
)

// TODO: piece maps
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
				if state.Board[i][j].Type == game.Knight {
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
				if state.Board[i][j].Type == game.Pawn {
					if state.Board[i][j].Owner == pov {
						if state.Board[i][j].Owner == state.Starter {
							res += PawnMap[i][j]
						} else {
							res += PawnMap[7-i][j]
						}
					} else {
						if state.Board[i][j].Owner == state.Starter {
							res -= PawnMap[i][j]
						} else {
							res -= PawnMap[7-i][j]
						}
					}
				}
			}
		}
	}
	return res
}

func evalMove(state *game.State, move game.Move) float32 {
	var res float32 = 0
	if move.Capture != nil {
		res += pieceTypeToValue[state.Board[move.Capture.X][move.Capture.Y].Type]
		res -= pieceTypeToValue[state.Board[move.Start.X][move.Start.Y].Type] * 0.1
	}
	if move.IsConversion {
		res += 10 + pieceTypeToValue[move.ConvertType]
	}
	return res
}

// TODO: move uistate.winner to state
func GetEngineMoves(state *game.State, player game.Player) []game.Move {
	moves := state.GetMoves(player)
	moveEvals := []float32{}
	for _, m := range moves {
		moveEvals = append(moveEvals, evalMove(state, m))
	}
	sort.Slice(moves, func(i, j int) bool {
		return moveEvals[i] > moveEvals[j]
	})
	return moves
}
