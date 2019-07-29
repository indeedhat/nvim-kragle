package main

import (
	"fmt"
	"github.com/neovim/go-client/nvim/plugin"
	"io/ioutil"
)

type kragleConfig struct {
	ServerName string
	LogPath    string
	SameRoot   bool
	ClientRoot string
}

var config kragleConfig

func readConfigFromClient(p *plugin.Plugin) error {
	config = kragleConfig{}

	res := make(map[string]interface{})
	err := p.Nvim.Call("kragle#getConfig", &res)

	if nil != err {
		return err
	}

	if val, ok := res["server_name"]; ok {
		config.ServerName = val.(string)
	}

	if val, ok := res["log_path"]; ok {
		config.LogPath = val.(string)
	}

	if val, ok := res["same_root"]; ok {
		config.SameRoot = val.(bool)
	}

	if val, ok := res["client_root"]; ok {
		config.ClientRoot = val.(string)
	}

	ioutil.WriteFile(
		"/home/indeedhat/Documents/krag_debug",
		[]byte(fmt.Sprintf("%v - %v - %v", err, res, config)),
		0644,
	)

	return nil
}
