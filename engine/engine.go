package engine

import (
	"chess/deepcopy"
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
	bigNum float32 = 10000000
)

type transpositionState struct {
}

var transpositions map[uint64]*transpositionState
var transpositionEvals map[uint64]float32 //pov of white

func Init() {
	initZobrist()
	transpositionEvals = make(map[uint64]float32)
}

func isEndGame(state *game.State) bool {
	numQueens := 0
	numMinors := 0
	numNonQueenPieces := 0
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				if state.Board[i][j].Type == game.Queen {
					numQueens++
				} else if state.Board[i][j].Type != game.King && state.Board[i][j].Type != game.Pawn {
					numNonQueenPieces++
				}
				if state.Board[i][j].Type == game.Bishop || state.Board[i][j].Type == game.Knight {
					numMinors++
				}
			}
		}
	}
	return numQueens == 0 || (numMinors <= 2 && numNonQueenPieces <= 1)
}

func GetBestMove(state *game.State, player game.Player, ch chan *game.Move) {
	isEndGame := isEndGame(state)
	fmt.Printf("isEndGame: %v\n", isEndGame)
	if isEndGame {
		pieceMaps[game.King] = kingMapEndGame
	}
	depth := 2
	moves := getEngineMoves(state, player)
	var best *game.Move
	var moveI int
	start := time.Now()
	for time.Since(start) < time.Second*5 {
		bestOppEval := -bigNum
		for i := 0; i < len(moves); i++ { //TODO: maybe use concurrency? or multithreaded
			copiedStateIface, _ := deepcopy.Anything(state)
			copiedState := copiedStateIface.(*game.State)
			copiedState.RunMove(moves[i])
			oppEval := -getBestEval(copiedState, copiedState.GetMoves((player+1)%2), (player+1)%2, depth-1, -bigNum, -bestOppEval, Hash(copiedState))
			if oppEval > bestOppEval {
				best = &moves[i]
				moveI = i
				bestOppEval = oppEval
			}
		}
		moves = util.RemoveIndex(moves, moveI)
		moves = append([]game.Move{*best}, moves...) //TODO: multiple move priority
		fmt.Printf("depth: %v, best: %v\n", depth, best)
		depth++
	}
	ch <- best
	state.RunMove(*best)
	fmt.Printf("eval for %v: %v\n", game.PlayerToString[player], evalState(state, player, Hash(state)))
}

func getBestEval(state *game.State, moves []game.Move, player game.Player, depth int, min, max float32, currHash uint64) float32 {
	state.Turn = player
	bestI := -1
	bestEval := -bigNum
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
			return bigNum - 1
		}
		oldHash := currHash

		currHash := RunMoveForHash(state, &m, currHash) //runs original RunMove
		var ev float32
		if depth == 1 {
			ev = evalState(state, player, currHash)
		} else {
			ev = getBestEval(state, getEngineMoves(state, (player+1)%2), (player+1)%2, depth-1, -max, -min, currHash)
			ev = -ev
			state.Turn = player
		}
		if ev > bestEval {
			bestEval = ev
			bestI = i
		}
		min = util.Max(min, bestEval)
		state.ReverseMove(m, captureType, convertType)
		currHash = oldHash

		if min >= max {
			break
		}
	}
	if bestI == -1 {
		return -bigNum
	}
	return bestEval
}

var (
	pawnMap [][]float32 = [][]float32{{0, 0, 0, 0, 0, 0, 0, 0},
		{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
		{0.1, 0.1, 0.2, 0.3, 0.3, 0.2, 0.1, 0.1},
		{0.05, 0.05, 0.1, 0.25, 0.25, 0.1, 0.05, 0.05},
		{0, 0, 0, 0.2, 0.2, 0, 0, 0},
		{0.05, -0.05, -0.1, 0, 0, -0.1, -0.05, 0.05},
		{0.05, 0.1, 0.1, -0.2, -0.2, 0.1, 0.1, 0.05},
		{0, 0, 0, 0, 0, 0, 0, 0}}
	knightMap [][]float32 = [][]float32{{-0.5, -0.4, -0.3, -0.3, -0.3, -0.3, -0.4, -0.5},
		{-0.4, -0.2, 0, 0, 0, 0, -0.2, -0.4},
		{-0.3, 0, 0.1, 0.15, 0.15, 0.1, 0, -0.3},
		{-0.3, 0.05, 0.15, 0.2, 0.2, 0.15, 0.05, -0.3},
		{-0.3, 0, 0.15, 0.2, 0.2, 0.15, 0, -0.3},
		{-0.3, 0.05, 0.1, 0.15, 0.15, 0.1, 0.05, -0.3},
		{-0.4, -0.2, 0, 0.05, 0.05, 0, -0.2, -0.4},
		{-0.5, -0.4, -0.3, -0.3, -0.3, -0.3, -0.4, -0.5}}
	bishopMap [][]float32 = [][]float32{{-0.2, -0.1, -0.1, -0.1, -0.1, -0.1, -0.1, -0.2},
		{-0.1, 0, 0, 0, 0, 0, 0, -0.1},
		{-0.1, 0, 0.05, 0.1, 0.1, 0.05, 0, -0.1},
		{-0.1, 0.05, 0.05, 0.1, 0.1, 0.05, 0.05, -0.1},
		{-0.1, 0, 0.1, 0.1, 0.1, 0.1, 0, -0.1},
		{-0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, -0.1},
		{-0.1, 0.05, 0, 0, 0, 0, 0.05, -0.1},
		{-0.2, -0.1, -0.1, -0.1, -0.1, -0.1, -0.1, -0.2}}
	rookMap [][]float32 = [][]float32{{0, 0, 0, 0, 0, 0, 0, 0},
		{0.05, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.05},
		{-0.05, 0, 0, 0, 0, 0, 0, -0.05},
		{-0.05, 0, 0, 0, 0, 0, 0, -0.05},
		{-0.05, 0, 0, 0, 0, 0, 0, -0.05},
		{-0.05, 0, 0, 0, 0, 0, 0, -0.05},
		{-0.05, 0, 0, 0, 0, 0, 0, -0.05},
		{0, 0, 0, 0.05, 0.05, 0, 0, 0}}
	queenMap [][]float32 = [][]float32{{-0.2, -0.1, -0.1, -0.05, -0.05, -0.1, -0.1, -0.2},
		{-0.1, 0, 0, 0, 0, 0, 0, -0.1},
		{-0.1, 0, 0.05, 0.05, 0.05, 0.05, 0, -0.1},
		{-0.05, 0, 0.05, 0.05, 0.05, 0.05, 0, -0.05},
		{0, 0, 0.05, 0.05, 0.05, 0.05, 0, -0.05},
		{-0.1, 0.05, 0.05, 0.05, 0.05, 0.05, 0, -0.1},
		{-0.1, 0, 0.05, 0, 0, 0, 0, -0.1},
		{-0.2, -0.1, -0.1, -0.05, -0.05, -0.1, -0.1, -0.2}}
	kingMapMiddleGame [][]float32 = [][]float32{{-0.3, -0.4, -0.4, -0.5, -0.5, -0.4, -0.4, -0.3},
		{-0.3, -0.4, -0.4, -0.5, -0.5, -0.4, -0.4, -0.3},
		{-0.3, -0.4, -0.4, -0.5, -0.5, -0.4, -0.4, -0.3},
		{-0.3, -0.4, -0.4, -0.5, -0.5, -0.4, -0.4, -0.3},
		{-0.2, -0.3, -0.3, -0.4, -0.4, -0.3, -0.3, -0.2},
		{-0.1, -0.2, -0.2, -0.2, -0.2, -0.2, -0.2, -0.1},
		{0.2, 0.2, 0, 0, 0, 0, 0.2, 0.2},
		{0.2, 0.3, 0.1, 0, 0, 0.1, 0.3, 0.2}}
	kingMapEndGame [][]float32 = [][]float32{{-0.5, -0.4, -0.3, -0.2, -0.2, -0.3, -0.4, -0.5},
		{-0.3, -0.2, -0.1, 0, 0, -0.1, -0.2, -0.3},
		{-0.3, -0.1, 0.2, 0.3, 0.3, 0.2, -0.1, -0.3},
		{-0.3, -0.1, 0.3, 0.4, 0.4, 0.3, -0.1, -0.3},
		{-0.3, -0.1, 0.3, 0.4, 0.4, 0.3, -0.1, -0.3},
		{-0.3, -0.1, 0.2, 0.3, 0.3, 0.2, -0.1, -0.3},
		{-0.3, -0.3, 0, 0, 0, 0, -0.3, -0.3},
		{-0.5, -0.3, -0.3, -0.3, -0.3, -0.3, -0.3, -0.5}}
	pieceMaps map[game.PieceType][][]float32 = map[game.PieceType][][]float32{game.Pawn: pawnMap, game.Knight: knightMap, game.Bishop: bishopMap, game.Rook: rookMap, game.Queen: queenMap, game.King: kingMapMiddleGame}
)

func evalState(state *game.State, pov game.Player, currHash uint64) float32 {
	stateHash := currHash
	if _, ok := transpositionEvals[stateHash]; ok {
		if pov == game.Black {
			return -transpositionEvals[stateHash]
		} else {
			return transpositionEvals[stateHash]
		}
	}
	var res float32 = 0

	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				if state.Board[i][j].Owner == pov {
					res += pieceTypeToValue[state.Board[i][j].Type]
				} else {
					res -= pieceTypeToValue[state.Board[i][j].Type]
				}
				if state.Board[i][j].Owner == pov {
					if state.Board[i][j].Owner == state.Starter {
						res += pieceMaps[state.Board[i][j].Type][i][j]
					} else {
						res += pieceMaps[state.Board[i][j].Type][7-i][j]
					}
				} else {
					if state.Board[i][j].Owner == state.Starter {
						res -= pieceMaps[state.Board[i][j].Type][i][j]
					} else {
						res -= pieceMaps[state.Board[i][j].Type][7-i][j]
					}
				}
			}
		}
	}
	//TODO: isolated pawns
	for j := 0; j <= 7; j++ {
		currPovPawns := 0
		currNonPovPawns := 0
		for i := 0; i <= 7; i++ {
			if state.Board[i][j] != nil && state.Board[i][j].Type == game.Pawn {
				if state.Board[i][j].Owner == pov {
					if state.Board[i][j].Owner == state.Starter { //blocked pawns
						if state.Board[i-1][j] != nil {
							res -= 0.5
						}
					} else {
						if state.Board[i+1][j] != nil {
							res -= 0.5
						}
					}
					currPovPawns++
				} else {
					if state.Board[i][j].Owner == state.Starter {
						if state.Board[i-1][j] != nil {
							res += 0.5
						}
					} else {
						if state.Board[i+1][j] != nil {
							res += 0.5
						}
					}
					currNonPovPawns++
				}
			}
		}
		if currPovPawns >= 2 { //doubled
			res -= 0.5
		}
		if currNonPovPawns >= 2 {
			res += 0.5
		}
	}
	if pov == game.Black {
		transpositionEvals[stateHash] = -res
	} else {
		transpositionEvals[stateHash] = res
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
func getEngineMoves(state *game.State, player game.Player) []game.Move {
	moves := state.GetMoves(player)
	moveEvals := []float32{}
	for i, m := range moves {
		if m.IsConversion {
			moves[i].ConvertType = game.Queen
		}
		moveEvals = append(moveEvals, evalMove(state, m))
	}
	sort.Slice(moves, func(i, j int) bool {
		return moveEvals[i] > moveEvals[j]
	})
	return moves
}
