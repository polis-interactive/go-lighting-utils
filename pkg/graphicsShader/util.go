package graphicsShader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetShaderPathIfAvailable(programName string) (string, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return "", errors.New("COULDN'T GET CWD")
	}
	if !strings.Contains(basePath, programName) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT INCLUDE PROGRAM %s", programName))
	}
	dataPath := strings.Split(basePath, programName)[0]
	dataPath = filepath.Join(dataPath, programName, "data")
	if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT EXIST: %s", dataPath))
	}
	shaderPath := filepath.Join(dataPath, "shaders")
	if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT EXIST: %s", shaderPath))
	}

	return shaderPath, nil
}

func GetQualifiedShaderPath(shaderPath string, shaderName string) (string, error) {
	basePath := filepath.Join(shaderPath, shaderName)
	fragPath := basePath + ".frag"
	if _, err := os.Stat(fragPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("COULDN'T FIND FRAGMENT SHADER")
	}
	vertPath := basePath + ".vert"
	if _, err := os.Stat(vertPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("COULDN'T FIND VERTEX SHADER")
	}
	return basePath, nil
}
