package config

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (r Route) Validate() error {
	var errs []error

	if strings.TrimSpace(r.Host) == "" {
		errs = append(errs, errors.New("host is required"))
	}

	if strings.TrimSpace(r.Destination) == "" {
		errs = append(errs, errors.New("destination is required"))
	}

	if r.Port < 1 || r.Port > 65535 {
		errs = append(errs, errors.New("destination port must be between 1 and 65535"))
	}

	if !strings.HasPrefix(r.Path, "/") {
		errs = append(errs, errors.New("path must begin with /"))
	} else {
		pattern := NormalizeHost(r.Host) + r.Path

		if err := registerPattern(http.NewServeMux(), pattern); err != nil {
			errs = append(errs, fmt.Errorf("path pattern is invalid: %w", err))
		}
	}

	return errors.Join(errs...)
}

func (rs Routes) Validate() error {
	var errs []error
	mux := http.NewServeMux()

	for i, r := range rs {
		if err := r.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("\t[Routes] route index %d: %w", i, err))
		}

		pattern := NormalizeHost(r.Host) + r.Path
		if err := registerPattern(http.NewServeMux(), pattern); err != nil {
			continue
		}

		if err := registerPattern(mux, pattern); err != nil {
			errs = append(errs, fmt.Errorf("\t[Routes] route index %d conflicts with another route: %w", i, err))
		}
	}

	return errors.Join(errs...)
}

// Helpers

func registerPattern(mux *http.ServeMux, pattern string) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%v", recovered)
		}
	}()

	mux.HandleFunc(pattern, func(http.ResponseWriter, *http.Request) {})
	return nil
}
