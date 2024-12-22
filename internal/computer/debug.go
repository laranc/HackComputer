package computer

func (c *Computer) GetRegs() (d int16, a int16, pc uint16, kbd int16) {
	return c.cpu.d, c.cpu.a, c.cpu.pc, c.ram[kbd]
}

func (c *Computer) GetValues(addrs []int) []int16 {
	out := make([]int16, 0, len(addrs))
	for _, v := range addrs {
		out = append(out, c.ram[v])
	}
	return out
}

func (c *Computer) GetIns() (ins uint16) {
	return c.rom[c.cpu.pc]
}

func (c *Computer) DumpRAM() []int16 {
	ram := make([]int16, ramSize)
	copy(ram, c.ram[:])
	return ram
}

func (c *Computer) DumpROM() []uint16 {
	rom := make([]uint16, romSize)
	copy(rom, c.rom[:])
	return rom
}
