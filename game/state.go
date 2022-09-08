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

func PrintBoard(board [][]*Piece) {
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if board[i][j] != nil {
				fmt.Printf("%v ", board[i][j].Type)
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Println()
	}
}

func NewStartState(starter Player) *State {
	nonStarter := (starter + 1) % 2
	state := &State{Turn: starter}
	state.Board = make([][]*Piece, 8)
	for i := 0; i <= 7; i++ {
		state.Board[i] = make([]*Piece, 8)
	}
	for i := 0; i <= 1; i++ {
		for j := 0; j <= 7; j++ {
			state.Add(Pos{i + 6, j}, Piece{StartPieces[i][j], starter})
		}
	}
	for i := 1; i >= 0; i-- {
		for j := 0; j <= 7; j++ {
			state.Add(Pos{1 - i, j}, Piece{StartPieces[i][j], nonStarter})
		}
	}
	return state
}

type State struct {
	Board      [][]*Piece
	Turn       Player
	HasCastled map[Player]bool
	Starter    Player
}

func (state *State) Add(pos Pos, piece Piece) {
	state.Board[pos.X][pos.Y] = &piece
}

func (state *State) Remove(pos Pos) {
	state.Board[pos.X][pos.Y] = nil
}

func (state *State) RunMove(move Move) {
	if move.Captures != nil {
		state.Board[move.Captures.X][move.Captures.Y] = nil
	}
	state.Board[move.End.X][move.End.Y] = state.Board[move.Start.X][move.Start.Y]
	state.Board[move.Start.X][move.Start.Y] = nil
	state.Turn = (state.Turn + 1) % 2
}

type Move struct {
	Start    Pos
	End      Pos
	Captures *Pos // can be nil
	Convert  bool // for pawns only
	IsCastle bool // for kings only
}

func (state *State) GetMoves(player Player) []Move {
	//otherPlayer := (player + 1) % 2
	moves := []Move{}
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if state.Board[i][j] != nil && state.Board[i][j].Owner == player {
				switch state.Board[i][j].Type {
				case Pawn:
					isUnmoved := (player == state.Starter && i == 6) || (player != state.Starter && i == 1)
					var dir int
					if player == state.Starter {
						dir = -1
					} else {
						dir = 1
					}
					if isUnmoved && state.Board[i+dir*2][j] == nil && state.Board[i+dir][j] == nil {
						moves = append(moves, Move{Pos{i, j}, Pos{i + 2*dir, j}, nil, false, false})
					}
					if state.Board[i+dir][j] == nil {
						convert := (i+dir) == 0 || (i+dir) == 7
						moves = append(moves, Move{Pos{i, j}, Pos{i + dir, j}, nil, convert, false})
					}
				}
			}
		}
	}
	return moves
}
