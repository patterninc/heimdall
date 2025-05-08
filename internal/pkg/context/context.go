package context

import (
	"encoding/json"
)

type Context map[string]any

func New(v interface{}) *Context {

	data, err := json.Marshal(v)
	if err != nil {
		panic(`cannot marshal json`)
	}

	value := make(map[string]any)

	if err := json.Unmarshal(data, &value); err != nil {
		panic(`cannot unmarshal json`)
	}

	return (*Context)(&value)

}

func (c *Context) UnmarshalYAML(unmarshal func(interface{}) error) error {

	value := make(map[string]any)

	if err := unmarshal(&value); err != nil {
		return err
	}

	*c = value

	return nil

}

func (c *Context) UnmarshalJSON(data []byte) error {

	value := make(map[string]any)

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*c = value

	return nil

}

func (c *Context) Unmarshal(v interface{}) error {

	// let's marshal our data first
	data, err := json.Marshal(*c)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)

}

func (c *Context) String() string {

	if c == nil {
		return ``
	}

	data, _ := json.Marshal(*c)

	return string(data)

}
