package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/neovim/go-client/nvim"
)

var (
	connections = make(map[string]*nvim.Nvim)
	peers       = make(map[string]*nvim.Nvim)
	blacklist   = make(map[string]*nvim.Nvim)
)

func connect(nvimPath string) (*nvim.Nvim, error) {
	log("checking blacklist for: %s", nvimPath)

	if clientIsBlacklisted(nvimPath) {
		log("client is blacklisted")
		return nil, nil
	}

	log("Dialing %s", nvimPath)

	client, err := nvim.Dial(nvimPath)
	if nil != err {
		log("Failed to connect to %s", nvimPath)
		return nil, err
	}

	connections[nvimPath] = client

	log("peer checking %s", nvimPath)
	if clientIsPeer(client, nvimPath) {
		peers[nvimPath] = client
	}

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

func cleanupClosedPeers() {
	for path, _ := range peers {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delete(peers, path)
		}
	}
}

func clientIsPeer(client *nvim.Nvim, path string) bool {
	if config.SameRoot {
		var result string

		log("requesting pwd")

		// TODO: this is not an ideal solution to the problem but the monkey patch will do for now
		//       this stops a lockup when scanning headless instances
		setTimeout(func() {
			_ = client.Call("getcwd", &result)
		}, 100*time.Millisecond)

		if 0 == len(result) {
			addToBlacklist(path)
			return false
		}

		log("same check: %v - %v", result, config.ClientRoot)
		return 0 < len(result) && result == config.ClientRoot
	}

	return true
}

func clientIsBlacklisted(path string) bool {
	return nil != blacklist[path]
}

func addToBlacklist(path string) {
	log("adding client %s to blacklist", path)
	blacklist[path] = peers[path]
	delete(peers, path)
}

func listPeers() map[string]*nvim.Nvim {
	connectAll()
	cleanupClosedPeers()

	log("peer list %v", peers)
	return peers
}

func focusClient(client *nvim.Nvim) error {
	return client.Command("call kragle#focus()")
}
