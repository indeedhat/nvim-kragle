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
