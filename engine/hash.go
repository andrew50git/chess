package engine

import (
	"chess/game"
	"math/rand"
)

var (
	pieceTable   map[game.Player]map[game.PieceType][][]uint64 //player, piece type
	castleTable  map[game.Player][]uint64                      //0: Long 1: Short
	passantTable [][]uint64
	blackToMove  uint64
)

func initZobrist() {
	passantTable = make([][]uint64, 8)
	for i := 0; i <= 7; i++ {
		passantTable[i] = make([]uint64, 8)
	}
	pieceTable = make(map[game.Player]map[game.PieceType][][]uint64)
	for _, p := range game.Players {
		pieceTable[p] = make(map[game.PieceType][][]uint64)
	}
	castleTable = make(map[game.Player][]uint64)
	for _, p := range game.Players {
		for _, t := range game.PieceTypes {
			pieceTable[p][t] = make([][]uint64, 8)
			for i := 0; i <= 7; i++ {
				pieceTable[p][t][i] = make([]uint64, 8)
			}
		}
	}
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			passantTable[i][j] = rand.Uint64()
			for _, p := range game.Players {
				for _, t := range game.PieceTypes {
					pieceTable[p][t][i][j] = rand.Uint64()
				}
			}
		}
	}
	for _, p := range game.Players {
		castleTable[p] = make([]uint64, 2)
		castleTable[p][0] = rand.Uint64()
		castleTable[p][1] = rand.Uint64()
	}
	blackToMove = rand.Uint64()
}

func Hash(state *game.State) uint64 {
	var ans uint64 = 0
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				ans ^= pieceTable[state.Board[i][j].Owner][state.Board[i][j].Type][i][j]
			}
		}
	}
	if state.CanCastleLong[game.White] {
		ans ^= castleTable[game.White][0]
	}
	if state.CanCastleShort[game.White] {
		ans ^= castleTable[game.White][1]
	}
	if state.CanCastleLong[game.Black] {
		ans ^= castleTable[game.Black][0]
	}
	if state.CanCastleShort[game.Black] {
		ans ^= castleTable[game.Black][1]
	}
	if state.Turn == game.Black {
		ans ^= blackToMove
	}
	if state.PassantPos != nil {
		ans ^= passantTable[state.PassantPos.X][state.PassantPos.Y]
	}
	return ans
}

func RunMoveForHash(state *game.State, m *game.Move, hash uint64) uint64 {
	//update hash
	hash ^= pieceTable[state.Board[m.Start.X][m.Start.Y].Owner][state.Board[m.Start.X][m.Start.Y].Type][m.Start.X][m.Start.Y]
	hash ^= pieceTable[state.Board[m.Start.X][m.Start.Y].Owner][state.Board[m.Start.X][m.Start.Y].Type][m.End.X][m.End.Y]
	if m.Capture != nil {
		hash ^= pieceTable[state.Board[m.Capture.X][m.Capture.Y].Owner][state.Board[m.Capture.X][m.Capture.Y].Type][m.Capture.X][m.Capture.Y]
	}
	//update hash passant
	if state.PassantPos != nil {
		hash ^= passantTable[state.PassantPos.X][state.PassantPos.Y]
	}
	if m.IsPassant {
		hash ^= passantTable[m.End.X][m.End.Y]
	}
	//update hash castling
	hash ^= castleTable[game.White][0]
	hash ^= castleTable[game.Black][0]
	hash ^= castleTable[game.White][0]
	hash ^= castleTable[game.Black][1]
	state.RunMove(*m)
	hash ^= castleTable[game.White][0]
	hash ^= castleTable[game.Black][0]
	hash ^= castleTable[game.White][0]
	hash ^= castleTable[game.Black][1]
	return hash
}
