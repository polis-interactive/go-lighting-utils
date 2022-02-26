package graphicsShader

import "github.com/go-gl/glfw/v3.3/glfw"

type windowProxy struct {
	*glfw.Window
}

func glfwInit() error {
	return glfw.Init()
}

func glfwTerminate() {
	glfw.Terminate()
}

func glfwPollEvents() {
	glfw.PollEvents()
}

func newWindow(width int32, height int32) (*windowProxy, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(int(width), int(height), "Shader Window", nil, nil)
	if err != nil {
		return nil, err
	}
	return &windowProxy{window}, nil
}

func windowKeyCallback(
	window *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey,
) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}
