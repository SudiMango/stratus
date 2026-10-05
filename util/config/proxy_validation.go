package config

import (
	"crypto/tls"
	"errors"
	"fmt"
)

func (p Proxy) Validate() error {
	var errs []error

	if !p.Http.Enabled && !p.Https.Enabled {
		errs = append(errs, errors.New("\t[Proxy] At least one listener must be enabled"))
	}

	if p.Http.RedirectToHttps && !p.Https.Enabled {
		errs = append(errs, errors.New("\t[Proxy] HTTP cannot redirect when HTTPS is disabled"))
	}

	if p.Http.RedirectToHttps && !p.Http.Enabled {
		errs = append(errs, errors.New("\t[Proxy] HTTP redirect cannot be enabled when HTTP is disabled"))
	}

	if p.Http.Enabled &&
		p.Https.Enabled &&
		p.Http.Port == p.Https.Port {
		errs = append(errs, errors.New("\t[Proxy] HTTP and HTTPS cannot use the same port"))
	}

	if p.MaxHeaderBytes < 0 {
		errs = append(errs, errors.New("\t[Proxy] MaxHeaderBytes cannot be less than 0"))
	}

	if p.ReadTimeout < 0 ||
		p.ReadHeaderTimeout < 0 ||
		p.WriteTimeout < 0 ||
		p.IdleTimeout < 0 {
		errs = append(errs, errors.New("\t[Proxy] Timeout values cannot be less than 0"))
	}

	errs = append(errs, p.Http.Validate())
	errs = append(errs, p.Https.Validate())

	return errors.Join(errs...)
}

func (h Http) Validate() error {
	if !h.Enabled {
		return nil
	}

	if h.Port < 1 || h.Port > 65535 {
		return errors.New("\t[Proxy] HTTP port must be between 1 and 65535")
	}

	return nil
}

func (h Https) Validate() error {
	if !h.Enabled {
		return nil
	}

	var errs []error

	if h.Port < 1 || h.Port > 65535 {
		errs = append(errs, errors.New("\t[Proxy] HTTPS port must be between 1 and 65535"))
	}

	if h.Cert == "" {
		errs = append(errs, errors.New("\t[Proxy] HTTPS certificate path is required"))
	}

	if h.Key == "" {
		errs = append(errs, errors.New("\t[Proxy] HTTPS key path is required"))
	}

	if _, err := tls.LoadX509KeyPair(h.Cert, h.Key); err != nil {
		errs = append(errs, fmt.Errorf("\t[Proxy] Error loading HTTP cert and key: %v", err))
	}

	return errors.Join(errs...)
}
