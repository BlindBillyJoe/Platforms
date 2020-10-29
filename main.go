package main

import (
	"math/rand"
	"time"
)

var gCounter = 0

func drawMap(mc *mapController) {
	counter := 0
	for i := range mc.sceneMap {
		for j := range mc.sceneMap[i] {
			if mc.sceneMap[j][i] != nil {
				counter++
			}
		}
	}

	if counter != gCounter && counter != 0 {
		for x := range mc.sceneMap {
			for y := range mc.sceneMap[x] {
				if mc.sceneMap[y][x] != nil {
					print("@")
				} else {
					print(" ")
				}
			}
			print("\n")
		}
		gCounter = counter
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var loader Loader
	loader.load("main.ini")
	loader.load("objects.ini")

	var window = loadWindowCfg(loader.xpath("Settings/Display"))
	var player = loadObjCfg(loader.xpath("Objects/Player"))

	ControllersManager.create("1Keyboard", 1)
	ControllersManager.create("3Physics", 2)
	ControllersManager.create("2Map", 3)

	var mapCtrl = ControllersManager.controllers["2Map"].(*mapController)
	var kbCtrl = ControllersManager.controllers["1Keyboard"].(*keyboardController)

	kbCtrl.activate(window)
	kbCtrl.add(&player)

	sf1 := loadObjCfg(loader.xpath("Objects/Surface1"))
	sf2 := loadObjCfg(loader.xpath("Objects/Surface2"))
	sf3 := loadObjCfg(loader.xpath("Objects/Surface3"))

	mapCtrl.add(&sf1)
	mapCtrl.add(&sf2)
	mapCtrl.add(&sf3)

	// var sr = strings.NewReader(createObjFile(generate(50, 50, 250)))
	// var sr = strings.NewReader(createObjFile(generate(20, 20, 100)))
	// generatedMapObj := createObjFile(generate(30, 30, 450))
	// println("map generated")
	// mapObjReader := strings.NewReader(generatedMapObj)
	// println("map readed")
	// mapObjs := loadMultiObjects(mapObjReader)
	// println("map parsed")
	// for i := range mapObjs {
	// 	mapObj := &mapObjs[i]
	// 	mapObj.attachShader(shaderDatabase.loaded["res/shaders/vShader.glsl"])
	// 	mapObj.attachShader(shaderDatabase.loaded["res/shaders/fShader.glsl"])
	// 	mapObj.col = point{
	// 		x: 0.5 + float32(i)/10,
	// 		y: 0.5,
	// 		z: 0.5,
	// 	}
	// 	mapCtrl.add(mapObj)
	// }

	println("map loaded")

	ControllersManager.controllers["3Physics"].add(&player)

	println("starts main loop...")

	for !window.shouldClose() {
		window.clear()
		window.pollEvents()
		ControllersManager.process()
		for i := range mapCtrl.units {
			draw(mapCtrl.units[i], window)
		}
		draw(&player, window)
		drawMap(mapCtrl)
		window.process()
	}
}
