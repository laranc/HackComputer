package video

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 512
	height = 256
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer
)

func init() {
	var err error
	if window, err = sdl.CreateWindow("Video Out", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN); err != nil {
		panic(err)
	}
	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		panic(err)
	}
}

func Draw(fb []int16) {
	if len(fb) != width*height/16 {
		panic(fmt.Errorf("Frame buffer error: %d", len(fb)))
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()
	renderer.SetDrawColor(0, 0, 0, 255)
	for i, word := range fb {
		row := int32((i * 16) / 512)
		colStart := (i * 16) % 512
		for b := 0; b < 16; b++ {
			col := int32(colStart + b)
			if (word & (1 << (15 - b))) != 0 {
				renderer.DrawPoint(col, row)
			}
		}
	}
	renderer.Present()
}
