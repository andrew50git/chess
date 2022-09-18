package game

import "fmt"

type Pos struct {
	X int
	Y int
}

func (p1 Pos) Add(p2 Pos) Pos {
	return Pos{p1.X + p2.X, p1.Y + p2.Y}
}

type Player int

var NilPlayer Player = -1

const (
	White Player = iota
	Black
	Both
)

type PieceType int

var NilPiece PieceType = -1

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
		Both:  "Both",
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
	state.CanCastleLong = map[Player]bool{
		White: true,
		Black: true,
	}
	state.CanCastleShort = map[Player]bool{
		White: true,
		Black: true,
	}
	state.IsGameEnd = false
	state.Winner = NilPlayer
	return state
}

type State struct {
	Board          [][]*Piece
	PassantPos     *Pos
	Turn           Player
	CanCastleLong  map[Player]bool
	CanCastleShort map[Player]bool
	Starter        Player
	IsGameEnd      bool
	Winner         Player
}

func (state *State) Add(pos Pos, piece Piece) {
	state.Board[pos.X][pos.Y] = &piece
}

func (state *State) Remove(pos Pos) {
	state.Board[pos.X][pos.Y] = nil
}

func (state *State) RunMove(move Move) bool {
	piece := state.Board[move.Start.X][move.Start.Y]
	if piece.Type == King {
		state.CanCastleLong[piece.Owner] = false
		state.CanCastleShort[piece.Owner] = false
	}
	if piece.Type == Rook {
		if move.Start.Y == 0 {
			state.CanCastleLong[piece.Owner] = false
		} else if move.Start.Y == 7 {
			state.CanCastleShort[piece.Owner] = false
		}
	}
	isGameEnd := false
	if move.Capture != nil {
		if state.Board[move.Capture.X][move.Capture.Y].Type == King {
			isGameEnd = true
		}
		state.Board[move.Capture.X][move.Capture.Y] = nil
	}

	state.PassantPos = nil
	if move.IsPassant {
		state.PassantPos = &Pos{move.End.X, move.End.Y}
	}
	state.Board[move.End.X][move.End.Y] = state.Board[move.Start.X][move.Start.Y]
	state.Board[move.Start.X][move.Start.Y] = nil
	if move.IsConversion && move.ConvertType != NilPiece {
		state.Board[move.End.X][move.End.Y].Type = move.ConvertType
	}
	return isGameEnd
}

func (state *State) ReverseMove(move Move, captureType PieceType, convertType PieceType) {
	state.Board[move.Start.X][move.Start.Y] = state.Board[move.End.X][move.End.Y]
	state.Board[move.End.X][move.End.Y] = nil
	if move.IsPassant {
		state.Board[move.Start.X][move.Start.Y] = &Piece{Type: Pawn, Owner: state.Turn}
	}
	if move.Capture != nil {
		state.Board[move.Capture.X][move.Capture.Y] = &Piece{Type: captureType, Owner: (state.Turn + 1) % 2}
	}
	if move.IsConversion {
		state.Board[move.Start.X][move.Start.Y].Type = Pawn
	}
	//TODO:castling, cancastle...
}

type Move struct {
	Start        Pos
	End          Pos
	Capture      *Pos // can be nil
	IsConversion bool // for pawns only
	ConvertType  PieceType
	IsCastle     bool // for kings only, moves the rook to middle of start and end
	IsPassant    bool
}

func MakeBasicMove(start Pos, end Pos, capture *Pos) Move {
	return Move{start, end, capture, false, NilPiece, false, false}
}

func OnBoard(p Pos) bool {
	return p.X >= 0 && p.Y >= 0 && p.X <= 7 && p.Y <= 7
}

func GenMovesByDirs(state *State, start Pos, dirs []Pos, otherPlayer Player) []Move {
	i, j := start.X, start.Y
	moves := []Move{}
	for _, dir := range dirs {
		currPos := Pos{i, j}.Add(dir)
		mvs := 0
		for OnBoard(currPos) && (state.Board[currPos.X][currPos.Y] == nil || state.Board[currPos.X][currPos.Y].Owner == otherPlayer) {
			var capture *Pos = nil
			if state.Board[currPos.X][currPos.Y] != nil && state.Board[currPos.X][currPos.Y].Owner == otherPlayer {
				capture = &Pos{currPos.X, currPos.Y}
			}
			moves = append(moves, MakeBasicMove(Pos{i, j}, currPos, capture))
			mvs++
			if state.Board[currPos.X][currPos.Y] != nil && state.Board[currPos.X][currPos.Y].Owner == otherPlayer {
				break
			}
			currPos = currPos.Add(dir)
		}
	}
	return moves
}

func (state *State) GetAttacks(player Player) [][]bool { //TODO: GET ATTACKS: FOR CASTLING
	attacks := make([][]bool, 8)
	for i := 0; i <= 7; i++ {
		attacks[i] = make([]bool, 8)
	}
	return attacks
}

func (state *State) GetMoves(player Player) []Move { //TODO: CASTLING
	oppPlayer := (player + 1) % 2
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
						moves = append(moves, Move{Pos{i, j}, Pos{i + 2*dir, j}, nil, false, NilPiece, false, true})
					}
					canMoveForward := (i+dir >= 0) && (i+dir <= 7)
					if canMoveForward && state.Board[i+dir][j] == nil {
						convert := (i+dir) == 0 || (i+dir) == 7
						moves = append(moves, Move{Pos{i, j}, Pos{i + dir, j}, nil, convert, NilPiece, false, false})
					}
					if canMoveForward && j-1 >= 0 && (state.Board[i+dir][j-1] != nil && state.Board[i+dir][j-1].Owner == oppPlayer) {
						convert := (i+dir) == 0 || (i+dir) == 7
						moves = append(moves, Move{Pos{i, j}, Pos{i + dir, j - 1}, &Pos{i + dir, j - 1}, convert, NilPiece, false, false})
					}
					if j-1 >= 0 && state.PassantPos != nil && state.PassantPos.X == i && state.PassantPos.Y == j-1 {
						moves = append(moves, Move{Pos{i, j}, Pos{i + dir, j - 1}, &Pos{i, j - 1}, false, NilPiece, false, true})
					}
					if canMoveForward && j+1 <= 7 && (state.Board[i+dir][j+1] != nil && state.Board[i+dir][j+1].Owner == oppPlayer) {
						convert := (i+dir) == 0 || (i+dir) == 7
						moves = append(moves, Move{Pos{i, j}, Pos{i + dir, j + 1}, &Pos{i + dir, j + 1}, convert, NilPiece, false, false})
					}
					if j+1 >= 0 && state.PassantPos != nil && state.PassantPos.X == i && state.PassantPos.Y == j+1 {
						moves = append(moves, Move{Pos{i, j}, Pos{i + dir, j + 1}, &Pos{i, j + 1}, false, NilPiece, false, true})
					}
				case Knight:
					knightDirs := []Pos{{1, 2}, {-1, 2}, {1, -2}, {-1, -2}, {2, 1}, {-2, 1}, {2, -1}, {-2, -1}}
					for _, dir := range knightDirs {
						newPos := Pos{i, j}.Add(dir)
						if OnBoard(newPos) && (state.Board[newPos.X][newPos.Y] == nil || state.Board[newPos.X][newPos.Y].Owner == oppPlayer) {
							var capture *Pos = nil
							if state.Board[newPos.X][newPos.Y] != nil && state.Board[newPos.X][newPos.Y].Owner == oppPlayer {
								capture = &Pos{newPos.X, newPos.Y}
							}
							moves = append(moves, MakeBasicMove(Pos{i, j}, newPos, capture))
						}
					}
				case Rook:
					rookDirs := []Pos{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
					moves = append(moves, GenMovesByDirs(state, Pos{i, j}, rookDirs, oppPlayer)...)
				case Bishop:
					bishopDirs := []Pos{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}}
					moves = append(moves, GenMovesByDirs(state, Pos{i, j}, bishopDirs, oppPlayer)...)
				case Queen:
					queenDirs := []Pos{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}}
					moves = append(moves, GenMovesByDirs(state, Pos{i, j}, queenDirs, oppPlayer)...)
				case King:
					kingDirs := []Pos{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}}
					for _, dir := range kingDirs {
						newPos := Pos{i, j}.Add(dir)
						if OnBoard(newPos) && (state.Board[newPos.X][newPos.Y] == nil || state.Board[newPos.X][newPos.Y].Owner == oppPlayer) {
							var capture *Pos = nil
							if state.Board[newPos.X][newPos.Y] != nil && state.Board[newPos.X][newPos.Y].Owner == oppPlayer {
								capture = &Pos{newPos.X, newPos.Y}
							}
							moves = append(moves, MakeBasicMove(Pos{i, j}, newPos, capture))
						}
					}
					/*
						if state.CanCastleLong[player] || state.CanCastleShort[player] {
							oppAttacks := state.GetAttacks(oppPlayer)
							var playerRank int
							if player == state.Starter {
								playerRank = 7
							} else {
								playerRank = 0
							}
							if state.CanCastleShort[player] {
								if state.Board[playerRank][5] == nil && state.Board[playerRank][6] == nil &&
									!oppAttacks[playerRank][4] && !oppAttacks[playerRank][5] && !oppAttacks[playerRank][6] {
									//can castle short

								}
							}
							if state.CanCastleLong[player] {
								if state.Board[playerRank][1] == nil && state.Board[playerRank][2] == nil && state.Board[playerRank][3] == nil &&
									!oppAttacks[playerRank][2] && !oppAttacks[playerRank][3] && !oppAttacks[playerRank][4] {
									//can castle long

								}
							}
						}
					*/
				}
			}
		}
	}
	return moves
}

func (state *State) GetEngineMoves(player Player) []Move {
	moves := state.GetMoves(player)
	resMoves := []Move{}
	for _, m := range moves {
		if m.IsConversion {
			m.ConvertType = Queen
			resMoves = append([]Move{m}, resMoves...)
			m.ConvertType = Knight
			resMoves = append([]Move{m}, resMoves...)
		} else {
			if m.Capture != nil {
				resMoves = append([]Move{m}, resMoves...)
			} else {
				resMoves = append(resMoves, m)
			}
		}
	}
	return resMoves
}
