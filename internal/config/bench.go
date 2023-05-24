package config

import (
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// K9sBench the name of the benchmarks config file.
var K9sBench = "bench"

type (
	// Bench tracks K9s styling options.
	Bench struct {
		Benchmarks *Benchmarks `yaml:"benchmarks"`
	}

	// Benchmarks tracks K9s benchmarks configuration.
	Benchmarks struct {
		Defaults   Benchmark              `yaml:"defaults"`
		Services   map[string]BenchConfig `yam':"services"`
		Containers map[string]BenchConfig `yam':"containers"`
	}

	// Auth basic auth creds.
	Auth struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}

	// Benchmark represents a generic benchmark.
	Benchmark struct {
		C int `yaml:"concurrency"`
		N int `yaml:"requests"`
	}

	// HTTP represents an http request.
	HTTP struct {
		Method  string      `yaml:"method"`
		Host    string      `yaml:"host"`
		Path    string      `yaml:"path"`
		HTTPS   bool        `yaml:"https"`
		HTTP2   bool        `yaml:"http2"`
		Body    string      `yaml:"body"`
		Headers http.Header `yaml:"headers"`
	}

	// ServiceResolution represents settings for how to resolve a k8s Service to an accessible IP or hostname.
	ServiceResolution struct {
		// Mode can be:
		//  - "ConfiguredHost": (default) Use the hostname or IP address configured in `HTTP.Host`. Assumes the service is reachable via that hostname or IP address.
		//  - "ClusterIP": Access the service directly via its ClusterIP. Assumes the ClusterIP is reachable, i.e. k9s is running within the cluster or has VPN access to it.
		//  - "LoadBalancerIngress": Access the service via its load-balancer ingress (`status.loadBalancer.ingress`). Only works for LoadBalancer services that have a load-balancer ingress; otherwise produces an error.
		Mode  string        `yaml:"mode"`
	}

	// BenchConfig represents a service benchmark.
	BenchConfig struct {
		Name string
		C    int  `yaml:"concurrency"`
		N    int  `yaml:"requests"`
		Auth Auth `yaml:"auth"`
		HTTP HTTP `yaml:"http"`
		ServiceResolution ServiceResolution `yaml:"serviceResolution"`
	}
)

const (
	// DefaultC default concurrency.
	DefaultC = 1
	// DefaultN default number of requests.
	DefaultN = 200
	// DefaultMethod default http verb.
	DefaultMethod = "GET"
)

func newBenchmark() Benchmark {
	return Benchmark{
		C: DefaultC,
		N: DefaultN,
	}
}

// Empty checks if the benchmark is set.
func (b Benchmark) Empty() bool {
	return b.C == 0 && b.N == 0
}

func newBenchmarks() *Benchmarks {
	return &Benchmarks{
		Defaults: newBenchmark(),
	}
}

// NewBench creates a new default config.
func NewBench(path string) (*Bench, error) {
	s := &Bench{Benchmarks: newBenchmarks()}
	err := s.load(path)
	return s, err
}

// Reload update the configuration from disk.
func (s *Bench) Reload(path string) error {
	return s.load(path)
}

// Load K9s benchmark configs from file.
func (s *Bench) load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(f, &s)
}

// DefaultBenchSpec returns a default bench spec.
func DefaultBenchSpec() BenchConfig {
	return BenchConfig{
		C: DefaultC,
		N: DefaultN,
		HTTP: HTTP{
			Method: DefaultMethod,
			Path:   "/",
		},
	}
}
