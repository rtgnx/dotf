package godots

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestExpand(t *testing.T) {
	vars := Variables{Global: map[string]string{
		"name":   "abc",
		"domain": "example.com",
		"home":   "/home/$USER",
		"email":  "$name@$domain",
	}}

	vars.Expand()

	validHome := path.Join("/home", os.Getenv("USER"))
	validEmail := fmt.Sprintf("%s@%s", vars.Global["name"], vars.Global["domain"])

	if home, ok := vars.Global["home"]; !ok || home != validHome {
		t.Errorf("TestExpand() => expected %s got %s", validHome, home)
	}

	if email, ok := vars.Global["email"]; !ok || email != validEmail {
		t.Errorf("TestExpand() => expected %s got %s", validEmail, email)
	}
}

func TestVariables_Lookup(t *testing.T) {
	type fields struct {
		Global map[string]string
	}
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "test existing key",
			fields: fields{Global: map[string]string{"key": "value"}},
			args:   args{"key", "default"},
			want:   "value",
		},
		{
			name:   "test existing key",
			fields: fields{Global: map[string]string{}},
			args:   args{"key", "default"},
			want:   "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := &Variables{
				Global: tt.fields.Global,
			}
			if got := vars.Lookup(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("Variables.Lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
