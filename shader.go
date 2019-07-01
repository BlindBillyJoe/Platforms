package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type shader struct {
	shaderType uint32
	source     string
	compiled   uint32
}

const (
	VERTEX_SHADER   = gl.VERTEX_SHADER
	FRAGMENT_SHADER = gl.FRAGMENT_SHADER
)

func createShaderFromFile(filename string, shaderType uint32) *shader {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()
	var builder strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		builder.WriteString(scanner.Text())
		builder.WriteByte('\n')
	}
	// builder.WriteByte('\x00')
	return createShader(builder.String(), shaderType)
}

func createShader(source string, shaderType uint32) *shader {
	return &shader{
		shaderType: shaderType,
		source:     source,
	}
}

func (s *shader) compile() error {
	if s.compiled != 0 {
		return fmt.Errorf("Shader already compiled")
	}

	var err error
	s.compiled, err = compileShader(s.source, s.shaderType)
	return err
}
