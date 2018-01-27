package cpu

import (
	"os"

	"encoding/binary"

	"fmt"

	"github.com/lasiqueira/chip8/util"
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
}

//Initialize sets up the initial state of the CPU
func (cpu *CPU) Initialize() {
	cpu.pc = 0x200
	cpu.opcode = 0
	cpu.regI = 0
	cpu.memory = make([]uint8, 4096)
	cpu.regV = make([]uint8, 16)
	cpu.gfx = make([]uint8, 2048)
	cpu.delayTimer = 60
	cpu.soundTimer = 60
	cpu.stack = make([]uint16, 16)
	cpu.sp = 0
	cpu.key = make([]uint8, 16)
	//loadfontset
	for i := 0; i < 80; i++ {
		//set font
		//cpu.memory[i] =
	}

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
