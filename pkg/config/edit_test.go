package config

import (
	"errors"
	"testing"

	"github.com/urfave/cli"
)

type TestRunner struct{}

func (t TestRunner) Run(cmd string, args ...string) ([]byte, error) {
	return []byte(""), nil
}

type TestFailRunner struct{}

func (t TestFailRunner) Run(cmd string, args ...string) ([]byte, error) {
	return []byte(""), errors.New("The file /app/asdhasdjasjd] does not exist.")
}

func TestHandleEdit(t *testing.T) {
	runner = TestRunner{}

	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		shouldErr bool
	}{
		{
			name: "should open editor with config file",
			args: args{
				c: &cli.Context{},
			},
			wantErr:   false,
			shouldErr: false,
		},
		{
			name: "should return error - file not found",
			args: args{
				c: &cli.Context{},
			},
			wantErr:   true,
			shouldErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldErr == true {
				runner = TestFailRunner{}
			} else {
				runner = TestRunner{}
			}
			if err := HandleEdit(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("HandleEdit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
