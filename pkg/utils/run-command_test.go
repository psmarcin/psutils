package utils

import (
	"reflect"
	"testing"
)

func TestRealRunner_Run(t *testing.T) {
	type args struct {
		command string
		args    []string
	}
	tests := []struct {
		name    string
		r       RealRunner
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "run echo",
			r:    RealRunner{},
			args: args{
				command: "echo",
				args:    []string{"123"},
			},
			want:    []byte{49, 50, 51, 10},
			wantErr: false,
		},
		{
			name: "run asdasda",
			r:    RealRunner{},
			args: args{
				command: "asdasda",
				args:    []string{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RealRunner{}
			got, err := r.Run(tt.args.command, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RealRunner.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RealRunner.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
