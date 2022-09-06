package main

import (
	"chess/game"
	"chess/util"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	White color.RGBA = color.RGBA{255, 255, 255, 255}
	Grey  color.RGBA = color.RGBA{127, 127, 127, 255}
	Black color.RGBA = color.RGBA{0, 0, 0, 255}
)

func Clear(renderer *sdl.Renderer, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.Clear()
}

func Rect(renderer *sdl.Renderer, rect *sdl.Rect, color color.RGBA) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillRect(rect)
}

func RenderState(renderer *sdl.Renderer, state *game.State, rect *sdl.Rect) {
	cellW, cellH := int(rect.W/8), int(rect.H/8)
	for i := 0; i <= 7; i++ {
		for j := 0; j <= 7; j++ {
			var sqCol color.RGBA
			if (i+j)%2 == 0 {
				sqCol = Grey
			} else {
				sqCol = Black
			}
			Rect(renderer, &sdl.Rect{X: int32(int(rect.X) + cellW*i), Y: int32(int(rect.Y) + cellH*j), W: int32(cellW), H: int32(cellH)}, sqCol)
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
