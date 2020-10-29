package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

//Window abstraction struct
type Window struct {
	w *glfw.Window
}

func (w *Window) setKeyCallback(callback glfw.KeyCallback) glfw.KeyCallback {
	return w.w.SetKeyCallback(callback)
}

func (w *Window) shouldClose() bool {
	return w.w.ShouldClose()
}

func (w *Window) clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (w *Window) process() {
	w.w.SwapBuffers()
}

func (w *Window) pollEvents() {
	glfw.PollEvents()
}

func initWindow(w, h int, title string) *Window {
	check(glfw.Init())
	var glfwWindow, err = glfw.CreateWindow(800, 600, "Test", nil, nil)
	check(err)
	glfwWindow.MakeContextCurrent()
	glfw.SwapInterval(1)
	check(gl.Init())

	var window = &Window{
		w: glfwWindow,
	}

	return window
}

func (obj *object) makeVao() {
	gl.GenBuffers(1, &obj.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, obj.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(obj._data), gl.Ptr(obj._data), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &obj.vao)
	gl.BindVertexArray(obj.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, obj.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	var shader = gl.CreateShader(shaderType)

	var csources, free = gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		var log = strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func draw(obj *object, w *Window) {
	gl.UseProgram(obj.program)

	gl.BindVertexArray(obj.vao)
	gl.Uniform3f(gl.GetUniformLocation(obj.program, gl.Str("pos\x00")), obj.pos.x, obj.pos.y, obj.pos.z)
	gl.Uniform3f(gl.GetUniformLocation(obj.program, gl.Str("color\x00")), obj.col.x, obj.col.y, obj.col.z)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(obj._data)))
}
