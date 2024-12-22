package keyboard

import "github.com/veandco/go-sdl2/sdl"

func init() {
	sdl.StartTextInput()
}

func Poll() int16 {
	if e := sdl.PollEvent(); e != nil {
		switch t := e.(type) {
		case *sdl.TextInputEvent:
			text := t.GetText()
			if len(text) > 0 {
				char := []rune(text)[0]
				return int16(char)
			}
		case *sdl.KeyboardEvent:
			return int16(t.Keysym.Sym)
		}
	}
	return 0
}
