package main

import (
	"errors"
	"fmt"

	"github.com/neovim/go-client/nvim/plugin"
)

var pluginPtr *plugin.Plugin

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		pluginPtr = p

		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleRemoteFocus"}, kragleRemoteFocus)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleRemoteFocusBuffer"}, kragleRemoteFocusBuffer)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleRemoteOpen"}, kragleRemoteOpen)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleInit"}, kragleInit)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleListAllFiles"}, func() ([]string, error) {
			return kragleListFiles(true)
		})
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleListRemoteFiles"}, func() ([]string, error) {
			return kragleListFiles(false)
		})
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleAdoptBuffer"}, kragleAdoptBuffer)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleListServers"}, kragleListServers)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleOrphanBuffer"}, kragleOrphanBuffer)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleCommandAll"}, kragleCommandAll)

		if err := recover(); nil != err {
			log("Fatal error: %v", err)
		}

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

func kragleRemoteFocus(args []string) error {
	if 1 < len(args) {
		return errors.New(fmt.Sprintf("No Server name Given %v", args))
	}

	serverName := args[0]

	connectAll()
	client, ok := peers[serverName]
	if !ok {
		return errors.New("Invalid client name")
	}

	return focusClient(client)
}

func kragleRemoteFocusBuffer(args []string) (string, error) {
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
		log("Error opening file %s", err)
	}
	err = client.Command("call kragle#focus()")
	log("calling foreground %v", err)

	return "opened", nil
}

func kragleListFiles(includeSelf bool) ([]string, error) {
	connectAll()

	var files []string
	if includeSelf {
		files = bufferNames(pluginPtr.Nvim)
	}

	for _, client := range peers {
		if bufferFiles := bufferNames(client); 0 < len(bufferFiles) {
			files = append(files, bufferFiles...)
		}
	}

	log("passing files to client %v", files)
	return files, nil
}

func kragleAdoptBuffer(args []string) error {
	if 1 < len(args) {
		return fmt.Errorf("No Path Given %v", args)
	}

	client := findSwapOwner(args[0])

	if nil == client {
		log("failed to find client for adoption")
		return errors.New("client not found")
	}

	buffer := bufferFromName(client, args[0])
	if buffer == nil {
		log("failed to find file for adoption")
		return errors.New("could not find buffer")
	}

	return moveBufferToClient(buffer, args[0], client, pluginPtr.Nvim)
}

func kragleListServers() ([]string, error) {
	log("gonna list peers")
	peers := listPeers()
	log("listed peers")
	serverNames := make([]string, 0, len(peers))

	for name := range peers {
		serverNames = append(serverNames, name)
	}

	log("Servers: %v", serverNames)
	return serverNames, nil
}

func kragleOrphanBuffer(args []string) error {
	if 2 != len(args) {
		return errors.New("Invalid input")
	}

	log("orphan input %v", args)
	bufferName := args[0]
	clientName := args[1]

	connectAll()
	client, ok := connections[clientName]
	if !ok {
		return errors.New("Invalid client name")
	}

	buffer := bufferFromName(pluginPtr.Nvim, bufferName)
	if nil == buffer {
		return errors.New("Invalid buffer name")
	}
	log("moving buffer %s", bufferName)

	return moveBufferToClient(buffer, bufferName, pluginPtr.Nvim, client)
}

func kragleCommandAll(args []string) error {
	if 1 < len(args) {
		return errors.New("No command sent")
	}

	for nvimPath, client := range listPeers() {
		kragleRemoteFocus([]string{nvimPath})
		client.Command(args[0])
	}

	focusClient(pluginPtr.Nvim)
	pluginPtr.Nvim.Command(args[0])
	return nil
}

func kragleRemoteOpen(args []string) error {
	if 2 < len(args) {
		return errors.New("Invalid input")
	}

	return openBufferOnClient(args[0], args[1])
}
