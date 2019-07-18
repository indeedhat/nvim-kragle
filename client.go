package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/neovim/go-client/nvim"
)

func connect(nvimPath string) (*nvim.Nvim, error) {
	client, err := nvim.Dial(nvimPath)
	if nil != err {
		log(fmt.Sprintf("Failed to connect to %s", nvimPath))
		return nil, err
	}

	connections[nvimPath] = client
	return client, err
}

func listUnconnectedPaths() []string {
	var paths []string

	files, err := ioutil.ReadDir(PATH_ROOT)
	if nil != err {
		return paths
	}

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), "nvim") || !f.IsDir() {
			continue
		}

		fpath := path.Join(PATH_ROOT, f.Name(), "0")
		if _, ok := connections[fpath]; ok {
			continue
		}

		if fpath == clientPath {
			continue
		}

		paths = append(paths, fpath)
	}

	log(fmt.Sprintf("Unconnected instances %v", paths))
	return paths
}

func connectAll() {
	for _, path := range listUnconnectedPaths() {
		connect(path)
	}
}
