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

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleFunction(&plugin.FunctionOptions{Name: "RemoteOpen"}, func(args []string) (string, error) {
			if 1 < len(args) {
				return fmt.Sprintf("No Path Given %v", args), errors.New(fmt.Sprintf("No Path Given %v", args))
			}

			filePath := args[0]
			client, _ := findSwapOwner(filePath)
			if nil == client {
				return "Not found", nil
			}

			err := client.Command(fmt.Sprintf("drop %s", filePath))
			if nil != err {
				log(fmt.Sprintf("Error opening file %s", err))
			}

			return "opened", nil
		})

		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleInit"}, func(args []string) error {
			readConfigFromClient(p)
			initLog(config.LogPath)

			return nil
		})

		return nil
	})

	log("closing plugin")
	logClose()
}

func initLog(path string) {
	logPath = path
	logOpen()
	log(fmt.Sprintf("New Client %s", config.ServerName))
}

func buffOpen(client *nvim.Nvim, filePath string) (bool, *nvim.Buffer) {
	buffers, err := client.Buffers()
	if nil != err {
		return false, nil
	}

	for _, b := range buffers {
		// check for same root
		if config.SameRoot {
			var result string
			_ = client.Call("getcwd", &result)
			log(fmt.Sprintf("same check: %v - %v", result, config.ClientRoot))
			if 0 < len(result) && result != config.ClientRoot {
				continue
			}
		}

		name, err := client.BufferName(b)
		log(fmt.Sprintf("Checking %s against %s", name, filePath))

		if nil == err && name == filePath {
			log("path found")
			return true, &b
		}
	}

	return false, nil
}

func findSwapOwner(path string) (*nvim.Nvim, *nvim.Buffer) {
	connectAll()

	for _, client := range connections {
		if open, buffer := buffOpen(client, path); open {
			return client, buffer
		}
	}

	return nil, nil
}
