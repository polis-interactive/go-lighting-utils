//go:build !windows
// +build !windows

package graphicsShader

import (
	"fmt"
	"github.com/go-gl/gl/v3.1/gles2"
	"log"
	"unsafe"
)

func newProgram(gs *GraphicsShader, shaderFile string) (p *program, err error) {

	defer func() {
		if err != nil && p != nil {
			p.delete()
		}
	}()

	p = &program{
		gs:         gs,
		handle:     gles2.CreateProgram(),
		shaderFile: shaderFile,
	}

	err = p.loadShaders()

	return p, nil
}

func (p *program) loadShaders() error {
	vertexShader, err := newShaderFromFile(p.shaderFile+".vert", gles2.VERTEX_SHADER)
	if err != nil {
		log.Println(fmt.Sprintf("Couldn't compile vertex shader %s", p.shaderFile))
		return err
	}
	fragmentShader, err := newShaderFromFile(p.shaderFile+".frag", gles2.FRAGMENT_SHADER)
	if err != nil {
		log.Println(fmt.Sprintf("Couldn't compile fragment shader %s", p.shaderFile))
		return err
	}

	p.attach(*vertexShader, *fragmentShader)

	err = p.link()
	if err != nil {
		log.Println(fmt.Sprintf("Couldn't link shaders %s", p.shaderFile))
		return err
	}
	return nil
}

func (p *program) reloadShaders() error {
	p.delete()
	p.handle = gles2.CreateProgram()
	return p.loadShaders()
}

func (p *program) attach(shaders ...shader) {
	for _, s := range shaders {
		gles2.AttachShader(p.handle, s.handle)
		p.shaders = append(p.shaders, s)
	}
}

func (p *program) link() error {
	gles2.LinkProgram(p.handle)
	return getGlError(p.handle, gles2.LINK_STATUS, gles2.GetProgramiv, gles2.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE")
}

func (p *program) use() {
	gles2.UseProgram(p.handle)
}

func (p *program) setUniform1f(name UniformKey, value float32) {
	chars := []uint8(name)
	loc := gles2.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 1f %s", name))
		return
	}
	gles2.Uniform1f(loc, value)
}

func (p *program) setUniform2fv(name string, value []float32, count int32) {
	chars := []uint8(name)
	loc := gles2.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 2f %s", name))
		return
	}
	gles2.Uniform2fv(loc, count, &value[0])
}

func (p *program) delete() {
	for _, s := range p.shaders {
		s.delete()
	}
	gles2.DeleteProgram(p.handle)
}

func (p *program) runProgram() error {
	p.use()
	p.setUniform2fv("resolution", []float32{p.gs.widthF, p.gs.heightF}, 1)
	p.gs.mu.RLock()
	for u, v := range p.gs.uniformDict {
		p.setUniform1f(u, v)
	}
	p.gs.mu.RUnlock()
	gles2.BindVertexArray(p.gs.rectHandle)
	gles2.DrawElements(gles2.TRIANGLE_FAN, 4, gles2.UNSIGNED_INT, unsafe.Pointer(nil))
	gles2.BindVertexArray(0)
	// should probably check for an error here, not sure what tho
	return nil
}
