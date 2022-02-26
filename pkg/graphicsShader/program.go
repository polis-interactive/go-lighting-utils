package graphicsShader

type program struct {
	gs         *GraphicsShader
	handle     uint32
	shaders    []shader
	shaderFile string
}
