package heimdall

import (
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
			// try new signature first: func(*context.Context) (*hp.Handlers, error)
			newPluginFunc, ok := newFunc.(func(*context.Context) (*hp.Handlers, error))
			if ok {
				plugins[stripExtension(file.Name())] = newPluginFunc
				continue
			}

			// make backward compatible with old signature: func(*context.Context) (hp.Handler, error)
			oldPluginFunc, ok := newFunc.(func(*context.Context) (hp.Handler, error))
			if ok {
				plugins[stripExtension(file.Name())] = func(ctx *context.Context) (*hp.Handlers, error) {
					handler, err := oldPluginFunc(ctx)
					if err != nil {
						return nil, err
					}
					return &hp.Handlers{
						Handler:        handler,
						CleanupHandler: nil,
					}, nil
				}
			}

		}
	}

	return plugins, nil

}

func stripExtension(name string) string {
	return name[:len(name)-pluginExtensionLength]
}
