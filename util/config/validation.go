package config

import (
	"errors"
)

type Validatable interface {
	Validate() error
}

func (c Config) Validate() error {
	return errors.Join(
		c.Proxy.Validate(),
		c.Routes.Validate(),
	)
}
