package client

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/indeedhat/nvim-kraggle/internal/config"
	"github.com/indeedhat/nvim-kraggle/internal/log"
	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

// FindSwapOwner for the given file path
func FindSwapOwner(path string) *nvim.Nvim {
	discoverUnconnectedClients()

	for _, client := range peers {
		if open := bufferBelongsToClient(client, path); open {
			return client
		}
	}

	return nil
}

// FocusBuffer by filepath in whatever nvim instance it belongs to
func FocusBuffer(filePath string) error {
	client := FindSwapOwner(filePath)
	if nil == client {
		return errors.New("Not found")
	}

	err := client.Command(fmt.Sprintf("drop %s", filePath))
	if nil != err {
		return errors.New("Failed to focus buffer")
	}

	err = client.Command("call kragle#focus()")
	if err != nil {
		return errors.New("Failed to focus buffer")
	}

	return nil
}

// ListBuffers on connected clients
//
// if the self client is passed then it will also be scanned
func ListBuffers(self ...*plugin.Plugin) []string {
	discoverUnconnectedClients()

	var files []string

	if len(self) == 1 {
		files = bufferNames(self[0].Nvim)
	}

	for _, client := range peers {
		clientBuffers := bufferNames(client)

		if len(clientBuffers) > 0 {
			files = append(files, clientBuffers...)
		}
	}

	return files
}

// AdoptBuffer from whatever client it belongs to by name
func AdoptBuffer(self *plugin.Plugin, bufferName string) error {
	peer := FindSwapOwner(bufferName)
	if peer == nil {
		log.Printf("failed to find client for adoption")
		return errors.New("client not found")
	}

	buffer := bufferFromName(peer, bufferName)
	if buffer == nil {
		log.Printf("failed to find file for adoption")
		return errors.New("could not find buffer")
	}

	return moveBufferToClient(buffer, bufferName, peer, self.Nvim)
}

// OrphanBuffer to the given remote client
func OrphanBuffer(self *plugin.Plugin, bufferName, clientName string) error {
	discoverUnconnectedClients()
	client, ok := connections[clientName]
	if !ok {
		return errors.New("Invalid client name")
	}

	buffer := bufferFromName(self.Nvim, bufferName)
	if nil == buffer {
		return errors.New("Invalid buffer name")
	}

	log.Printf("moving buffer %s", bufferName)

	return moveBufferToClient(buffer, bufferName, self.Nvim, client)
}

func bufferBelongsToClient(client *nvim.Nvim, filePath string) bool {
	for _, name := range bufferNames(client) {
		if name == filePath {
			log.Printf("path found")
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

	for _, b := range buffers {
		name, err := client.BufferName(b)

		if nil != err || "" == name {
			continue
		}

		// i need to find a better way of ignoring these buffers
		// TODO: find out if this is necesarry not or if the IsBufferLoaded thing sorts it
		//       i suspect it might cover ctrlP but not nerdtree
		_, fileName := path.Split(name)
		if "CtonrolP" == fileName || strings.HasPrefix(fileName, "NERD_tree") {
			continue
		}

		if is, _ := client.IsBufferLoaded(b); !is {
			continue
		}

		files = append(files, name)
	}

	return files
}

func bufferFromName(client *nvim.Nvim, bufferName string) *nvim.Buffer {
	buffers, err := client.Buffers()
	if nil != err {
		return nil
	}

	for _, b := range buffers {
		name, _ := client.BufferName(b)
		if bufferName == name {
			return &b
		}
	}

	return nil
}

func moveBufferToClient(buf *nvim.Buffer, bufName string, from, to *nvim.Nvim) error {
	log.Printf("detaching file %s from parent", bufName)
	err := from.Command(fmt.Sprintf("bd %d", int(*buf)))
	if nil != err {
		return err
	}

	err = to.Command(fmt.Sprintf("%s %s", config.Get().OpenCommand, bufName))
	log.Printf("opening in new parent %v", err)
	if nil != err {
		return err
	}

	err = to.Command("call kragle#focus()")
	log.Printf("calling foreground %v", err)

	return err
}

func openBufferOnClient(bufferName, clientName string) error {
	discoverUnconnectedClients()

	client := peers[clientName]
	if nil == client {
		return errors.New("Bad client name")
	}

	err := client.Command(fmt.Sprintf("%s %s", config.Get().OpenCommand, bufferName))
	if nil != err {
		return err
	}

	err = client.Command("call kragle#focus()")
	log.Printf("calling foreground %v", err)

	return err
}
