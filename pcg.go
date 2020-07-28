package main

import (
	"math/rand"
)

var sqrc int

var counter int

type cell struct {
	walls [][]*bool
}

func boundCell(c cell) {
	for x := range c.walls {
		for y := range c.walls[x] {
			if x == 0 || x == len(c.walls)-1 {
				*c.walls[x][y] = true
			} else if y == 0 || y == len(c.walls[x])-1 {
				*c.walls[x][y] = true
			}
		}
	}
}

func drawCells(c *[][]bool) {
	for y := range (*c)[0] {
		for x := range *c {
			if (*c)[x][y] {
				print("@")
			} else {
				print(" ")
			}
		}
		print("\n")
	}
}

func split(c cell, cells []cell, ind int) []cell {
	if counter == sqrc {
		return cells
	}
	var x, y int
	var xy int
	if len(c.walls) > len(c.walls[0]) {
		xy = 1
	}
	if xy == 0 {
		var l = len(c.walls[0])
		var d int
		if l > 6 {
			d = l / 3
		} else {
			d = l / 2
		}
		y = l/2 - d/2 + rand.Int()%(d+1)
	} else {
		var l = len(c.walls)
		var d int
		if l > 6 {
			d = l / 3
		} else {
			d = l / 2
		}
		x = l/2 - d/2 + rand.Int()%(d+1)
	}

	var create = func(w, h int) cell {
		var cl cell
		cl.walls = make([][]*bool, w)
		for i := range cl.walls {
			cl.walls[i] = make([]*bool, h)
		}
		return cl
	}
	var bound = func(cl cell, sx, sy int) {
		for x := range cl.walls {
			for y := range cl.walls[0] {
				cl.walls[x][y] = c.walls[x+sx][y+sy]
			}
		}
	}
	var c1, c2 cell
	if xy == 0 {
		c1 = create(len(c.walls), len(c.walls[0])-y)
		c2 = create(len(c.walls), y+1)
		bound(c1, 0, 0)
		bound(c2, 0, len(c1.walls[0])-1)
	} else {
		c1 = create(len(c.walls)-x, len(c.walls[0]))
		c2 = create(x+1, len(c.walls[0]))
		bound(c1, 0, 0)
		bound(c2, len(c1.walls)-1, 0)
	}
	boundCell(c1)
	boundCell(c2)
	counter++
	cells[ind] = c1
	cells = append(cells, c2)
	var index int
	var sqr int
	for i := range cells {
		if sqr < len(cells[i].walls)*len(cells[i].walls[0]) {
			index = i
			sqr = len(cells[i].walls) * len(cells[i].walls[0])
		}
	}
	return split(cells[index], cells, index)
}

func makeDoors(cells []cell) {
	var checkDoors = func(cl cell) (bool, bool, bool, bool) {
		var t bool
		var b bool
		var l bool
		var r bool
		for x := range cl.walls {
			if !*cl.walls[x][0] {
				t = true
			}
			if !*cl.walls[x][len(cl.walls[x])-1] {
				b = true
			}
		}
		for y := range cl.walls[0] {
			if !*cl.walls[0][y] {
				l = true
			}
			if !*cl.walls[len(cl.walls)-1][y] {
				r = true
			}
		}
		return t, r, b, l
	}
	for i := range cells {
		var cl = cells[i]
		var t, r, b, l = checkDoors(cl)
		if !t {
			*cl.walls[len(cl.walls)/2][0] = false
		}
		if !b {
			*cl.walls[len(cl.walls)/2][len(cl.walls[0])-1] = false
		}
		if !l {
			*cl.walls[0][len(cl.walls[0])/2] = false
		}
		if !r {
			*cl.walls[len(cl.walls)-1][len(cl.walls[0])/2] = false
		}
	}
}

func generate(w, h, k int) [][]bool {
	counter = 1
	sqrc = w * h / k
	println(sqrc)
	var wcells = make([][]bool, w)
	for i := range wcells {
		wcells[i] = make([]bool, h)
	}
	var cells = make([]cell, 0, 1)
	var mCell cell
	mCell.walls = make([][]*bool, w)
	for i := range mCell.walls {
		mCell.walls[i] = make([]*bool, h)
		for j := range mCell.walls[i] {
			mCell.walls[i][j] = &wcells[i][j]
		}
	}
	cells = append(cells, mCell)
	boundCell(mCell)
	cells = split(mCell, cells, 0)
	makeDoors(cells)
	boundCell(mCell)
	drawCells(&wcells)
	return wcells
}
