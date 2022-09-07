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
	White       color.RGBA = color.RGBA{255, 255, 255, 255}
	Grey        color.RGBA = color.RGBA{127, 127, 127, 255}
	Black       color.RGBA = color.RGBA{0, 0, 0, 255}
	Yellow      color.RGBA = color.RGBA{242, 202, 92, 255}
	LightYellow color.RGBA = color.RGBA{251, 245, 222, 255}
	DarkBlue    color.RGBA = color.RGBA{0, 68, 116, 255}
)

var (
	PieceImages map[game.Player]map[game.PieceType]*sdl.Texture
)

const (
	assetsFolder string = "assets"
)

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

func RenderState(renderer *sdl.Renderer, state *game.State, rect *sdl.Rect) {
	cellW, cellH := float32(rect.W)/8.0, float32(rect.H)/8.0
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			var sqCol color.RGBA
			if (i+j)%2 == 0 {
				sqCol = LightYellow
			} else {
				sqCol = Yellow
			}
			sqRect := &sdl.FRect{X: float32(rect.X) + cellW*float32(j), Y: float32(rect.Y) + cellH*float32(i), W: cellW, H: cellH}
			RectF(renderer, sqRect, sqCol)
			if state.Board[i][j] != nil {
				renderer.CopyF(PieceImages[state.Board[i][j].Owner][state.Board[i][j].Type], nil, sqRect)
			}
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Chess", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
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
	PieceImages, err = LoadPieceImages(renderer)
	if err != nil {
		panic(err)
	}

	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	state := game.NewStartState(game.White)

	running := true
	for running {
		Clear(renderer, White)
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
		RenderState(renderer, state, boardRect)
		renderer.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
	}
}
