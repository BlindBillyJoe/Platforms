package main

import (
	"fmt"
	"strings"
)

type face struct {
	v1, v2, v3 int
}

func (p point) String() string {
	return fmt.Sprintf("%.5f %.5f %.5f", p.x, p.y, p.z)
}

func (f face) String() string {
	return fmt.Sprintf("%v %v %v", f.v1, f.v2, f.v3)
}

func createObjFile(data [][]bool) string {
	var w, h = len(data), len(data[0])
	var builder strings.Builder
	var kx, ky = 2. / float32(w), 2. / float32(h)
	var vDatabase = make(map[int]point)
	var fDatabase = make([]face, 0)
	var contains = func(p point) int {
		for i := range vDatabase {
			if vDatabase[i].String() == p.String() {
				return i
			}
		}
		return -1
	}
	var counter int
	var add = func(p point) int {
		var ind = contains(p)
		if ind == -1 {
			vDatabase[counter] = p
			ind = counter
			counter++
		}
		return ind + 1
	}
	var xFlag = false
	var yFlag = false
	var tl point
	var tr point
	var bl point
	var br point

	var addFace = func() {
		fDatabase = append(fDatabase, face{v1: add(tl), v2: add(br), v3: add(bl)})
		fDatabase = append(fDatabase, face{v1: add(tl), v2: add(tr), v3: add(br)})
	}

	var yStart = func(x, y int) {
		tl.x = float32(x)*kx - 1.
		tl.y = 1. - float32(y)*ky
		tr.x = float32(x)*kx + kx - 1.
		tr.y = 1. - float32(y)*ky

		// println("yStart tl:", tl.x, tl.y, "tr:", tr.x, tr.y, x, y)

		yFlag = true
	}

	var yEnd = func(x, y int) {
		bl.x = float32(x)*kx - 1.
		bl.y = 1. - float32(y)*ky - ky
		br.x = float32(x)*kx + kx - 1.
		br.y = 1. - float32(y)*ky - ky

		// println("yEnd bl:", bl.x, bl.y, "br:", br.x, br.y, x, y)

		addFace()
		yFlag = false
	}

	var xStart = func(x, y int) {
		tl.x = float32(x)*kx - 1.
		tl.y = 1. - float32(y)*ky
		bl.x = float32(x)*kx - 1.
		bl.y = 1. - float32(y)*ky - ky

		// println("xStart tl:", tl.x, tl.y, "bl:", bl.x, bl.y, x, y)

		xFlag = true
	}

	var xEnd = func(x, y int) {
		tr.x = float32(x)*kx + kx - 1.
		tr.y = 1. - float32(y)*ky
		br.x = float32(x)*kx + kx - 1.
		br.y = 1. - float32(y)*ky - ky

		// println("xEnd tr:", tr.x, tr.y, "br:", br.x, br.y, x, y)

		addFace()
		xFlag = false
	}

	type _callback func(int, int)

	var check = func(x, y int, l *int, _flag *bool, start _callback, end _callback, isX bool) {
		if data[x][y] {
			if !*_flag {
				start(x, y)
			}
			*l++
		} else {
			if *_flag && *l > 1 {
				if isX {
					end(x-1, y)
				} else {
					end(x, y-1)
				}
			}
			*_flag = false
			*l = 0
		}
	}

	var checkY = func(x, y int, l *int) {
		check(x, y, l, &yFlag, yStart, yEnd, false)
	}

	var checkX = func(x, y int, l *int) {
		check(x, y, l, &xFlag, xStart, xEnd, true)
	}

	for x := range data {
		var l = 0
		for y := range data[x] {
			checkY(x, y, &l)
		}
		if yFlag && l > 1 {
			// println("line")
			yEnd(x, len(data[0])-1)
		}
		l = 0
		yFlag = false
	}

	for y := range data[0] {
		var l = 0
		for x := range data {
			checkX(x, y, &l)
		}
		if xFlag && l > 1 {
			// println("line")
			xEnd(len(data)-1, y)
		}
		l = 0
		xFlag = false
	}

	for i := 0; i != len(vDatabase); i++ {
		builder.WriteString(fmt.Sprintf("v %v\n", vDatabase[i].String()))
	}
	for i := range fDatabase {
		builder.WriteString(fmt.Sprintf("f %v\n", fDatabase[i].String()))
	}
	println("#v:", len(vDatabase), "#f:", len(fDatabase))
	// println(builder.String())
	return builder.String()
}
