package main

import (
	"fmt"
	"path"
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
	if !clientIsPeer(client) {
		return files
	}

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
	log("detaching file %s from parent", bufName)
	err := from.Command(fmt.Sprintf("bd %d", int(*buf)))
	if nil != err {
		return err
	}

	err = to.Command(fmt.Sprintf("tabe %s", bufName))
	log("opening in new parent %v", err)

	return err
}
