package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"
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

	for i, c := range p.Certificates {
		if strings.TrimSpace(c.Host) == "" {
			errs = append(errs, fmt.Errorf("\t[Proxy] HTTPS certificate host for cert index %d is required", i))
		}

		if strings.TrimSpace(c.Cert) == "" {
			errs = append(errs, fmt.Errorf("\t[Proxy] HTTPS certificate path for host %s is required", c.Host))
		}

		if strings.TrimSpace(c.Key) == "" {
			errs = append(errs, fmt.Errorf("\t[Proxy] HTTPS key path for host %s is required", c.Host))
		}

		loadedCert, err := tls.LoadX509KeyPair(c.Cert, c.Key)
		if err != nil {
			errs = append(errs, fmt.Errorf("\t[Proxy] Error loading HTTPS cert and key for host %s: %v", c.Host, err))
			continue
		}

		x509Cert, err := x509.ParseCertificate(loadedCert.Certificate[0])
		if err != nil {
			errs = append(errs, fmt.Errorf("[Proxy] Failed to parse certificate for host %s: %w", c.Host, err))
			continue
		}

		if err := x509Cert.VerifyHostname(c.Host); err != nil {
			errs = append(errs, fmt.Errorf("[Proxy] Certificate %q is not valid for host %q: %w", c.Cert, c.Host, err))
		}

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

	return errors.Join(errs...)
}
