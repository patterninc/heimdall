package auth

import (
	"fmt"
	"net/http"
	"plugin"

	hp "github.com/patterninc/heimdall/pkg/plugin"
	"gopkg.in/yaml.v3"
)

const (
	pluginKey         = `plugin`
	contextKey        = `context`
	newPluginFunction = `New`
)

var (
	ErrNoPluginKey       = fmt.Errorf(`no 'plugin' key in auth`)
	ErrPluginKeyNoString = fmt.Errorf(`'plugin' key in auth is not a string`)
	ErrNoAuthPlugin      = fmt.Errorf(`unexpected plugin in auth`)
)

type Auth struct {
	Plugin hp.Auth
}

type hasInit interface {
	Init() error
}

func (a *Auth) UnmarshalYAML(unmarshal func(interface{}) error) error {

	value := make(map[string]any)

	if err := unmarshal(&value); err != nil {
		return err
	}

	pluginNameValue, found := value[pluginKey]
	if !found {
		return ErrNoPluginKey
	}

	pluginName, ok := pluginNameValue.(string)
	if !ok {
		return ErrPluginKeyNoString
	}

	// load plugin
	p, err := plugin.Open(pluginName)
	if err != nil {
		return err
	}

	newFunc, err := p.Lookup(newPluginFunction)
	if err != nil {
		return err
	}

	// is it our plugin?
	newPluginFunc, ok := newFunc.(func() (hp.Auth, error))
	if !ok {
		return ErrNoAuthPlugin
	}

	authPlugin, err := newPluginFunc()
	if err != nil {
		return err
	}

	// let's unmarshal context to the plugin object
	if context, found := value[contextKey]; found {
		yamlContext, err := yaml.Marshal(context)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(yamlContext, authPlugin); err != nil {
			return err
		}
	}

	// do we have Init function?
	if init, ok := authPlugin.(hasInit); ok {
		if err := init.Init(); err != nil {
			return err
		}
	}

	*a = Auth{
		Plugin: authPlugin,
	}

	return nil

}

func (a *Auth) GetUser(r *http.Request) (string, error) {
	return a.Plugin.GetUser(r)
}
