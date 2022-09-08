package main

import (
	"chess/game"
	"chess/util"
	"fmt"
	"image/color"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
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
)

const (
	assetsFolder string = "assets"
)

type UIState struct {
	gameState *game.State
	selected  *game.Pos
}

func LoadPieceImages(renderer *sdl.Renderer) (map[game.Player]map[game.PieceType]*sdl.Texture, error) {
	pieceImages := map[game.Player]map[game.PieceType]*sdl.Texture{
		game.White: {},
		game.Black: {},
	}

	for _, v := range game.PieceTypes {
		for _, p := range game.Players {
			var err error
			pieceImages[p][v], err = img.LoadTexture(renderer, fmt.Sprintf("%v/%v/%v.png", assetsFolder, game.PlayerToString[p], game.PieceTypeToString[v]))
			if err != nil {
				return nil, err
			}
		}
	}
	return pieceImages, nil
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

func RenderState(renderer *sdl.Renderer, uiState *UIState, rect *sdl.Rect) {
	state := uiState.gameState
	cellW, cellH := float32(rect.W)/8.0, float32(rect.H)/8.0
	moves := state.GetMoves(state.Turn)
	movingPoints := []game.Pos{}

	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			if uiState.selected != nil && uiState.selected.X == j && uiState.selected.Y == i {
				for _, m := range moves {
					if m.Start.X == i && m.Start.Y == j { //NOT INVERTED!!!
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
			if util.Contains(movingPoints, game.Pos{X: i, Y: j}) {
				RectF(renderer, &sdl.FRect{X: float32(rect.X) + cellW*float32(j) + 0.375*float32(cellW), Y: float32(rect.Y) + cellH*float32(i) + 0.375*float32(cellH), W: cellW * 0.25, H: cellH * 0.25}, darkRed)
			}
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

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
	pieceImages, err = LoadPieceImages(renderer)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	state := game.NewStartState(game.White)
	uiState := &UIState{state, nil}
	running := true
	for running {
		Clear(renderer, white)
		w, h := window.GetSize()
		boardRect := &sdl.Rect{}
		if w < h {
			boardRect.X = 0
			boardRect.Y = (h - w) / 2
		} else {
			boardRect.X = (w - h) / 2
			boardRect.Y = 0
		}
		boardRect.W = int32(util.Min(int(w), int(h)))
		boardRect.H = int32(util.Min(int(w), int(h)))
		RenderState(renderer, uiState, boardRect)
		renderer.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_LEFT && e.Type == sdl.MOUSEBUTTONDOWN {
					mX, mY, _ := sdl.GetMouseState()
					if mX >= boardRect.X && mY >= boardRect.Y && mX <= boardRect.X+boardRect.W && mY <= boardRect.Y+boardRect.H {
						relX, relY := float32(mX-boardRect.X), float32(mY-boardRect.Y)
						cellW, cellH := float32(boardRect.W)/8.0, float32(boardRect.H)/8.0
						sqX, sqY := util.Min(int(relX/cellH), 7), util.Min(int(relY/cellW), 7)

						//selecting own piece
						if state.Board[sqY][sqX] != nil && state.Board[sqY][sqX].Owner == state.Turn { //INVERTED X AND Y!!!
							if uiState.selected != nil && uiState.selected.X == sqX && uiState.selected.Y == sqY {
								uiState.selected = nil
							} else {
								uiState.selected = &game.Pos{X: sqX, Y: sqY}
							}
						} else if uiState.selected != nil { //selecting place to move
							moves := state.GetMoves(state.Turn)
							for _, m := range moves {
								if m.Start.X == uiState.selected.Y && m.Start.Y == uiState.selected.X && m.End.X == sqY && m.End.Y == sqX {
									state.RunMove(m)
									uiState.selected = nil
									break
								}
							}
						}
					}
				}
			}
		}
	}
}
