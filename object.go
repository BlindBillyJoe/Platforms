package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type object struct {
	pos point
	col point

	_data []float32

	vao     uint32
	vbo     uint32
	program uint32

	speedX float32
	speedY float32
	speedZ float32
}

func (obj *object) move(x, y, z float32) {
	obj.pos.x += x
	obj.pos.y += y
	obj.pos.z += z
}

func (obj *object) update() {
	if obj.speedX != 0 || obj.speedY != 0 || obj.speedZ != 0 {
		obj.move(obj.speedX, obj.speedY, obj.speedZ)
	}
}

func loadObject(filename string) object {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	var vBuf []float32
	var result []float32

	scanner := bufio.NewScanner(f)

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
		_data: result,
	}
	obj.program = gl.CreateProgram()
	obj.makeVao()
	return obj
}

func (obj *object) attachShader(s *shader) {
	gl.AttachShader(obj.program, s.compiled)
	gl.LinkProgram(obj.program)
}
