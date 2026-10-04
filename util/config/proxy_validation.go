package config

import "errors"

func (p Proxy) Validate() error {
	var errs []error

	if !p.Http.Enabled && !p.Https.Enabled {
		errs = append(errs, errors.New("At least one listener must be enabled"))
	}

	if p.Http.RedirectToHttps && !p.Https.Enabled {
		errs = append(errs, errors.New("HTTP cannot redirect when HTTPS is disabled"))
	}

	if p.Http.Enabled &&
		p.Https.Enabled &&
		p.Http.Port == p.Https.Port {
		errs = append(errs, errors.New("HTTP and HTTPS cannot use the same port"))
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
		return errors.New("HTTP port must be between 1 and 65535")
	}

	return nil
}

func (h Https) Validate() error {
	if !h.Enabled {
		return nil
	}

	var errs []error

	if h.Port < 1 || h.Port > 65535 {
		errs = append(errs, errors.New("HTTPS port must be between 1 and 65535"))
	}

	if h.Cert == "" {
		errs = append(errs, errors.New("HTTPS certificate path is required"))
	}

	if h.Key == "" {
		errs = append(errs, errors.New("HTTPS key path is required"))
	}

	return errors.Join(errs...)
}
