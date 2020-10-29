package main

type point struct {
	x, y, z float32
}

type triangle struct {
	pt1, pt2, pt3 point
}

type square struct {
	tl, tr, bl, br point
}

type model struct {
	data []triangle
}

func (p point) cAdd(other point) point {
	return *p.add(other)
}

func (p *point) add(other point) *point {
	p.x += other.x
	p.y += other.y
	p.z += other.z
	return p
}

func (p *point) inc(v float32) *point {
	p.x *= v
	p.y *= v
	p.z *= v
	return p
}

func (p point) cInc(v float32) point {
	return *p.inc(v)
}

func sliceToPoints(points []float32) []point {
	var result = make([]point, len(points)/3)
	for i := range result {
		result[i] = point{
			x: points[i*3],
			y: points[i*3+1],
			z: points[i*3+2],
		}
	}
	return result
}

func buildBox(slice []float32) square {
	var points = sliceToPoints(slice)
	var result = square{
		tl: points[0],
		tr: points[0],
		bl: points[0],
		br: points[0],
	}
	for i := range points {
		if points[i].x >= result.tr.x && points[i].y >= result.tr.y {
			result.tr = points[i]
		}
		if points[i].x >= result.br.x && points[i].y <= result.br.y {
			result.br = points[i]
		}
		if points[i].x <= result.tl.x && points[i].y >= result.tl.y {
			result.tl = points[i]
		}
		if points[i].x <= result.bl.x && points[i].y <= result.bl.y {
			result.bl = points[i]
		}
	}
	return result
}

func (box *square) add(p point) *square {
	box.bl.add(p)
	box.br.add(p)
	box.tl.add(p)
	box.tr.add(p)
	return box
}

func (box square) cAdd(p point) square {
	return *box.add(p)
}

const (
	noIntersection = 0x0

	topIntersection    = 0x1
	bottomIntersection = 0x2
	leftIntersection   = 0x4
	rightIntersection  = 0x8

	vIntersection = rightIntersection | leftIntersection
	hIntersevtion = topIntersection | bottomIntersection
)

func (box *square) intersects(other *square) uint8 {
	var intersection uint8 = noIntersection
	if intersects(box, other) {
		var tt = triangle{
			pt1: other.tl,
			pt2: other.tr,
			pt3: other.center(),
		}
		var bt = triangle{
			pt1: other.center(),
			pt2: other.bl,
			pt3: other.br,
		}
		var lt = triangle{
			pt1: other.tl,
			pt2: other.center(),
			pt3: other.bl,
		}
		var rt = triangle{
			pt1: other.tr,
			pt2: other.center(),
			pt3: other.br,
		}

		if bt.contains(&box.tl) || bt.contains(&box.tr) {
			intersection |= topIntersection
		}

		if tt.contains(&box.bl) || tt.contains(&box.br) {
			intersection |= bottomIntersection
		}

		if lt.contains(&box.tr) || lt.contains(&box.br) {
			intersection |= rightIntersection
		}

		if rt.contains(&box.tl) || rt.contains(&box.bl) {
			intersection |= leftIntersection
		}
	}

	return intersection
}

func (box *square) center() point {
	return point{
		x: (box.br.x + box.tl.x) / 2,
		y: (box.br.y + box.tl.y) / 2,
	}
}

func (box *square) contains(p *point) bool {
	return p.x >= box.tl.x && p.x <= box.br.x && p.y >= box.br.y && p.y <= box.tl.y
}

func (p *point) equal(other *point) bool {
	return p.x == other.x && p.y == other.y && p.z == other.z
}

func (t *triangle) contains(p *point) bool {
	var f = func(a, b, p *point) float32 {
		return (p.x-a.x)*(b.y-a.y) - (p.y-a.y)*(b.x-a.x)
	}

	if p.equal(&t.pt1) || p.equal(&t.pt2) || p.equal(&t.pt3) {
		return true
	}

	var a = f(&t.pt1, &t.pt2, p) * f(&t.pt1, &t.pt2, &t.pt3)
	var b = f(&t.pt1, &t.pt3, p) * f(&t.pt1, &t.pt3, &t.pt2)
	var c = f(&t.pt2, &t.pt3, p) * f(&t.pt2, &t.pt3, &t.pt1)

	return a >= 0. && b >= 0. && c >= 0.
}

func intersects(o1, o2 *square) bool {
	return !(o1.tl.y < o2.br.y || o1.br.y > o2.tl.y || o1.tl.x > o2.br.x || o1.br.x < o2.tl.x)
}
