package heimdall

import (
	"fmt"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/patterninc/heimdall/pkg/context"
	hp "github.com/patterninc/heimdall/pkg/plugin"
)

const (
	pluginExtension       = `.so`
	newPluginFunction     = `New`
	pluginExtensionLength = len(pluginExtension)
)

func (h *Heimdall) loadPlugins() (map[string]func(*context.Context) (*hp.Handlers, error), error) {

	plugins := make(map[string]func(*context.Context) (*hp.Handlers, error))

	files, err := os.ReadDir(h.PluginsDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filename := path.Join(h.PluginsDirectory, file.Name()); !file.IsDir() && strings.HasSuffix(filename, pluginExtension) {
			p, err := plugin.Open(filename)
			if err != nil {
				return nil, err
			}
			newFunc, err := p.Lookup(newPluginFunction)
			if err != nil {
				return nil, err
			}
			// plugins must return *Handlers
			newPluginFunc, ok := newFunc.(func(*context.Context) (*hp.Handlers, error))
			if !ok {
				return nil, fmt.Errorf("plugin %s must return *plugin.Handlers", stripExtension(file.Name()))
			}
			plugins[stripExtension(file.Name())] = newPluginFunc
		}
	}

	return plugins, nil

}

func stripExtension(name string) string {
	return name[:len(name)-pluginExtensionLength]
}
