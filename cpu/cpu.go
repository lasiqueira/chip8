package cpu

import (
	"os"

	"encoding/binary"

	"fmt"

	"math/rand"

	"github.com/lasiqueira/chip8/util"

	"time"
)

//CPU struct holds the cpu state during emulation
type CPU struct {
	opcode     uint16
	memory     []uint8
	regV       []uint8
	regI       uint16
	pc         uint16
	gfx        []uint8
	delayTimer uint8
	soundTimer uint8
	stack      []uint16
	sp         uint16
	key        []uint8
	drawFlag   bool
}

//Initialize sets up the initial state of the CPU
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
	cpu.gfx = make([]uint8, 2048)
	cpu.delayTimer = 0
	cpu.soundTimer = 0
	cpu.stack = make([]uint16, 16)
	cpu.sp = 0
	cpu.key = make([]uint8, 16)
	//loadfontset
	for i := 0; i < 80; i++ {
		//set font
		cpu.memory[i] = font[i]
	}
	cpu.drawFlag = true
}

//LoadGame loads the game into memory and reads it
func (cpu *CPU) LoadGame(game string) {
	f, err := os.OpenFile(game, os.O_RDONLY, 0777)
	util.HandleError(err)
	fStat, err := f.Stat()
	b := make([]byte, 1)
	for i := 0; int64(i) < fStat.Size(); i++ {
		data, err := f.Read(b)
		util.HandleError(err)
		cpu.memory[i+512] = uint8(data)
	}
	defer f.Close()

}

//EmulateCycle fetch, decode and execute opcode from program
func (cpu *CPU) EmulateCycle() {
	cpu.opcode = binary.BigEndian.Uint16(cpu.memory[cpu.pc : cpu.pc+2])

	switch cpu.opcode & 0xF000 {
	case 0x0000:
		switch cpu.opcode & 0x000F {
		case 0x0000:
			//clear screen
			cpu.gfx = make([]uint8, 2048)
			cpu.drawFlag = true
			cpu.pc += 2
			break
		case 0x000E:
			//return from subroutine
			cpu.pc = cpu.stack[cpu.sp]
			cpu.sp--
			cpu.pc += 2
			break
		default:
			fmt.Printf("Unknown opcode [0x0000]: 0x%X\n", cpu.opcode)
		}
	case 0x1000:
		//jumps to
		cpu.pc = cpu.opcode & 0x0FFF
		break
	case 0x2000:
		//calls subroutine
		cpu.stack[cpu.sp] = cpu.pc
		cpu.sp++
		cpu.pc = cpu.opcode & 0x0FFF
		break
	case 0x3000:
		if cpu.regV[uint8(cpu.opcode&0x0F00)] == uint8(cpu.opcode&0x00FF) {
			cpu.pc += 4
		} else {
			cpu.pc += 2
		}
		break
	case 0x4000:
		if cpu.regV[uint8(cpu.opcode&0x0F00)] != uint8(cpu.opcode&0x00FF) {
			cpu.pc += 4
		} else {
			cpu.pc += 2
		}
		break
	case 0x5000:
		if cpu.regV[uint8(cpu.opcode&0x0F00)] == cpu.regV[uint8(cpu.opcode&0x00F0)] {
			cpu.pc += 4
		} else {
			cpu.pc += 2
		}
		break
	case 0x6000:
		cpu.regV[uint8(cpu.opcode&0x0F00)] = uint8(cpu.opcode & 0x00FF)
		cpu.pc += 2
		break
	case 0x7000:
		cpu.regV[uint8(cpu.opcode&0x0F00)] += uint8(cpu.opcode & 0x00FF)
		cpu.pc += 2
		break
	case 0x8000:
		switch cpu.opcode & 0x000F {
		case 0x0000:
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0001:
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x0F00)] | cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0002:
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x0F00)] & cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0003:
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x0F00)] ^ cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0004:
			if (uint16(cpu.regV[uint8(cpu.opcode&0x0F00)]) + uint16(cpu.regV[uint8(cpu.opcode&0x00F0)])) > 255 {
				cpu.regV[0xF] = 1
			} else {
				cpu.regV[0xF] = 0
			}
			cpu.regV[uint8(cpu.opcode&0x0F00)] += cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0005:
			if cpu.regV[uint8(cpu.opcode&0x0F00)] >= cpu.regV[uint8(cpu.opcode&0x00F0)] {
				cpu.regV[0xF] = 1
			} else {
				cpu.regV[0xF] = 0
			}
			cpu.regV[uint8(cpu.opcode&0x0F00)] -= cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0006:
			cpu.regV[0xF] = cpu.regV[uint8(cpu.opcode&0x00F0)] & 1
			cpu.regV[uint8(cpu.opcode&0x00F0)] = cpu.regV[uint8(cpu.opcode&0x00F0)] >> 1
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		case 0x0007:
			if cpu.regV[uint8(cpu.opcode&0x00F0)] >= cpu.regV[uint8(cpu.opcode&0x0F00)] {
				cpu.regV[0xF] = 1
			} else {
				cpu.regV[0xF] = 0
			}
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x00F0)] - cpu.regV[uint8(cpu.opcode&0x0F00)]
			cpu.pc += 2
			break
		case 0x000E:
			cpu.regV[0xF] = cpu.regV[uint8(cpu.opcode&0x00F0)] & 128
			cpu.regV[uint8(cpu.opcode&0x00F0)] = cpu.regV[uint8(cpu.opcode&0x00F0)] << 1
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.regV[uint8(cpu.opcode&0x00F0)]
			cpu.pc += 2
			break
		default:
			fmt.Printf("Unknown opcode [0x8000]: 0x%X\n", cpu.opcode)
		}
	case 0x9000:
		if cpu.regV[uint8(cpu.opcode&0x0F00)] != cpu.regV[uint8(cpu.opcode&0x00F0)] {
			cpu.pc += 4
		} else {
			cpu.pc += 2
		}
		break
	case 0xA000:
		cpu.regI = cpu.opcode & 0x0FFF
		cpu.pc += 2
		break
	case 0xB000:
		cpu.pc = cpu.opcode&0x0FFF + uint16(cpu.regV[0])
		break
	case 0xC000:
		rand.Seed(time.Now().UTC().UnixNano())
		cpu.regV[cpu.opcode&0x0f00] = uint8(rand.Intn(255)) & 0x00FF
		cpu.pc += 2
		break
	case 0xD000:
		var x uint16 = uint16(cpu.regV[uint8(cpu.opcode&0x0F00)])
		var y uint16 = uint16(cpu.regV[uint8(cpu.opcode&0x00F0)])
		var height uint16 = cpu.opcode & 0x000F
		var pixel uint16
		var yline uint16
		var xline uint16
		cpu.regV[0xF] = 0
		for yline = 0; yline < height; yline++ {
			pixel = uint16(cpu.memory[cpu.regI+yline])
			for xline = 0; xline < 8; xline++ {
				if pixel&(0x80>>xline) != 0 {
					if cpu.gfx[x+xline+((y+yline)*64)] == 1 {
						cpu.regV[0xF] = 1
					}
					cpu.gfx[x+xline+((y+yline)*64)] ^= 1
				}
			}
		}
		cpu.drawFlag = true
		cpu.pc += 2
		break
	case 0xE000:
		switch cpu.opcode & 0x00FF {
		case 0x009E:
			if cpu.key[cpu.regV[uint8(cpu.opcode&0x0F00)]] != 0 {
				cpu.pc += 4
			} else {
				cpu.pc += 2
			}
			break
		case 0x00A1:
			if cpu.key[cpu.regV[uint8(cpu.opcode&0x0F00)]] == 0 {
				cpu.pc += 4
			} else {
				cpu.pc += 2
			}
			break
		default:
			fmt.Printf("Unknown opcode [0xE000]: 0x%X\n", cpu.opcode)
		}
	case 0xF000:
		switch cpu.opcode & 0x00FF {
		case 0x0007:
			cpu.regV[uint8(cpu.opcode&0x0F00)] = cpu.delayTimer
			cpu.pc += 2
			break
		case 0x000A:
			keyPress := false
			for i := 0; i < 16; i++ {
				if cpu.key[i] != 0 {
					cpu.regV[uint8(cpu.opcode&0x0F00)] = uint8(i)
					keyPress = true
				}
			}
			if keyPress {
				cpu.pc += 2
			}

			break
		case 0x0015:
			cpu.delayTimer = cpu.regV[uint8(cpu.opcode&0x0F00)]
			cpu.pc += 2
			break
		case 0x0018:
			cpu.soundTimer = cpu.regV[uint8(cpu.opcode&0x0F00)]
			cpu.pc += 2
			break
		case 0x001E:
			cpu.regI += uint16(cpu.regV[uint8(cpu.opcode&0x0F00)])
			cpu.pc += 2
			break
		case 0x0029:
			cpu.regI = uint16(cpu.regV[uint8(cpu.opcode&0x0F00)] * 0x05)
			cpu.pc += 2
			break
		case 0x0033:
			cpu.memory[cpu.regI] = cpu.regV[uint8(cpu.opcode&0x0F00)] / 100
			cpu.memory[cpu.regI+1] = (cpu.regV[uint8(cpu.opcode&0x0F00)] / 10) % 10
			cpu.memory[cpu.regI+2] = (cpu.regV[uint8(cpu.opcode&0x0F00)] % 100) % 10
			cpu.pc += 2
			break
		case 0x0055:
			var i uint8
			for i = 0; i <= uint8(cpu.opcode&0x0F00); i++ {
				cpu.memory[cpu.regI+uint16(i)] = cpu.regV[i]
			}
			cpu.regI = uint16(uint8(cpu.opcode&0x0F00) + 1)
			cpu.pc += 2
			break
		case 0x0065:
			var i uint8
			for i = 0; i <= uint8(cpu.opcode&0x0F00); i++ {
				cpu.regV[i] = cpu.memory[cpu.regI+uint16(i)]
			}
			cpu.regI = uint16(uint8(cpu.opcode&0x0F00) + 1)
			cpu.pc += 2
			break
		default:
			fmt.Printf("Unknown opcode [0xF000]: 0x%X\n", cpu.opcode)
		}

	default:
		fmt.Printf("Unknown opcode: 0x%X\n", cpu.opcode)
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
