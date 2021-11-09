package client

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/indeedhat/nvim-kraggle/internal/config"
	"github.com/neovim/go-client/nvim"
)

var (
	connections = make(map[string]*nvim.Nvim)
	peers       = make(map[string]*nvim.Nvim)
	blacklist   = make(map[string]*nvim.Nvim)
)

// Connect the given nvimPath to kraggle
func Connect(nvimPath string) (*nvim.Nvim, error) {
	log.Printf("checking blacklist for: %s", nvimPath)

	if clientIsBlacklisted(nvimPath) {
		log.Printf("client is blacklisted")
		return nil, nil
	}

	log.Printf("Dialing %s", nvimPath)

	client, err := nvim.Dial(nvimPath)
	if nil != err {
		log.Printf("Failed to connect to %s", nvimPath)
		return nil, err
	}

	connections[nvimPath] = client

	log.Printf("peer checking %s", nvimPath)
	if clientIsPeer(client, nvimPath) {
		peers[nvimPath] = client
	}

	return client, err
}

// Focus a client by its server name
func Focus(serverName string) error {
	discoverUnconnectedClients()

	client, ok := peers[serverName]
	if !ok {
		return errors.New("Invalid client name")
	}

	return focusClient(client)
}

// ListPeers
func ListPeers() map[string]*nvim.Nvim {
	discoverUnconnectedClients()
	cleanupClosedPeers()

	log.Printf("peer list %v", peers)
	return peers
}

func listUnconnectedPaths() []string {
	var paths []string
	conf := config.Get()

	files, err := ioutil.ReadDir(conf.ServerRoot)
	if nil != err {
		return paths
	}

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "nvim") || !f.IsDir() {
			continue
		}

		fpath := path.Join(conf.ServerRoot, f.Name(), "0")
		if _, ok := connections[fpath]; ok {
			continue
		}

		if fpath == conf.ServerName {
			continue
		}

		paths = append(paths, fpath)
	}

	log.Printf("Unconnected instances %v", paths)
	return paths
}

func cleanupClosedPeers() {
	for path, _ := range peers {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delete(peers, path)
		}
	}
}

func clientIsPeer(client *nvim.Nvim, path string) bool {
	conf := config.Get()
	if conf.SameRoot {
		var result string

		log.Printf("requesting pwd")

		// TODO: this is not an ideal solution to the problem but the monkey patch will do for now
		//       this stops a lockup when scanning headless instances
		race(func() {
			_ = client.Call("getcwd", &result)
		}, 100*time.Millisecond)

		if 0 == len(result) {
			addToBlacklist(path)
			return false
		}

		log.Printf("same check: %v - %v", result, conf.ClientRoot)
		return 0 < len(result) && result == conf.ClientRoot
	}

	return true
}

func clientIsBlacklisted(path string) bool {
	return nil != blacklist[path]
}

func addToBlacklist(path string) {
	log.Printf("adding client %s to blacklist", path)
	blacklist[path] = peers[path]
	delete(peers, path)
}

func focusClient(client *nvim.Nvim) error {
	return client.Command("call kragle#focus()")
}

// discoverUnconnectedClients and connect them to kraggle
func discoverUnconnectedClients() {
	for _, path := range listUnconnectedPaths() {
		Connect(path)
	}
}
