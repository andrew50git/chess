package game

import "fmt"

type Pos struct {
	X int
	Y int
}

type Player int

const (
	White Player = iota
	Black
)

type PieceType int

const (
	King PieceType = iota
	Rook
	Bishop
	Queen
	Knight
	Pawn
)

var (
	PieceTypes []PieceType = []PieceType{
		King, Rook, Bishop, Queen, Knight, Pawn,
	}
	Players []Player = []Player{
		White, Black,
	}
	PieceTypeToString map[PieceType]string = map[PieceType]string{
		King:   "King",
		Rook:   "Rook",
		Bishop: "Bishop",
		Queen:  "Queen",
		Knight: "Knight",
		Pawn:   "Pawn",
	}
	PlayerToString map[Player]string = map[Player]string{
		White: "White",
		Black: "Black",
	}
)

type Piece struct {
	Type  PieceType
	Owner Player
}

var (
	StartPieces [][]PieceType = [][]PieceType{{Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn}, {Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}}
)

func NewStartState(starter Player) *State {
	nonStarter := (starter + 1) % 2
	state := &State{Turn: starter}
	state.Board = make([][]*Piece, 8)
	for i := 0; i <= 7; i++ {
		state.Board[i] = make([]*Piece, 8)
	}
	for i := 0; i <= 1; i++ {
		for j := 0; j <= 7; j++ {
			state.AddPiece(Pos{i + 6, j}, Piece{StartPieces[i][j], starter})
		}
	}
	for i := 1; i >= 0; i-- {
		for j := 0; j <= 7; j++ {
			state.AddPiece(Pos{1 - i, j}, Piece{StartPieces[i][j], nonStarter})
		}
	}
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil {
				fmt.Printf("%v ", state.Board[i][j].Type)
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Println()
	}
	return state
}

type State struct {
	Board [][]*Piece
	Turn  Player
}

func (state *State) AddPiece(pos Pos, piece Piece) {
	state.Board[pos.X][pos.Y] = &piece
}

func (state *State) RemovePiece(pos Pos) {
	state.Board[pos.X][pos.Y] = nil
}

func (state *State) MovePiece(start Pos, end Pos) {
	state.Board[end.X][end.Y] = state.Board[start.X][start.Y]
	state.Board[start.X][start.Y] = nil
}

type Move struct {
	Start    Pos
	End      Pos
	Captures []Pos
}
