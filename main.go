package main

func main() {
	var loader Loader
	loader.load("main.ini")
	loader.load("objects.ini")

	var window = loadWindowCfg(loader.xpath("Settings/Display"))
	var player = loadObjCfg(loader.xpath("Objects/Player"))
	var surface = loadObjCfg(loader.xpath("Objects/Surface"))

	MainKeyboardController.units = append(MainKeyboardController.units, &player)
	MainKeyboardController.activate(window)

	for !window.shouldClose() {
		window.clear()
		player.update()
		draw(&surface, window)
		draw(&player, window)
		window.process()
	}
}
