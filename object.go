package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type object struct {
	pos point
	col point

	_data    []float32
	colision square

	vao     uint32
	vbo     uint32
	program uint32

	speed point

	xAcceleration  float32
	nxAcceleration float32

	yAcceleration  float32
	nyAcceleration float32

	intersection uint8

	updated bool
}

func (obj *object) move(pos point) {
	obj.pos.add(pos)

	if obj.pos.x > 1. {
		obj.pos.x = -1.
	}

	if obj.pos.x < -1. {
		obj.pos.x = 1.
	}

	if obj.pos.y > 1. {
		obj.pos.y = -1.
	}

	if obj.pos.y < -1. {
		obj.pos.y = 1.
	}
}

func (obj *object) update() {
	if obj.speed.x < speed && obj.xAcceleration > 0 {
		obj.speed.x += obj.xAcceleration
	} else if obj.speed.x > 0 && obj.xAcceleration == 0 {
		obj.speed.x -= (speed + obj.speed.x) / acceleration
		if obj.speed.x < 0 {
			obj.speed.x = 0
		}
	}

	if obj.speed.x > -speed && obj.nxAcceleration > 0 {
		obj.speed.x -= obj.nxAcceleration
	} else if obj.speed.x < 0 && obj.nxAcceleration == 0 {
		obj.speed.x += (speed - obj.speed.x) / acceleration
		if obj.speed.x > 0 {
			obj.speed.x = 0
		}
	}

	if obj.speed.y < speed && obj.yAcceleration > 0 {
		obj.speed.y += obj.yAcceleration
	} else if obj.speed.y > 0 && obj.yAcceleration == 0 {
		obj.speed.y -= (speed + obj.speed.y) / acceleration
		if obj.speed.y < 0 {
			obj.speed.y = 0
		}
	}

	if obj.speed.y > -speed && obj.nyAcceleration > 0 {
		obj.speed.y -= obj.nyAcceleration
	} else if obj.speed.y < 0 && obj.nyAcceleration == 0 {
		obj.speed.y += (speed - obj.speed.y) / acceleration
		if obj.speed.y > 0 {
			obj.speed.y = 0
		}
	}

	obj.move(obj.speed)

	obj.updated = false
}

type intersectionCallback func(*object, *object)

func intCallback(obj, other *object) {
	var h = (obj.colision.tl.y - obj.colision.br.y) / 2.
	var nw = (obj.colision.tl.x - obj.colision.br.x) / 2.
	var pw = (obj.colision.br.x - obj.colision.tl.x) / 2.
	if obj.intersection&topIntersection != 0 {
		obj.pos.y = other.colision.br.cAdd(other.pos).y - h
		if obj.speed.y > 0. {
			obj.speed.y = 0.
		}
	}
	if obj.intersection&bottomIntersection != 0 {
		obj.pos.y = other.colision.tl.cAdd(other.pos).y + h
		if obj.speed.y < 0. {
			obj.speed.y = 0.
		}
	}
	if obj.intersection&leftIntersection != 0 {
		obj.pos.x = other.colision.br.cAdd(other.pos).x + pw
		if obj.speed.x < 0. {
			obj.speed.x = 0.
		}
	}
	if obj.intersection&rightIntersection != 0 {
		obj.pos.x = other.colision.tl.cAdd(other.pos).x + nw
		if obj.speed.x > 0. {
			obj.speed.x = 0.
		}
	}
}

func (obj *object) updatePhysics() {
	var gravity float32 = 0.005

	if obj.speed.y > -speed {
		obj.yAcceleration -= speed / acceleration
		obj.speed.y -= gravity
	}

	checkIntersections(obj, intCallback)
}

func loadMultiObjects(r io.Reader) []object {
	var objs []object
	var vBuf []float32
	var result []float32

	flip := true

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var lStr = strings.Split(scanner.Text(), " ")
		if lStr[0] == "v" {
			for i := 1; i != len(lStr); i++ {
				var val, err = strconv.ParseFloat(lStr[i], 32)
				check(err)
				vBuf = append(vBuf, float32(val))
			}
		} else if lStr[0] == "f" {
			if flip {
				result = []float32{}
			}
			for i := 1; i != len(lStr); i++ {
				var strVal = lStr[i]
				if strings.Contains(strVal, "/") {
					strVal = strings.Split(strVal, "/")[0]
				}
				var val, err = strconv.Atoi(strVal)
				check(err)
				result = append(result, vBuf[(val-1)*3])
				result = append(result, vBuf[(val-1)*3+1])
				result = append(result, vBuf[(val-1)*3+2])
			}
			if !flip {
				var obj = object{
					_data:    result,
					colision: buildBox(result),
					updated:  false,
				}
				obj.program = gl.CreateProgram()
				obj.makeVao()
				objs = append(objs, obj)
			}
			flip = !flip
		}
	}
	return objs
}

func loadObject(r io.Reader) object {
	var vBuf []float32
	var result []float32

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var lStr = strings.Split(scanner.Text(), " ")
		if lStr[0] == "v" {
			for i := 1; i != len(lStr); i++ {
				var val, err = strconv.ParseFloat(lStr[i], 32)
				check(err)
				vBuf = append(vBuf, float32(val))
			}
		} else if lStr[0] == "f" {
			for i := 1; i != len(lStr); i++ {
				var strVal = lStr[i]
				if strings.Contains(strVal, "/") {
					strVal = strings.Split(strVal, "/")[0]
				}
				var val, err = strconv.Atoi(strVal)
				check(err)
				result = append(result, vBuf[(val-1)*3])
				result = append(result, vBuf[(val-1)*3+1])
				result = append(result, vBuf[(val-1)*3+2])
			}
		}
	}
	var obj = object{
		_data:    result,
		colision: buildBox(result),
		updated:  false,
	}
	obj.program = gl.CreateProgram()
	obj.makeVao()
	return obj
}

func (obj *object) attachShader(s *shader) {
	gl.AttachShader(obj.program, s.compiled)
	gl.LinkProgram(obj.program)
}

func (obj *object) intersects(other *object) uint8 {
	var box1 = obj.colision
	var box2 = other.colision

	return box1.add(obj.pos).intersects(box2.add(other.pos))
}

func checkIntersections(obj *object, callback intersectionCallback) uint8 {
	var mc = ControllersManager.controllers["2Map"].(*mapController)
	var box = obj.colision.cAdd(obj.pos)
	obj.intersection = noIntersection
	for x := box.tl.x; x <= box.br.x; x += 0.01 {
		for y := box.br.y; y <= box.tl.y; y += 0.01 {
			var target *object
			{
				pX := int((x+1.)*hMapSize/2 + 0.5)
				pY := int((y+1.)*vMapSize/2 + 0.5)
				if pX >= hMapSize {
					pX = hMapSize - 1
				} else if pX < 0 {
					pX = 0
				}
				if pY >= vMapSize {
					pY = vMapSize - 1
				} else if pY < 0 {
					pY = 0
				}
				target = mc.sceneMap[pX][pY]
			}

			if target != nil {
				bInter := obj.intersects(target)
				obj.intersection |= bInter
				if bInter != 0 {
					callback(obj, target)
				}
			}
		}
	}
	return obj.intersection
}
