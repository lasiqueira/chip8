package main

import (
	"os"

	"chip8/cpu"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/vbsw/glut"
)

var chip8 = cpu.CPU{}
var displayWidth = 640
var displayHeight = 320
var modifier = 10

func initWindow() {
	gl.Init()
	glut.Init()
	glut.InitDisplayMode(glut.DOUBLE | glut.RGBA)
	glut.InitWindowSize(displayWidth, displayHeight)
	glut.InitWindowPosition(320, 320)
	glut.CreateWindow("chip8")
	glut.ReshapeFunc(reshapeWindow)
	glut.DisplayFunc(display)
	glut.IdleFunc(display)
	glut.KeyboardFunc(keyboardDown)
	glut.KeyboardUpFunc(keyboardUp)
	glut.MainLoop()

}
func updateQuads() {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if chip8.Gfx[(y*64)+x] == 0 {
				gl.Color3f(0.0, 0.0, 0.0)
			} else {
				gl.Color3f(1.0, 1.0, 1.0)
			}
			drawPixel(x, y)
		}
	}
}
func drawPixel(x int, y int) {
	gl.Begin(gl.QUADS)
	gl.Vertex3f(float32(x*modifier), float32(y*modifier), 0.0)
	gl.Vertex3f(float32(x*modifier), float32((y*modifier)+modifier), 0.0)
	gl.Vertex3f(float32((x*modifier)+modifier), float32((y*modifier)+modifier), 0.0)
	gl.Vertex3f(float32((x*modifier)+modifier), float32(y*modifier), 0.0)
	gl.End()
}
func display() {
	chip8.EmulateCycle()
	if chip8.DrawFlag {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		updateQuads()
		glut.SwapBuffers()
		chip8.DrawFlag = false
	}
}
func keyboardDown(key uint8, x int, y int) {
	if key == 27 {
		os.Exit(0)
	}

	switch key {
	case '1':
		chip8.Key[0x1] = 1
		break
	case '2':
		chip8.Key[0x2] = 1
		break
	case '3':
		chip8.Key[0x3] = 1
		break
	case '4':
		chip8.Key[0xC] = 1
		break
	case 'q':
		chip8.Key[0x4] = 1
		break
	case 'w':
		chip8.Key[0x5] = 1
		break
	case 'e':
		chip8.Key[0x6] = 1
		break
	case 'r':
		chip8.Key[0xD] = 1
		break
	case 'a':
		chip8.Key[0x7] = 1
		break
	case 's':
		chip8.Key[0x8] = 1
		break
	case 'd':
		chip8.Key[0x9] = 1
		break
	case 'f':
		chip8.Key[0xE] = 1
		break
	case 'z':
		chip8.Key[0xA] = 1
		break
	case 'x':
		chip8.Key[0x0] = 1
		break
	case 'c':
		chip8.Key[0xB] = 1
		break
	case 'v':
		chip8.Key[0xF] = 1
		break
	}
}
func keyboardUp(key uint8, x int, y int) {
	if key == 27 {
		os.Exit(0)
	}

	switch key {
	case '1':
		chip8.Key[0x1] = 0
		break
	case '2':
		chip8.Key[0x2] = 0
		break
	case '3':
		chip8.Key[0x3] = 0
		break
	case '4':
		chip8.Key[0xC] = 0
		break
	case 'q':
		chip8.Key[0x4] = 0
		break
	case 'w':
		chip8.Key[0x5] = 0
		break
	case 'e':
		chip8.Key[0x6] = 0
		break
	case 'r':
		chip8.Key[0xD] = 0
		break
	case 'a':
		chip8.Key[0x7] = 0
		break
	case 's':
		chip8.Key[0x8] = 0
		break
	case 'd':
		chip8.Key[0x9] = 0
		break
	case 'f':
		chip8.Key[0xE] = 0
		break
	case 'z':
		chip8.Key[0xA] = 0
		break
	case 'x':
		chip8.Key[0x0] = 0
		break
	case 'c':
		chip8.Key[0xB] = 0
		break
	case 'v':
		chip8.Key[0xF] = 0
		break
	}
}
func reshapeWindow(w int, h int) {
	gl.ClearColor(0.0, 0.0, 0.5, 0.0)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0.0, float64(w), float64(h), 0.0, -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.Viewport(0, 0, int32(w), int32(h))

	// Resize quad
	displayWidth = w
	displayHeight = h
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Need to inform the game path")
	}
	game := args[0]

	chip8.Initialize()
	chip8.LoadGame(game)
	initWindow()

	return
}
