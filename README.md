# HackComputer
Nand2Tetris Hack Computer Emulator written in Go

## Script Example
```lua
return {
  file = "path/to/file", -- Path to the binary
  ticks = 100, -- Number of CPU cycles to execute, leave as zero to run infinitly
  delay = 16, -- Default SDL delay, set to a higher number to slow execution
  output = { 0, 1, 25 }, -- The memory addresses to print the values from after program execution
}
```
