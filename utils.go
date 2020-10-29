package main

import (
	"math"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func closeAndNil(f **os.File) {
	(*f).Close()
	(*f) = nil
}

func cosInterpolate(start, end float64, count int) []float64 {
	var mod = math.Pi / float64(count)
	var sumMod = (end - start) / float64(count)
	var generator = func(v int) float64 {
		return math.Cos(float64(v)*mod - math.Pi)
	}
	var sum = 0.
	var result = make([]float64, count)
	for i := 0; i != count; i++ {
		var cos = (generator(i+1) + 1) * sumMod
		sum += cos
		result[i] = sum
	}
	return result
}

func roundTo(p point, prec float32) point {
	var e = 1 / prec
	p.x = float32(int(p.x*e)) / e
	p.y = float32(int(p.y*e)) / e
	p.z = float32(int(p.z*e)) / e
	return p
}
