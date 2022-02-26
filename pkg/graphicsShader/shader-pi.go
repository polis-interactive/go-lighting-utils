//go:build !windows
// +build !windows

package graphicsShader

import (
	"github.com/go-gl/gl/v3.1/gles2"
	"io/ioutil"
)

func newShaderFromFile(file string, sType uint32) (*shader, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	handle := gles2.CreateShader(sType)
	glSrc, freeFn := gles2.Strs(string(src) + "\x00")
	defer freeFn()
	gles2.ShaderSource(handle, 1, glSrc, nil)
	gles2.CompileShader(handle)
	err = getGlError(handle, gles2.COMPILE_STATUS, gles2.GetShaderiv, gles2.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::"+file)
	if err != nil {
		return nil, err
	}
	return &shader{handle: handle}, nil
}

func (s *shader) delete() {
	gles2.DeleteShader(s.handle)
}
