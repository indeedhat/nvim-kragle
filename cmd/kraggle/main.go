package main

import (
	"errors"
	"fmt"

	"github.com/indeedhat/nvim-kraggle/internal/client"
	"github.com/indeedhat/nvim-kraggle/internal/config"
	"github.com/indeedhat/nvim-kraggle/internal/log"

	"github.com/neovim/go-client/nvim/plugin"
)

var pluginPtr *plugin.Plugin

func main() {
	defer log.Printf("closing plugin")
	defer log.Close()

	plugin.Main(func(p *plugin.Plugin) error {
		pluginPtr = p

		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleRemoteFocus"}, kragleRemoteFocus)
		p.HandleFunction(&plugin.FunctionOptions{Name: "KragleRemoteFocusBuffer"}, kragleRemoteFocusBuffer)
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
			log.Printf("Fatal error: %v", err)
		}

		return nil
	})
}

// kragleInit initializes the nvim client
func kragleInit(args []string) error {
	config.ReadFromClient(pluginPtr)
	log.Init(config.Get().LogPath)

	return nil
}

// kragleRemoteFocus will switch focus to the given client
func kragleRemoteFocus(args []string) error {
	if 1 < len(args) {
		return errors.New(fmt.Sprintf("No Server name Given %v", args))
	}

	return client.Focus(args[0])
}

// kragleRemoteFocusBuffer switch to the nvim instance owning the given file and bring it into focus
func kragleRemoteFocusBuffer(args []string) (string, error) {
	if 1 < len(args) {
		return fmt.Sprintf("No Path Given %v", args), errors.New(fmt.Sprintf("No Path Given %v", args))
	}

	err := client.FocusBuffer(args[0])
	if err != nil {
		return err.Error(), nil
	}

	return "Opened", nil
}

// kragleListFiles will display alist of buffers open in all connected instances
func kragleListFiles(includeSelf bool) ([]string, error) {
	var files []string

	if includeSelf {
		files = client.ListBuffers(pluginPtr)
	} else {
		files = client.ListBuffers()
	}

	log.Printf("passing files to client %v", files)

	return files, nil
}

// kragleAdoptBuffer from unknown peer
func kragleAdoptBuffer(args []string) error {
	if 1 < len(args) {
		return fmt.Errorf("No Path Given %v", args)
	}

	return client.AdoptBuffer(pluginPtr, args[0])
}

// kragleListServers connected to this one
func kragleListServers() ([]string, error) {
	log.Printf("gonna list peers")
	peers := client.ListPeers()

	var serverNames []string
	for name := range peers {
		serverNames = append(serverNames, name)
	}

	log.Printf("Servers: %v", serverNames)

	return serverNames, nil
}

// kragleCommandAll runs a vim command on all peers
//
// The command will be ran on the current instance last
func kragleCommandAll(args []string) error {
	if 1 < len(args) {
		return errors.New("No command sent")
	}

	for nvimPath, peer := range client.ListPeers() {
		kragleRemoteFocus([]string{nvimPath})
		peer.Command(args[0])
	}

	client.Focus(config.Get().ServerName)
	pluginPtr.Nvim.Command(args[0])

	return nil
}

// kragleOrphanBuffer will transfer the given buffer to the given remote client
func kragleOrphanBuffer(args []string) error {
	if 2 != len(args) {
		return errors.New("Invalid input")
	}

	log.Printf("orphan input %v", args)

	return client.OrphanBuffer(pluginPtr, args[0], args[1])
}
