//go:build windows
// +build windows

package graphicsShader

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
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
		handle:     gl.CreateProgram(),
		shaderFile: shaderFile,
	}

	err = p.loadShaders()

	return p, nil
}

func (p *program) loadShaders() error {
	vertexShader, err := newShaderFromFile(p.shaderFile+".vert", gl.VERTEX_SHADER)
	if err != nil {
		log.Println(fmt.Sprintf("Couldn't compile vertex shader %s", p.shaderFile))
		return err
	}
	fragmentShader, err := newShaderFromFile(p.shaderFile+".frag", gl.FRAGMENT_SHADER)
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

func (p *program) attach(shaders ...shader) {
	for _, s := range shaders {
		gl.AttachShader(p.handle, s.handle)
		p.shaders = append(p.shaders, s)
	}
}

func (p *program) link() error {
	gl.LinkProgram(p.handle)
	return getGlError(p.handle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE")
}

func (p *program) use() {
	gl.UseProgram(p.handle)
}

func (p *program) setUniform1f(name UniformKey, value float32) {
	chars := []uint8(name)
	loc := gl.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 1f %s", name))
		return
	}
	gl.Uniform1f(loc, value)
}

func (p *program) setUniform2fv(name string, value []float32, count int32) {
	chars := []uint8(name)
	loc := gl.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 2f %s", name))
		return
	}
	gl.Uniform2fv(loc, count, &value[0])
}

func (p *program) delete() {
	for _, s := range p.shaders {
		s.delete()
	}
	gl.DeleteProgram(p.handle)
}

func (p *program) runProgram() error {
	p.use()
	p.setUniform2fv("resolution", []float32{p.gs.widthF, p.gs.heightF}, 1)
	p.gs.mu.RLock()
	for u, v := range p.gs.uniformDict {
		p.setUniform1f(u, v)
	}
	p.gs.mu.RUnlock()
	gl.BindVertexArray(p.gs.rectHandle)
	gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, unsafe.Pointer(nil))
	gl.BindVertexArray(0)
	// should probably check for an error here, not sure what tho
	return nil
}
