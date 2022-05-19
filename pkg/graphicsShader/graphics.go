package graphicsShader

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
)

type ShaderKey string

type ShaderIdentifiers map[ShaderKey]string

type UniformKey string

type UniformDict map[UniformKey]float32

type GraphicsShader struct {
	shaderPath    string
	width         int32
	widthF        float32
	height        int32
	heightF       float32
	uniformDict   UniformDict
	mu            *sync.RWMutex
	window        *windowProxy
	programs      map[ShaderKey]*program
	currentShader ShaderKey
	rectHandle    uint32
}

func NewGraphicsShader(
	shaderPath string, width int32, height int32,
	uniformDict UniformDict, mu *sync.RWMutex,
) (*GraphicsShader, error) {

	// required by glfw
	runtime.LockOSThread()

	gs := &GraphicsShader{
		shaderPath:    shaderPath,
		width:         width,
		widthF:        float32(width),
		height:        height,
		heightF:       float32(height),
		uniformDict:   uniformDict,
		mu:            mu,
		programs:      make(map[ShaderKey]*program),
		currentShader: "",
	}

	err := glfwInit()
	if err != nil {
		gs.Cleanup()
		log.Fatalln("failed to initialize glfw:", err)
		return nil, err
	}

	window, err := newWindow(width, height)
	if err != nil {
		gs.Cleanup()
		log.Fatalln("failed to create glfw window:", err)
		return nil, err
	}
	window.MakeContextCurrent()
	gs.window = window

	err = glInit()
	if err != nil {
		gs.Cleanup()
		log.Fatalln("failed to create gl context:", err)
		return nil, err
	}

	gs.rectHandle = createFillRect()

	window.SetKeyCallback(windowKeyCallback)

	return gs, nil
}

func (gs *GraphicsShader) AttachShader(id ShaderKey, fileName string) error {
	if _, ok := gs.programs[id]; ok {
		return errors.New(fmt.Sprintf("shader with key %s already exists", id))
	}
	qualifiedPath, err := GetQualifiedShaderPath(gs.shaderPath, fileName)
	if err != nil {
		return err
	}
	p, err := newProgram(gs, qualifiedPath)
	if err != nil {
		return err
	}
	if gs.currentShader == "" {
		gs.currentShader = id
	}
	gs.programs[id] = p
	return nil
}

func (gs *GraphicsShader) AttachShaders(si ShaderIdentifiers) error {
	for k, v := range si {
		err := gs.AttachShader(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *GraphicsShader) SetShader(key ShaderKey) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if _, ok := gs.programs[key]; !ok {
		return errors.New(fmt.Sprintf("couldn't find shader with key %s", key))
	}
	gs.currentShader = key
	return nil
}

func (gs *GraphicsShader) ReloadShader() error {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	if gs.currentShader == "" {
		return errors.New("nothing to reload")
	}
	p, ok := gs.programs[gs.currentShader]
	if !ok {
		return errors.New(fmt.Sprintf("couldn't find shader with key %s", gs.currentShader))
	}
	return p.reloadShaders()
}

func (gs *GraphicsShader) RunShader() error {

	gs.mu.RLock()
	defer gs.mu.RUnlock()

	glfwPollEvents()
	gs.ClearBuffer()

	if gs.window.ShouldClose() {
		return errors.New("force close window")
	}

	p, ok := gs.programs[gs.currentShader]
	if !ok {
		return errors.New(fmt.Sprintf("couldn't find shader %s", gs.currentShader))
	}

	err := p.runProgram()
	if err != nil {
		return errors.New("shader failed")
	}

	gs.window.SwapBuffers()

	return nil
}

func (gs *GraphicsShader) Cleanup() {
	for _, p := range gs.programs {
		if p != nil {
			p.delete()
		}
	}

	glfwTerminate()
	runtime.UnlockOSThread()
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)
