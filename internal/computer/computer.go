package computer

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/laranc/HackComputer/internal/keyboard"
	"github.com/laranc/HackComputer/internal/video"
)

const (
	fbLow   = 0x4000
	fbHigh  = 0x5FFF
	kbd     = 0x6000
	ramSize = 24577
	romSize = 32768
)

type Computer struct {
	ram [ramSize]int16
	rom [romSize]uint16
	cpu struct {
		d, a int16
		pc   uint16
	}
}

func NewComputer() *Computer {
	return &Computer{
		ram: [ramSize]int16{},
		rom: [romSize]uint16{},
		cpu: struct {
			d, a int16
			pc   uint16
		}{
			d:  0,
			a:  0,
			pc: 0,
		},
	}
}

func (comp *Computer) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, romSize*2)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	for i := 0; i < n/2; i++ {
		comp.rom[i] = binary.LittleEndian.Uint16(buf[i*2 : (i+1)*2])
	}
	return nil
}

func (comp *Computer) Execute(reset bool) {
	ins := comp.rom[comp.cpu.pc]
	a := ins & 0b100_0_000000_000_000
	c := (ins & 0b000_1_111111_000_000) >> 6
	d := (ins & 0b000_0_000000_111_000) >> 3
	j := ins & 0b000_0_000000_000_111
	if a == 0 {
		comp.cpu.a = int16(ins)
		comp.cpu.pc++
		return
	}
	var res int16
	var jmp uint16
	switch c {
	case 0b0101010:
		res = 0
	case 0b0111111:
		res = 1
	case 0b0111010:
		res = -1
	case 0b0001100:
		res = comp.cpu.d
	case 0b0110000:
		res = comp.cpu.a
	case 0b1110000:
		res = comp.ram[comp.cpu.a]
	case 0b0001101:
		res = ^comp.cpu.d
	case 0b0110001:
		res = ^comp.cpu.a
	case 0b1110001:
		res = ^comp.ram[comp.cpu.a]
	case 0b0001111:
		res = ^comp.cpu.d + 1
	case 0b0110011:
		res = -comp.cpu.a + 1
	case 0b1110011:
		res = -comp.ram[comp.cpu.a] + 1
	case 0b0011111:
		res = comp.cpu.d + 1
	case 0b0110111:
		res = comp.cpu.a + 1
	case 0b1110111:
		res = comp.ram[comp.cpu.a] + 1
	case 0b0001110:
		res = comp.cpu.d - 1
	case 0b0110010:
		res = comp.cpu.a - 1
	case 0b1110010:
		res = comp.ram[comp.cpu.a] - 1
	case 0b0000010:
		res = comp.cpu.d + comp.cpu.a
	case 0b1000010:
		res = comp.cpu.d + comp.ram[comp.cpu.a]
	case 0b0010011:
		res = comp.cpu.d - comp.cpu.a
	case 0b1010011:
		res = comp.cpu.d - comp.ram[comp.cpu.a]
	case 0b0000111:
		res = comp.cpu.a - comp.cpu.d
	case 0b1000111:
		res = comp.ram[comp.cpu.a] - comp.cpu.d
	case 0b0000000:
		res = comp.cpu.d & comp.cpu.a
	case 0b1000000:
		res = comp.cpu.d & comp.ram[comp.cpu.a]
	case 0b0010101:
		res = comp.cpu.d | comp.cpu.a
	case 0b1010101:
		res = comp.cpu.d | comp.ram[comp.cpu.a]
	}
	switch d {
	case 0b001:
		comp.ram[comp.cpu.a] = res
	case 0b010:
		comp.cpu.d = res
	case 0b011:
		comp.ram[comp.cpu.a] = res
		comp.cpu.d = res
	case 0b100:
		comp.cpu.a = res
	case 0b101:
		comp.cpu.a = res
		comp.ram[comp.cpu.a] = res
	case 0b110:
		comp.cpu.a = res
		comp.cpu.d = res
	case 0b111:
		comp.cpu.a = res
		comp.ram[comp.cpu.a] = res
		comp.cpu.d = res
	}
	switch j {
	case 0b001:
		if res > 0 {
			jmp = 1
		}
	case 0b010:
		if res == 0 {
			jmp = 1
		}
	case 0b011:
		if res >= 0 {
			jmp = 1
		}
	case 0b100:
		if res < 0 {
			jmp = 1
		}
	case 0b101:
		if res != 0 {
			jmp = 1
		}
	case 0b110:
		if res <= 0 {
			jmp = 1
		}
	case 0b111:
		jmp = 1
	}
	if reset {
		comp.cpu.pc = 0
	} else if jmp != 0 {
		comp.cpu.pc = uint16(comp.cpu.a)
	} else {
		comp.cpu.pc++
	}
}

func (c *Computer) ReadKeyboard() {
	c.ram[kbd] = keyboard.Poll()
}

func (c *Computer) DrawScreen() {
	video.Draw(c.ram[fbLow : fbHigh+1])
}
