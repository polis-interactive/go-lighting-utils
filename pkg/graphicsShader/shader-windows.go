//go:build windows
// +build windows

package graphicsShader

import (
	"github.com/go-gl/gl/v2.1/gl"
	"io/ioutil"
)

func newShaderFromFile(file string, sType uint32) (*shader, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(string(src) + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err = getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::"+file)
	if err != nil {
		return nil, err
	}
	return &shader{handle: handle}, nil
}

func (s *shader) delete() {
	gl.DeleteShader(s.handle)
}
