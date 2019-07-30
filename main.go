package main

import (
	"errors"
	"fmt"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

const (
	PATH_ROOT = "/tmp"
)

var connections = make(map[string]*nvim.Nvim)
var pluginPtr *plugin.Plugin

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		pluginPtr = p

		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleRemoteOpen"}, kragleRemoteOpen)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleInit"}, kragleInit)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleListFiles"}, kragleListFiles)

		return nil
	})

	log("closing plugin")
	logClose()
}

func kragleInit(args []string) error {
	readConfigFromClient(pluginPtr)
	initLog(config.LogPath)

	return nil
}

func kragleRemoteOpen(args []string) (string, error) {
	if 1 < len(args) {
		return fmt.Sprintf("No Path Given %v", args), errors.New(fmt.Sprintf("No Path Given %v", args))
	}

	filePath := args[0]
	client := findSwapOwner(filePath)
	if nil == client {
		return "Not found", nil
	}

	err := client.Command(fmt.Sprintf("drop %s", filePath))
	if nil != err {
		log(fmt.Sprintf("Error opening file %s", err))
	}
	err = client.Command("call foreground()")
	log(fmt.Sprintf("calling foreground %v", err))

	err = client.Call("foreground", nil)
	log(fmt.Sprintf("calling foreground 2 %v", err))

	return "opened", nil
}

func kragleListFiles() ([]string, error) {
	connectAll()

	files := bufferNames(pluginPtr.Nvim)

	for _, client := range connections {
		if bufferFiles := bufferNames(client); 0 < len(bufferFiles) {
			files = append(files, bufferFiles...)
		}
	}

	log(fmt.Sprintf("passing files to client %v", files))
	return files, nil
}
