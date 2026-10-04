package config

import (
	"errors"
	"fmt"
	"strings"
)

func (r Route) Validate() error {
	var errs []error

	if r.Host == "" {
		errs = append(errs, errors.New("host is required"))
	}

	if r.Destination == "" {
		errs = append(errs, errors.New("destination is required"))
	}

	if r.Port < 1 || r.Port > 65535 {
		errs = append(errs, errors.New("destination port must be between 1 and 65535"))
	}

	if !strings.HasPrefix(r.Path, "/") {
		errs = append(errs, errors.New("path must begin with /"))
	}

	return errors.Join(errs...)
}

func (rs Routes) Validate() error {
	var errs []error
	seen := make(map[string]int)

	for i, r := range rs {
		if err := r.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("route %d: %w", i, err))
		}

		key := strings.ToLower(r.Host) + r.Path

		if prev, exists := seen[key]; exists {
			errs = append(errs, fmt.Errorf("route index %d conflicts with route index %d", i, prev))
		}

		seen[key] = i
	}

	return errors.Join(errs...)
}
