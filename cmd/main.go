package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"

	"github.com/laranc/HackComputer/internal/computer"
	"github.com/laranc/HackComputer/internal/script"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/yuin/gopher-lua"
)

func init() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
}

func main() {
	defer sdl.Quit()
	if len(os.Args) <= 1 {
		panic("No input script provided\n")
	}
	config := os.Args[1]
	tbl, err := script.Load(config)
	if err != nil {
		panic(err)
	}
	path := string(script.Get[lua.LString](tbl, "file"))
	ticks := int(script.Get[lua.LNumber](tbl, "ticks"))
	if ticks == 0 {
		ticks = math.MaxInt
	}
	delay := uint32(script.Get[lua.LNumber](tbl, "delay"))
	outSlice := script.SliceFromTable[lua.LNumber](script.Get[*lua.LTable](tbl, "output"))
	outAddrs := make([]int, len(outSlice))
	for i, v := range outSlice {
		outAddrs[i] = int(v)
	}
	comp := computer.NewComputer()
	if err = comp.Load(path); err != nil {
		panic(err)
	}
	for range ticks {
		comp.ReadKeyboard()
		comp.Execute(false)
		comp.DrawScreen()
		flush()
		d, a, pc, kbd := comp.GetRegs()
		ins := comp.GetIns()
		fmt.Printf("D: %04x, A: %04x, PC: %04x\nROM[%04x]: %016b\nKBD: %04x\n", d, a, pc, pc, ins, kbd)
		sdl.Delay(delay)
	}
	outVals := comp.GetValues(outAddrs)
	fmt.Println("_____________________________________")
	for i, v := range outVals {
		fmt.Printf("M[%d]\t|   %016b\n", outAddrs[i], v)
	}
	fmt.Println("_____________________________________")
	for {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
	}
}

func flush() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
