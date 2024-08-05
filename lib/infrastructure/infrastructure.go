package infrastructure

import (
	"fmt"
)

type infrastructure struct {
	cache map[string]*Component
}

type Component struct {
	enable bool
	name   string

	c ConnectClose
}

func (c *Component) hasConfiged() bool {
	return c.enable
}

func (c *Component) Init() error {
	if !c.hasConfiged() {
		return fmt.Errorf("component not configured, nest name: %s", c.name)
	}

	if err := c.c.Connect(); err != nil {
		return fmt.Errorf("connect to %s failure, nest error: %v", c.name, err)
	}

	return nil
}

func (c *Component) Close() error {
	if c.c == nil {
		return nil
	}

	return c.c.Close()
}

type ConnectClose interface {
	Connect() error
	Close() error

	UnMarshalConfig([]byte) error
}

var ins = &infrastructure{
	cache: map[string]*Component{},
}

func Register(name string, c ConnectClose) {
	ins.cache[name] = &Component{
		name: name,
		c:    c,
	}
}

func LoadConfig(c Config) (*Component, error) {
	config, err := c.MarshalConfig()
	if err != nil {
		return nil, err
	}

	component, ok := ins.cache[c.Name()]
	if !ok {
		return nil, fmt.Errorf("[%s] not register", c.Name())
	}

	if err := component.c.UnMarshalConfig(config); err != nil {
		return nil, err
	}
	component.enable = true
	ins.cache[c.Name()] = component

	return component, nil
}
