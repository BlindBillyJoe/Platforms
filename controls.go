package main

import "github.com/go-gl/glfw/v3.2/glfw"

//MainKeyboardController instance
var MainKeyboardController keyboardController

type keyboardController struct {
	units []*object
}

func (controller *keyboardController) activate(w *Window) {
	w.setKeyCallback(keyboardCallback)
}

func keyboardCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyRight {
		for i := range MainKeyboardController.units {
			MainKeyboardController.units[i].speedX = speed * float32(action)
		}
	}

	if key == glfw.KeyLeft {
		for i := range MainKeyboardController.units {
			MainKeyboardController.units[i].speedX = -speed * float32(action)
		}
	}

	if key == glfw.KeyUp {
		for i := range MainKeyboardController.units {
			MainKeyboardController.units[i].speedY = speed * float32(action)
		}
	}

	if key == glfw.KeyDown {
		for i := range MainKeyboardController.units {
			MainKeyboardController.units[i].speedY = -speed * float32(action)
		}
	}
}
