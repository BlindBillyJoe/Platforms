package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

//ControllersManager is main controller manager
var ControllersManager = controllersManager{
	controllers: make(map[string]controller),
}

const (
	vMapSize = 29
	hMapSize = 29
)

func (m *controllersManager) process() {
	// println("-------------------")
	for i := range m.controllers {
		// println(i)
		m.controllers[i].process()
	}
}

var (
	keyboardControllerType = 1
	physicsControllerType  = 2
	mapControllerType      = 3
)

type controller interface {
	activate(arg interface{})
	add(arg interface{})
	process()
}

type controllersManager struct {
	controllers map[string]controller
}

func (m *controllersManager) create(name string, controllerType int) {
	switch controllerType {
	case keyboardControllerType:
		m.controllers[name] = &keyboardController{}
	case physicsControllerType:
		m.controllers[name] = &physicsController{}
	case mapControllerType:
		m.controllers[name] = &mapController{}
	}
}

type keyboardController struct {
	units []*object
}

func (controller *keyboardController) activate(arg interface{}) {
	arg.(*Window).setKeyCallback(keyboardCallback)
}

func (controller *keyboardController) add(arg interface{}) {
	controller.units = append(controller.units, arg.(*object))
}

func (controller *keyboardController) process() {
	for i := range controller.units {
		controller.units[i].update()
	}
}

func keyboardCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	var control = ControllersManager.controllers["1Keyboard"].(*keyboardController)

	if key == glfw.KeyRight {
		if action == glfw.Press || action == glfw.Repeat {
			for i := range control.units {
				control.units[i].xAcceleration = speed / acceleration
			}
		} else if action == glfw.Release {
			for i := range control.units {
				control.units[i].xAcceleration = 0.
			}
		}
	}

	if key == glfw.KeyLeft {
		if action == glfw.Press || action == glfw.Repeat {
			for i := range control.units {
				control.units[i].nxAcceleration = speed / acceleration
			}
		} else if action == glfw.Release {
			for i := range control.units {
				control.units[i].nxAcceleration = 0.
			}
		}
	}

	if key == glfw.KeyUp {
		if action == glfw.Press || action == glfw.Repeat {
			for i := range control.units {
				control.units[i].yAcceleration = speed / acceleration
			}
		} else if action == glfw.Release {
			for i := range control.units {
				control.units[i].yAcceleration = 0.
			}
		}
	}

	if key == glfw.KeyDown {
		if action == glfw.Press || action == glfw.Repeat {
			for i := range control.units {
				control.units[i].nyAcceleration = speed / acceleration
			}
		} else if action == glfw.Release {
			for i := range control.units {
				control.units[i].nyAcceleration = 0.
			}
		}
	}

	if key == glfw.KeySpace {
		if action == glfw.Press {
			for i := range control.units {
				var obj = control.units[i]
				println(obj.intersection)
				if obj.intersection != 0 && obj.intersection&topIntersection == 0 {
					obj.yAcceleration = speed / acceleration * 10
					if obj.intersection&rightIntersection != 0 {
						obj.speed.x = -speed
					}
					if obj.intersection&leftIntersection != 0 {
						obj.speed.x = speed
					}
				}
			}
		} else if action == glfw.Release {
			for i := range control.units {
				control.units[i].yAcceleration = 0.
			}
		}
	}
}

type physicsController struct {
	units []*object
}

func (controller *physicsController) activate(arg interface{}) {

}

func (controller *physicsController) add(arg interface{}) {
	controller.units = append(controller.units, arg.(*object))
}

func (controller *physicsController) process() {
	for i := range controller.units {
		controller.units[i].updatePhysics()
	}
}

type mapController struct {
	units    []*object
	sceneMap [vMapSize][hMapSize]*object
}

func (controller *mapController) activate(arg interface{}) {
}

func (controller *mapController) add(arg interface{}) {
	controller.units = append(controller.units, arg.(*object))
}

type trackingObject struct {
	box square
	obj *object
}

func (controller *mapController) process() {
	var track = make([]trackingObject, len(controller.units))

	for i := range controller.units {
		var target = controller.units[i]
		if target.updated {
			continue
		}
		var tl = target.colision.tl.cAdd(target.pos)
		var br = target.colision.br.cAdd(target.pos)
		track[i] = trackingObject{
			box: square{
				tl: tl,
				br: br,
			},
			obj: target,
		}
		target.updated = true
	}

	if track[0].obj == nil {
		return
	}

	controller.clearMap()

	for x := range controller.sceneMap {
		for y := range controller.sceneMap[x] {
			for i := range track {
				ptr := point{x: float32(x)/(float32(hMapSize-1)/2.) - 1., y: float32(y)/(float32(vMapSize-1)/2.) - 1.}
				if track[i].box.contains(&ptr) {
					controller.sceneMap[x][y] = track[i].obj
				}
			}
		}
	}
}

func (controller *mapController) clearMap() {
	for x := 0; x != hMapSize; x++ {
		for y := 0; y != vMapSize; y++ {
			controller.sceneMap[x][y] = nil
		}
	}
}
