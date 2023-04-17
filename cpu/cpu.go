package cpu

import (
	"os"

	"fmt"

	"math/rand"

	"chip8/util"

	"time"
)

// CPU struct holds the cpu state during emulation
type CPU struct {
	opcode     uint16
	memory     []uint8
	regV       []uint8
	regI       uint16
	pc         uint16
	Gfx        []uint8
	delayTimer uint8
	soundTimer uint8
	stack      []uint16
	sp         uint16
	Key        []uint8
	DrawFlag   bool
}

// Initialize sets up the initial state of the CPU
func (cpu *CPU) Initialize() {
	font := [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	cpu.pc = 0x200
	cpu.opcode = 0
	cpu.regI = 0
	cpu.memory = make([]uint8, 4096)
	cpu.regV = make([]uint8, 16)
	cpu.Gfx = make([]uint8, 2048)
	cpu.delayTimer = 0
	cpu.soundTimer = 0
	cpu.stack = make([]uint16, 16)
	cpu.sp = 0
	cpu.Key = make([]uint8, 16)
	//loadfontset
	for i := 0; i < 80; i++ {
		//set font
		cpu.memory[i] = font[i]
	}
	cpu.DrawFlag = true
}

// LoadGame loads the game into memory and reads it
func (cpu *CPU) LoadGame(game string) {
	f, err := os.Open(game)
	util.HandleError(err)
	fStat, err := f.Stat()
	util.HandleError(err)
	b := make([]byte, 1)
	for i := 0; int64(i) < fStat.Size(); i++ {
		_, err := f.Read(b)
		util.HandleError(err)

		cpu.memory[i+512] = uint8(b[0])
	}
	defer f.Close()

}

// lookup table
var opcodes = map[uint16]func(*CPU){
	0x0000: (*CPU).execute00E0,
	0x000E: (*CPU).execute00EE,
	0x1000: (*CPU).execute1NNN,
	0x2000: (*CPU).execute2NNN,
	0x3000: (*CPU).execute3XNN,
	0x4000: (*CPU).execute4XNN,
	0x5000: (*CPU).execute5XY0,
	0x6000: (*CPU).execute6XNN,
	0x7000: (*CPU).execute7XNN,
	0x8000: (*CPU).execute8XY0,
	0x8001: (*CPU).execute8XY1,
	0x8002: (*CPU).execute8XY2,
	0x8003: (*CPU).execute8XY3,
	0x8004: (*CPU).execute8XY4,
	0x8005: (*CPU).execute8XY5,
	0x8006: (*CPU).execute8XY6,
	0x8007: (*CPU).execute8XY7,
	0x800E: (*CPU).execute8XYE,
	0x9000: (*CPU).execute9XY0,
	0xA000: (*CPU).executeANNN,
	0xB000: (*CPU).executeBNNN,
	0xC000: (*CPU).executeCXNN,
	0xD000: (*CPU).executeDXYN,
	0xE09E: (*CPU).executeEX9E,
	0xE0A1: (*CPU).executeEXA1,
	0xF007: (*CPU).executeFX07,
	0xF00A: (*CPU).executeFX0A,
	0xF015: (*CPU).executeFX15,
	0xF018: (*CPU).executeFX18,
	0xF01E: (*CPU).executeFX1E,
	0xF029: (*CPU).executeFX29,
	0xF033: (*CPU).executeFX33,
	0xF055: (*CPU).executeFX55,
	0xF065: (*CPU).executeFX65,
}

// EmulateCycle fetch, decode and execute opcode from program
func (cpu *CPU) EmulateCycle() {
	cpu.opcode = uint16(cpu.memory[cpu.pc])<<8 | uint16(cpu.memory[cpu.pc+1])
	opcode := getOpcode(cpu.opcode)
	if opcodeFunc, found := opcodes[opcode]; found {
		opcodeFunc(cpu)
	} else {
		// Invalid opcode
		fmt.Printf("Unknown opcode: 0x%X\n", opcode)
		cpu.pc += 2
	}

	if cpu.delayTimer > 0 {
		cpu.delayTimer--
	}
	if cpu.soundTimer > 0 {
		if cpu.soundTimer == 1 {
			fmt.Printf("BEEP!\n")
		}
		cpu.soundTimer--
	}
}

// gets the correct mapping for the lookup table
func getOpcode(opcode uint16) uint16 {
	prefix := opcode & 0xF000
	var newOpcode uint16
	if prefix == 0x0000 {
		suffix := opcode & 0x000F
		if suffix == 0x000E {
			newOpcode = prefix | suffix
		} else {
			newOpcode = prefix
		}
	} else if prefix == 0x8000 {
		suffix := opcode & 0x000F
		newOpcode = prefix | suffix
	} else if prefix == 0xE000 || prefix == 0xF000 {
		suffix := opcode & 0x00FF
		newOpcode = prefix | suffix
	} else {
		newOpcode = prefix
	}
	return newOpcode
}

// opcode execution
func (cpu *CPU) execute00E0() {
	//clear screen
	cpu.Gfx = make([]uint8, 2048)
	cpu.DrawFlag = true
	cpu.pc += 2
}
func (cpu *CPU) execute00EE() {
	//return from subroutine
	cpu.sp--
	cpu.pc = cpu.stack[cpu.sp]
	cpu.pc += 2
}
func (cpu *CPU) execute1NNN() {
	//jumps to
	cpu.pc = cpu.opcode & 0x0FFF
}
func (cpu *CPU) execute2NNN() {
	//calls subroutine
	cpu.stack[cpu.sp] = cpu.pc
	cpu.sp++
	cpu.pc = cpu.opcode & 0x0FFF
}
func (cpu *CPU) execute3XNN() {
	if cpu.regV[(cpu.opcode&0x0F00)>>8] == uint8(cpu.opcode&0x00FF) {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}
func (cpu *CPU) execute4XNN() {
	if cpu.regV[(cpu.opcode&0x0F00)>>8] != uint8(cpu.opcode&0x00FF) {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}
func (cpu *CPU) execute5XY0() {
	if cpu.regV[(cpu.opcode&0x0F00)>>8] == cpu.regV[(cpu.opcode&0x00F0)>>4] {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}
func (cpu *CPU) execute6XNN() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] = uint8(cpu.opcode & 0x00FF)
	cpu.pc += 2
}
func (cpu *CPU) execute7XNN() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] += uint8(cpu.opcode & 0x00FF)
	cpu.pc += 2
}
func (cpu *CPU) execute8XY0() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.regV[(cpu.opcode&0x00F0)>>4]
	cpu.pc += 2

}
func (cpu *CPU) execute8XY1() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.regV[(cpu.opcode&0x0F00)>>8] | cpu.regV[(cpu.opcode&0x00F0)>>4]
	cpu.pc += 2
}
func (cpu *CPU) execute8XY2() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.regV[(cpu.opcode&0x0F00)>>8] & cpu.regV[(cpu.opcode&0x00F0)>>4]
	cpu.pc += 2
}
func (cpu *CPU) execute8XY3() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.regV[(cpu.opcode&0x0F00)>>8] ^ cpu.regV[(cpu.opcode&0x00F0)>>4]
	cpu.pc += 2
}
func (cpu *CPU) execute8XY4() {
	tempX := cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.regV[(cpu.opcode&0x0F00)>>8] += cpu.regV[(cpu.opcode&0x00F0)>>4]
	if (uint16(tempX) + uint16(cpu.regV[(cpu.opcode&0x00F0)>>4])) > 255 {
		cpu.regV[0xF] = 1
	} else {
		cpu.regV[0xF] = 0
	}
	cpu.pc += 2
}
func (cpu *CPU) execute8XY5() {
	tempX := cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.regV[(cpu.opcode&0x0F00)>>8] -= cpu.regV[(cpu.opcode&0x00F0)>>4]
	if tempX >= cpu.regV[(cpu.opcode&0x00F0)>>4] {
		cpu.regV[0xF] = 1
	} else {
		cpu.regV[0xF] = 0
	}
	cpu.pc += 2
}
func (cpu *CPU) execute8XY6() {
	tempX := cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.regV[(cpu.opcode&0x0F00)>>8] = tempX >> 1
	cpu.regV[0xF] = tempX & 1
	cpu.pc += 2
}
func (cpu *CPU) execute8XY7() {
	tempX := cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.regV[(cpu.opcode&0x00F0)>>4] - tempX
	if cpu.regV[(cpu.opcode&0x00F0)>>4] >= tempX {
		cpu.regV[0xF] = 1
	} else {
		cpu.regV[0xF] = 0
	}
	cpu.pc += 2
}
func (cpu *CPU) execute8XYE() {
	tempX := cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.regV[(cpu.opcode&0x0F00)>>8] << 1
	cpu.regV[0xF] = tempX >> 7
	cpu.pc += 2
}
func (cpu *CPU) execute9XY0() {
	if cpu.regV[(cpu.opcode&0x0F00)>>8] != cpu.regV[(cpu.opcode&0x00F0)>>4] {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}
func (cpu *CPU) executeANNN() {
	cpu.regI = cpu.opcode & 0x0FFF
	cpu.pc += 2
}
func (cpu *CPU) executeBNNN() {
	cpu.pc = cpu.opcode&0x0FFF + uint16(cpu.regV[0])
}
func (cpu *CPU) executeCXNN() {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	cpu.regV[(cpu.opcode&0x0F00)>>8] = uint8(r.Intn(255)) & 0x00FF
	cpu.pc += 2
}
func (cpu *CPU) executeDXYN() {
	x := uint16(cpu.regV[(cpu.opcode&0x0F00)>>8])
	y := uint16(cpu.regV[(cpu.opcode&0x00F0)>>4])
	height := cpu.opcode & 0x000F
	var pixel uint16
	var yline uint16
	var xline uint16
	cpu.regV[0xF] = 0
	for yline = 0; yline < height; yline++ {
		pixel = uint16(cpu.memory[cpu.regI+yline])
		for xline = 0; xline < 8; xline++ {
			if pixel&(0x80>>xline) != 0 {
				gfxIndex := (int(x)+int(xline))%64 + (((int(y) + int(yline)) % 32) * 64)
				if cpu.Gfx[gfxIndex] == 1 {
					cpu.regV[0xF] = 1
				}
				cpu.Gfx[gfxIndex] ^= 1
			}
		}
	}
	cpu.DrawFlag = true
	cpu.pc += 2
}
func (cpu *CPU) executeEX9E() {
	if cpu.Key[cpu.regV[(cpu.opcode&0x0F00)>>8]] != 0 {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}
func (cpu *CPU) executeEXA1() {
	if cpu.Key[cpu.regV[(cpu.opcode&0x0F00)>>8]] == 0 {
		cpu.pc += 4
	} else {
		cpu.pc += 2
	}
}
func (cpu *CPU) executeFX07() {
	cpu.regV[(cpu.opcode&0x0F00)>>8] = cpu.delayTimer
	cpu.pc += 2
}
func (cpu *CPU) executeFX0A() {
	keyPress := false
	for i := 0; i < 16; i++ {
		if cpu.Key[i] != 0 {
			cpu.regV[(cpu.opcode&0x0F00)>>8] = uint8(i)
			keyPress = true
		}
	}
	if keyPress {
		cpu.pc += 2
	}
}
func (cpu *CPU) executeFX15() {
	cpu.delayTimer = cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.pc += 2
}
func (cpu *CPU) executeFX18() {
	cpu.soundTimer = cpu.regV[(cpu.opcode&0x0F00)>>8]
	cpu.pc += 2
}
func (cpu *CPU) executeFX1E() {
	cpu.regI += uint16(cpu.regV[(cpu.opcode&0x0F00)>>8])
	cpu.pc += 2
}
func (cpu *CPU) executeFX29() {
	cpu.regI = uint16(cpu.regV[(cpu.opcode&0x0F00)>>8] * 0x05)
	cpu.pc += 2
}
func (cpu *CPU) executeFX33() {
	cpu.memory[cpu.regI] = cpu.regV[(cpu.opcode&0x0F00)>>8] / 100
	cpu.memory[cpu.regI+1] = (cpu.regV[(cpu.opcode&0x0F00)>>8] / 10) % 10
	cpu.memory[cpu.regI+2] = (cpu.regV[(cpu.opcode&0x0F00)>>8] % 100) % 10
	cpu.pc += 2
}
func (cpu *CPU) executeFX55() {
	var i uint16
	for i = 0; i <= ((cpu.opcode & 0x0F00) >> 8); i++ {
		cpu.memory[cpu.regI+i] = cpu.regV[i]
	}
	cpu.regI += uint16(((cpu.opcode & 0x0F00) >> 8) + 1)
	cpu.pc += 2
}
func (cpu *CPU) executeFX65() {
	var i uint16
	for i = 0; i <= (cpu.opcode&0x0F00)>>8; i++ {
		cpu.regV[i] = cpu.memory[cpu.regI+i]
	}
	cpu.regI += uint16(((cpu.opcode & 0x0F00) >> 8) + 1)
	cpu.pc += 2
}
