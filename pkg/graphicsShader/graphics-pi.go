//go:build !windows
// +build !windows

package graphicsShader

import (
	"fmt"
	"github.com/go-gl/gl/v3.1/gles2"
	"strings"
	"unsafe"
)

func glInit() error {
	return gles2.Init()
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
	gles2.GenVertexArrays(1, &VAO)

	var VBO uint32
	gles2.GenBuffers(1, &VBO)

	var EBO uint32
	gles2.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gles2.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gles2.BindBuffer(gles2.ARRAY_BUFFER, VBO)
	gles2.BufferData(gles2.ARRAY_BUFFER, len(vertices)*4, gles2.Ptr(vertices), gles2.STATIC_DRAW)

	// copy indices into element buffer
	gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, EBO)
	gles2.BufferData(gles2.ELEMENT_ARRAY_BUFFER, len(indices)*4, gles2.Ptr(indices), gles2.STATIC_DRAW)

	// position
	gles2.VertexAttribPointer(0, 3, gles2.FLOAT, false, 6*4, gles2.PtrOffset(0))
	gles2.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gles2.BindVertexArray(0)

	// should probably check for an error here, not sure what tho

	return VAO
}

func (gs *GraphicsShader) ClearBuffer() {
	gles2.ClearColor(0.2, 0.2, 0.2, 1.0)
	gles2.Clear(gles2.COLOR_BUFFER_BIT)
}

func (gs *GraphicsShader) ReadToPixels(pb unsafe.Pointer) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gles2.ReadPixels(0, 0, gs.width, gs.height, gles2.RGBA, gles2.UNSIGNED_BYTE, pb)
	// should probably check for error
	return nil
}

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gles2.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gles2.INFO_LOG_LENGTH, &logLength)

		outMsg := gles2.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, outMsg)

		return fmt.Errorf("%s: %s", failMsg, gles2.GoStr(outMsg))
	}

	return nil
}
