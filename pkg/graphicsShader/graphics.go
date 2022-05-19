package graphicsShader

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
)

type ShaderKey string

type ShaderIdentifier struct {
	Key      ShaderKey
	Filename string
}

type ShaderIdentifiers []ShaderIdentifier

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

func (gs *GraphicsShader) AttachShader(id ShaderIdentifier) error {
	if _, ok := gs.programs[id.Key]; ok {
		return errors.New(fmt.Sprintf("shader with key %s already exists", id.Key))
	}
	qualifiedPath, err := GetQualifiedShaderPath(gs.shaderPath, id.Filename)
	if err != nil {
		return err
	}
	p, err := newProgram(gs, qualifiedPath)
	if err != nil {
		return err
	}
	if gs.currentShader == "" {
		gs.currentShader = id.Key
	}
	gs.programs[id.Key] = p
	return nil
}

func (gs *GraphicsShader) AttachShaders(si ShaderIdentifiers) error {
	for _, s := range si {
		err := gs.AttachShader(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *GraphicsShader) SetShader(key ShaderKey) error {
	// might need to lock here
	if _, ok := gs.programs[key]; !ok {
		return errors.New(fmt.Sprintf("couldn't find shader with key %s", key))
	}
	gs.currentShader = key
	return nil
}

func (gs *GraphicsShader) ReloadShader() error {
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
