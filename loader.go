package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	headerRe, _    = regexp.Compile("^\\[(\\w+)::(\\w+)]$")
	dataRe, _      = regexp.Compile("^\\s*(\\S+)\\s*=\\s*(\\S+)\\s*$")
	shaderDatabase = database{
		loaded: make(map[string]*shader),
	}
)

//DataNode is data node struct
type DataNode struct {
	data     map[string]string
	children map[string]*DataNode
}

type database struct {
	loaded map[string]*shader
}

func (l *Loader) xpath(path string) *DataNode {
	var p = strings.Split(path, "/")
	var node = l.root
	for i := range p {
		var ptr = node.children[p[i]]
		if ptr != nil {
			node = ptr
		}
	}
	return node
}

//Loader is loader for data declaration files
type Loader struct {
	root        *DataNode
	currentNode *DataNode
}

func createDataNode() *DataNode {
	return &DataNode{
		data:     make(map[string]string),
		children: make(map[string]*DataNode),
	}
}

func (l *Loader) load(filename string) {
	if l.root == nil {
		l.root = createDataNode()
	}

	var f, err = os.Open(filename)
	check(err)
	defer f.Close()

	var scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		l.parse(scanner.Bytes())
	}
}

func (l *Loader) parse(line []byte) {
	if headerRe.Match(line) {
		var all = headerRe.FindSubmatch(line)
		l.currentNode = l.root
		l.startGroup(all[1])
		l.startGroup(all[2])
	} else if dataRe.Match(line) {
		var all = dataRe.FindSubmatch(line)
		l.currentNode.data[string(all[1])] = string(all[2])
	}
}

func (l *Loader) startGroup(data []byte) {
	var nameStr = string(data)
	if l.currentNode.children[nameStr] == nil {
		l.currentNode.children[nameStr] = createDataNode()
	}
	l.currentNode = l.currentNode.children[nameStr]
}

func loadObjCfg(config *DataNode) object {
	var obj = loadObject(config.data["Model"])

	var checkShaderInDb = func(sPath string, sType uint32) {
		if shaderDatabase.loaded[sPath] == nil {
			shaderDatabase.loaded[sPath] = createShaderFromFile(sPath, sType)
			check(shaderDatabase.loaded[sPath].compile())
		}
	}

	var vShaderPath = config.data["VShader"]
	checkShaderInDb(vShaderPath, gl.VERTEX_SHADER)
	obj.attachShader(shaderDatabase.loaded[vShaderPath])

	var fShaderPath = config.data["FShader"]
	if len(fShaderPath) != 0 {
		checkShaderInDb(fShaderPath, gl.FRAGMENT_SHADER)
		obj.attachShader(shaderDatabase.loaded[fShaderPath])
	}

	if len(config.data["Color"]) != 0 {
		var color = strings.Split(config.data["Color"], "/")
		var col, _ = strconv.ParseFloat(color[0], 32)
		obj.col.x = float32(col)
		col, _ = strconv.ParseFloat(color[1], 32)
		obj.col.y = float32(col)
		col, _ = strconv.ParseFloat(color[2], 32)
		obj.col.z = float32(col)
	}
	return obj
}

func loadWindowCfg(config *DataNode) *Window {
	var w, _ = strconv.Atoi(config.data["Width"])
	var h, _ = strconv.Atoi(config.data["Height"])
	if len(config.data["Title"]) != 0 {
		return initWindow(w, h, config.data["Title"])
	}
	return initWindow(w, h, "Main")
}
