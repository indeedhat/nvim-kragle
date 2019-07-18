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

var clientPath string
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
				return "Not found", errors.New("Not found")
			}

			err := client.Command(fmt.Sprintf("drop %s", filePath))
			if nil != err {
				log(fmt.Sprintf("Error opening file %s", err))
			}

			return "opened", nil
		})

		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleInit"}, func(args []string) error {
			if 1 < len(args) {
				return errors.New("no server name")
			}

			clientPath = args[0]
			return nil
		})

		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleLog"}, func(args []string) {
			if 1 <= len(args) {
				logPath = args[0]
				logOpen()
				log(fmt.Sprintf("New Client %s", clientPath))
			}
		})

		return nil
	})

	log("closing plugin")
	logClose()
}

func buffOpen(client *nvim.Nvim, filePath string) (bool, *nvim.Buffer) {
	buffers, err := client.Buffers()
	if nil != err {
		return false, nil
	}

	for _, b := range buffers {
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
