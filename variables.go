package godots

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Expand variable values with environment or
// previously defined variables
func (vars *Variables) Expand() {
	for k, v := range vars.Global {
		vars.Global[k] = os.Expand(v, func(s string) string {
			if val, ok := vars.Global[s]; ok {
				return val
			}
			return os.Getenv(s)
		})
	}
}

// Lookup variable and return its value or default string
func (vars *Variables) Lookup(key, fallback string) string {
	if v, ok := vars.Global[key]; ok {
		return v
	}

	return fallback
}

func (vars *Variables) ReadIn(r io.Reader) error {
	defer vars.Expand()
	return yaml.NewDecoder(r).Decode(&vars.Global)
}
