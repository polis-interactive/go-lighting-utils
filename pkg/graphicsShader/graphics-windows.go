//go:build windows
// +build windows

package graphicsShader

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"strings"
	"unsafe"
)

func glInit() error {
	return gl.Init()
}

func createFillRect() uint32 {

	vertices := []float32{
		// bottom left
		-1.0, -1.0, 0.0, // position
		0.0, 0.0, 0.0, // Color

		// bottom right
		1.0, -1.0, 0.0,
		0.0, 0.0, 0.0,

		// top right
		1.0, 1.0, 0.0,
		0.0, 0.0, 0.0,

		// top left
		-1.0, 1.0, 0.0,
		0.0, 0.0, 0.0,
	}

	indices := []uint32{
		0, 1, 2, 3,
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	// should probably check for an error here, not sure what tho

	return VAO
}

func (gs *GraphicsShader) ClearBuffer() {
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (gs *GraphicsShader) ReadToPixels(pb unsafe.Pointer) error {
	gl.ReadPixels(0, 0, gs.width, gs.height, gl.RGBA, gl.UNSIGNED_BYTE, pb)
	// should probably check for error
	return nil
}

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		outMsg := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, outMsg)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(outMsg))
	}

	return nil
}
