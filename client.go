package main

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func connect(nvimPath string) (*nvim.Nvim, error) {
	client, err := nvim.Dial(nvimPath)
	if nil != err {
		log("Failed to connect to %s", nvimPath)
		return nil, err
	}

	connections[nvimPath] = client
	return client, err
}

func listUnconnectedPaths() []string {
	var paths []string

	files, err := ioutil.ReadDir(config.ServerRoot)
	if nil != err {
		return paths
	}

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "nvim") || !f.IsDir() {
			continue
		}

		fpath := path.Join(config.ServerRoot, f.Name(), "0")
		if _, ok := connections[fpath]; ok {
			continue
		}

		if fpath == config.ServerName {
			continue
		}

		paths = append(paths, fpath)
	}

	log("Unconnected instances %v", paths)
	return paths
}

func connectAll() {
	for _, path := range listUnconnectedPaths() {
		connect(path)
	}
}

func clientIsPeer(client *nvim.Nvim) bool {
	if config.SameRoot {
		var result string
		_ = client.Call("getcwd", &result)
		log("same check: %v - %v", result, config.ClientRoot)
		return 0 < len(result) && result == config.ClientRoot
	}

	return true
}

func listPeers() map[string]*nvim.Nvim {
	connectAll()

	peers := make(map[string]*nvim.Nvim)

	for serverName, client := range connections {
		if !clientIsPeer(client) {
			continue
		}

		peers[serverName] = client
	}

	return peers
}
