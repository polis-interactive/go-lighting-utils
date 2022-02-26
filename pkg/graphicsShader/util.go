package graphicsShader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func checkProgramPath(programName string) error {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	if !strings.Contains(basePath, programName) {
		return errors.New(fmt.Sprintf("PATH DOES NOT INCLUDE PROGRAM %s", programName))
	}
	return nil
}

func getQualifiedShaderPath(programName string, shaderName string) (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	dataPath := strings.Split(basePath, programName)[0]
	dataPath = filepath.Join(dataPath, programName, "data", shaderName)
	fragPath := dataPath + ".frag"
	if _, err := os.Stat(fragPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("COULDN'T FIND FRAGMENT SHADER")
	}
	vertPath := dataPath + ".vert"
	if _, err := os.Stat(vertPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("COULDN'T FIND VERTEX SHADER")
	}
	return dataPath, nil
}
