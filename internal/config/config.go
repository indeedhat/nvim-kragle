package config

import (
	"os"
	"strings"

	"github.com/neovim/go-client/nvim/plugin"
)

type Config struct {
	ServerRoot  string
	ServerName  string
	LogPath     string
	SameRoot    bool
	ClientRoot  string
	UseTabs     bool
	OpenCommand string
}

var config Config

// Get the existing config instance
func Get() Config {
	return config
}

// ReadFromClient into a config object
func ReadFromClient(p *plugin.Plugin) error {
	config = Config{
		OpenCommand: "e",
	}

	res := make(map[string]interface{})
	err := p.Nvim.Call("kragle#getConfig", &res)

	if nil != err {
		return err
	}

	if val, ok := res["server_name"]; ok {
		config.ServerName = val.(string)

		parts := strings.Split(config.ServerName, string(os.PathSeparator))
		if 3 <= len(parts) {
			parts = parts[:len(parts)-2]
			config.ServerRoot = strings.Join(parts, string(os.PathSeparator))
		}
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

	if val, ok := res["use_tabs"]; ok {
		config.UseTabs = val.(bool)
		config.OpenCommand = "tabe"
	}

	return nil
}
