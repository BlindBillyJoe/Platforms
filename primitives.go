package main

type point struct {
	x, y, z float32
}

type triangle struct {
	pt1, pt2, pt3 point
}

type model struct {
	data []triangle
}
