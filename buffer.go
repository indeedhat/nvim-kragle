package main

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func findSwapOwner(path string) *nvim.Nvim {
	connectAll()

	for _, client := range connections {
		if open := bufferBelongsToClient(client, path); open {
			return client
		}
	}

	return nil
}

func bufferBelongsToClient(client *nvim.Nvim, filePath string) bool {
	for _, name := range bufferNames(client) {
		if name == filePath {
			log("path found")
			return true
		}
	}

	return false
}

func bufferNames(client *nvim.Nvim) []string {
	var files []string
	buffers, err := client.Buffers()
	if nil != err {
		return files
	}

	if config.SameRoot {
		var result string
		_ = client.Call("getcwd", &result)
		log(fmt.Sprintf("same check: %v - %v", result, config.ClientRoot))
		if 0 < len(result) && result != config.ClientRoot {
			return files
		}
	}

	for _, b := range buffers {
		name, err := client.BufferName(b)

		if nil != err || "" == name {
			continue
		}

		// i need to find a better way of ignoring these buffers
		if strings.HasSuffix(name, "/ControlP") || strings.HasSuffix(name, "NERDTree") {
			continue
		}

		files = append(files, name)
	}

	return files
}
