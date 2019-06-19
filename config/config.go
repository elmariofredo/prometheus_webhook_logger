package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// Secret is a string that must not be revealed on marshaling.
type Secret string

// MarshalYAML implements the yaml.Marshaler interface.
func (s Secret) MarshalYAML() (interface{}, error) {
	if s != "" {
		return "<secret>", nil
	}
	return nil, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Secrets.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Secret
	return unmarshal((*plain)(s))
}

// LoadConfig parses the YAML input into a Config.
func LoadConfig(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	log.V(1).Infof("Loaded config:\n%+v", cfg)
	return cfg, nil
}

// LoadConfigFile parses the given YAML file into a Config.
func LoadConfigFile(filename string) (*Config, []byte, error) {
	log.V(1).Infof("Loading configuration from %q", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := LoadConfig(string(content))
	if err != nil {
		return nil, nil, err
	}

	//resolveFilepaths(filepath.Dir(filename), cfg)
	return cfg, content, nil
}

// Config is the top-level configuration for JIRAlert's config file.
type Config struct {
	WebhookAddress string `yaml:"WebhookAddress" json:"WebhookAddress"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline" json:"-"`
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return checkOverflow(c.XXX, "config")
}

func checkOverflow(m map[string]interface{}, ctx string) error {
	if len(m) > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		log.Warningf("unknown fields in %s: %s", ctx, strings.Join(keys, ", "))
	}
	return nil
}
