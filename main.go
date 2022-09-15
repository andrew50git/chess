package main

import (
	"chess/engine"
	"chess/game"
	"chess/util"
	"fmt"
	"image/color"
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	white       color.RGBA = color.RGBA{255, 255, 255, 255}
	grey        color.RGBA = color.RGBA{127, 127, 127, 255}
	black       color.RGBA = color.RGBA{0, 0, 0, 255}
	yellow      color.RGBA = color.RGBA{242, 202, 92, 255}
	lightYellow color.RGBA = color.RGBA{251, 245, 222, 255}
	darkBlue    color.RGBA = color.RGBA{0, 68, 116, 255}
	darkRed     color.RGBA = color.RGBA{102, 0, 0, 255}
)

var (
	pieceImages map[game.Player]map[game.PieceType]*sdl.Texture
	openSans    *ttf.Font
)

const (
	assetsFolder string = "assets"
)

type UIState struct {
	gameState     *game.State
	selected      *game.Pos //selected, convertMenu are inverted from screen coordinates
	convertMenu   *game.Pos
	winner        game.Player
	prevMoveStart *game.Pos
}

func LoadPieceImages(renderer *sdl.Renderer) error {
	pieceImages = map[game.Player]map[game.PieceType]*sdl.Texture{
		game.White: {},
		game.Black: {},
	}

	for _, v := range game.PieceTypes {
		for _, p := range game.Players {
			var err error
			pieceImages[p][v], err = img.LoadTexture(renderer, fmt.Sprintf("%v/%v/%v.png", assetsFolder, game.PlayerToString[p], game.PieceTypeToString[v]))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadFonts(renderer *sdl.Renderer) error {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	openSans, err = ttf.OpenFont(fmt.Sprintf("%v/%v/opensans.ttf", wd, assetsFolder), 60)
	return err
}

func Clear(renderer *sdl.Renderer, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.Clear()
}

func Rect(renderer *sdl.Renderer, rect *sdl.Rect, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRect(rect)
}

func RectF(renderer *sdl.Renderer, rect *sdl.FRect, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRectF(rect)
}

func TextF(renderer *sdl.Renderer, text string, posX float32, posY float32, font *ttf.Font, color color.RGBA) {
	textSurf, err := openSans.RenderUTF8Solid(text, sdl.Color{R: color.R, G: color.G, B: color.B, A: color.A})
	if err != nil {
		panic(err)
	}
	textTex, err := renderer.CreateTextureFromSurface(textSurf)
	if err != nil {
		panic(err)
	}
	textW := float32(textSurf.W)
	textH := float32(textSurf.H)
	textRect := &sdl.FRect{X: posX - textW/2, Y: posY - textH/2, W: textW, H: textH}
	renderer.CopyF(textTex, nil, textRect)
}

func RenderState(renderer *sdl.Renderer, uiState *UIState, rect *sdl.FRect) {
	state := uiState.gameState
	cellW, cellH := float32(rect.W)/8.0, float32(rect.H)/8.0
	moves := state.GetMoves(state.Turn)
	movingPoints := []game.Pos{}

	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if uiState.selected != nil && uiState.selected.X == j && uiState.selected.Y == i {
				for _, m := range moves {
					if m.Start.X == i && m.Start.Y == j {
						movingPoints = append(movingPoints, m.End)
					}
				}
			}
		}
	}

	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			var sqCol color.RGBA
			if (i+j)%2 == 0 {
				sqCol = lightYellow
			} else {
				sqCol = yellow
			}
			if uiState.selected != nil && uiState.selected.X == j && uiState.selected.Y == i {
				sqCol = darkBlue
			}
			sqRect := &sdl.FRect{X: float32(rect.X) + cellW*float32(j), Y: float32(rect.Y) + cellH*float32(i), W: cellW, H: cellH}
			RectF(renderer, sqRect, sqCol)
			if state.Board[i][j] != nil {
				renderer.CopyF(pieceImages[state.Board[i][j].Owner][state.Board[i][j].Type], nil, sqRect)
			}
			smallSq := &sdl.FRect{X: float32(rect.X) + cellW*float32(j) + 0.375*float32(cellW), Y: float32(rect.Y) + cellH*float32(i) + 0.375*float32(cellH), W: cellW * 0.25, H: cellH * 0.25}
			if util.Contains(movingPoints, game.Pos{X: i, Y: j}) {
				RectF(renderer, smallSq, darkRed)
			}
			if uiState.prevMoveStart != nil && uiState.prevMoveStart.X == i && uiState.prevMoveStart.Y == j {
				RectF(renderer, smallSq, darkBlue)
			}
		}
	}

	if uiState.convertMenu != nil {
		sqRect := &sdl.FRect{X: float32(rect.X) + cellW*float32(uiState.convertMenu.Y), Y: float32(rect.Y) + cellH*float32(uiState.convertMenu.X), W: cellW, H: cellH} //INVERTED !!!!
		RectF(renderer, sqRect, darkRed)
		sqRect.W /= 2.0
		sqRect.H /= 2.0
		renderer.CopyF(pieceImages[state.Turn][game.Bishop], nil, sqRect)
		sqRect.X += cellW / 2.0
		renderer.CopyF(pieceImages[state.Turn][game.Knight], nil, sqRect)
		sqRect.X -= cellW / 2.0
		sqRect.Y += cellH / 2.0
		renderer.CopyF(pieceImages[state.Turn][game.Queen], nil, sqRect)
		sqRect.X += cellW / 2.0
		renderer.CopyF(pieceImages[state.Turn][game.Rook], nil, sqRect)
	}

	if uiState.winner != game.NilPlayer {
		TextF(renderer, fmt.Sprintf("%v won", game.PlayerToString[uiState.winner]), rect.X+rect.W/2, rect.Y+rect.H/2, openSans, black)
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err := ttf.Init()
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Chess", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 800, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	window.SetResizable(true)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	err = LoadPieceImages(renderer)
	if err != nil {
		panic(err)
	}
	err = LoadFonts(renderer)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	humanPlayer := game.White
	state := game.NewStartState(humanPlayer)
	uiState := &UIState{state, nil, nil, game.NilPlayer, nil}
	running := true
	isGameEnd := false
	for running {
		Clear(renderer, white)
		w, h := window.GetSize()
		boardRect := &sdl.FRect{}
		if w < h {
			boardRect.X = 0
			boardRect.Y = float32(h-w) / 2
		} else {
			boardRect.X = float32(w-h) / 2
			boardRect.Y = 0
		}
		boardRect.W = float32(util.Min(int(w), int(h)))
		boardRect.H = float32(util.Min(int(w), int(h)))
		RenderState(renderer, uiState, boardRect)
		renderer.Present()
	eventLoop:
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if isGameEnd {
					break eventLoop
				}
				if e.Button == sdl.BUTTON_LEFT && e.Type == sdl.MOUSEBUTTONDOWN {
					mXInt, mYInt, _ := sdl.GetMouseState()
					mX, mY := float32(mXInt), float32(mYInt)
					if mX >= boardRect.X && mY >= boardRect.Y && mX <= boardRect.X+boardRect.W && mY <= boardRect.Y+boardRect.H {
						relX, relY := mX-boardRect.X, mY-boardRect.Y
						sqW, sqH := boardRect.W/8.0, boardRect.H/8.0
						sqC, sqR := util.Min(int(relX/sqH), 7), util.Min(int(relY/sqW), 7)                                                           //col, row C is X, R is Y
						if uiState.convertMenu != nil && uiState.convertMenu.X == sqR && uiState.convertMenu.Y == sqC && state.Turn == humanPlayer { // convert menu
							menuX, menuY := int((relX-sqW*float32(sqC))/(sqW/2.0)), int((relY-sqH*float32(sqR))/(sqH/2.0))
							if menuX < 0 || menuY < 0 || menuX >= 2 || menuY >= 2 {
								break
							}
							if menuX == 0 && menuY == 0 {
								state.Board[uiState.convertMenu.X][uiState.convertMenu.Y].Type = game.Bishop
							} else if menuX == 1 && menuY == 0 {
								state.Board[uiState.convertMenu.X][uiState.convertMenu.Y].Type = game.Knight
							} else if menuX == 0 && menuY == 1 {
								state.Board[uiState.convertMenu.X][uiState.convertMenu.Y].Type = game.Queen
							} else {
								state.Board[uiState.convertMenu.X][uiState.convertMenu.Y].Type = game.Rook
							}
							state.Turn = (state.Turn + 1) % 2
							uiState.convertMenu = nil
							break
						}

						//selecting own piece
						if state.Board[sqR][sqC] != nil && state.Board[sqR][sqC].Owner == state.Turn && state.Turn == humanPlayer {
							if uiState.selected != nil && uiState.selected.X == sqC && uiState.selected.Y == sqR {
								uiState.selected = nil
							} else {
								uiState.selected = &game.Pos{X: sqC, Y: sqR}
							}
						} else if uiState.selected != nil && state.Turn == humanPlayer { //selecting place to move
							moves := state.GetMoves(state.Turn)
							for _, m := range moves {
								if m.Start.X == uiState.selected.Y && m.Start.Y == uiState.selected.X && m.End.X == sqR && m.End.Y == sqC {
									isGameEnd = state.RunMove(m)
									uiState.prevMoveStart = &game.Pos{X: m.Start.X, Y: m.Start.Y}
									if isGameEnd {
										uiState.winner = state.Turn
									}
									if m.IsConversion && m.ConvertType == game.NilPiece {
										uiState.convertMenu = &game.Pos{X: m.End.X, Y: m.End.Y}
									} else {
										state.Turn = (state.Turn + 1) % 2
									}
									uiState.selected = nil
									break
								}
							}
						}
					}
				}
			}
		}
		//end event loop
		if state.Turn != humanPlayer { //engine move
			engineMoves := state.GetEngineMoves((humanPlayer + 1) % 2) //TODO: other player func
			m := engineMoves[engine.GetBestMove(state, engineMoves)]
			state.RunMove(m)
			uiState.prevMoveStart = &game.Pos{X: m.Start.X, Y: m.Start.Y}
			state.Turn = humanPlayer
		}
	}
}
