package config_test

import (
	"testing"

	"github.com/SudiMango/stratus/util/config"
)

func TestRouteValidate(t *testing.T) {
	validRoute := config.Route{
		Host:        "app.example.com",
		Destination: "localhost",
		Port:        3000,
		Path:        "/",
	}

	tests := []struct {
		name        string
		mutate      func(*config.Route)
		wantMessage string
	}{
		{name: "literal route"},
		{
			name: "remainder wildcard route",
			mutate: func(route *config.Route) {
				route.Path = "/test/{rest...}"
			},
		},
		{
			name: "host is required",
			mutate: func(route *config.Route) {
				route.Host = ""
			},
			wantMessage: "host is required",
		},
		{
			name: "whitespace-only host is invalid",
			mutate: func(route *config.Route) {
				route.Host = "   "
			},
			wantMessage: "host",
		},
		{
			name: "destination is required",
			mutate: func(route *config.Route) {
				route.Destination = ""
			},
			wantMessage: "destination is required",
		},
		{
			name: "whitespace-only destination is invalid",
			mutate: func(route *config.Route) {
				route.Destination = "   "
			},
			wantMessage: "destination",
		},
		{
			name: "destination port must be positive",
			mutate: func(route *config.Route) {
				route.Port = 0
			},
			wantMessage: "destination port",
		},
		{
			name: "path must start with slash",
			mutate: func(route *config.Route) {
				route.Path = "test"
			},
			wantMessage: "path must begin",
		},
		{
			name: "invalid wildcard pattern",
			mutate: func(route *config.Route) {
				route.Path = "/test/{rest...}/more"
			},
			wantMessage: "path",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			route := validRoute
			if test.mutate != nil {
				test.mutate(&route)
			}

			assertValidationResult(t, route.Validate(), test.wantMessage)
		})
	}
}

func TestRoutesValidateConflicts(t *testing.T) {
	base := config.Route{
		Host:        "app.example.com",
		Destination: "localhost",
		Port:        3000,
		Path:        "/users/{id}",
	}

	tests := []struct {
		name        string
		routes      config.Routes
		wantMessage string
	}{
		{
			name: "different exact and subtree routes",
			routes: config.Routes{
				{
					Host:        "app.example.com",
					Destination: "localhost",
					Port:        3000,
					Path:        "/test",
				},
				{
					Host:        "app.example.com",
					Destination: "localhost",
					Port:        4000,
					Path:        "/test/{rest...}",
				},
			},
		},
		{
			name: "identical route",
			routes: config.Routes{
				base,
				base,
			},
			wantMessage: "conflicts",
		},
		{
			name: "host comparison is case insensitive",
			routes: config.Routes{
				base,
				{
					Host:        "APP.EXAMPLE.COM",
					Destination: "localhost",
					Port:        4000,
					Path:        "/users/{id}",
				},
			},
			wantMessage: "conflicts",
		},
		{
			name: "trailing DNS period is equivalent",
			routes: config.Routes{
				base,
				{
					Host:        "app.example.com.",
					Destination: "localhost",
					Port:        4000,
					Path:        "/users/{id}",
				},
			},
			wantMessage: "conflicts",
		},
		{
			name: "wildcard names do not distinguish patterns",
			routes: config.Routes{
				base,
				{
					Host:        "app.example.com",
					Destination: "localhost",
					Port:        4000,
					Path:        "/users/{name}",
				},
			},
			wantMessage: "conflicts",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertValidationResult(t, test.routes.Validate(), test.wantMessage)
		})
	}
}
