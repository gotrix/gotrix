package app

import (
	"testing"
)

func Test_readPathsRecursive(t *testing.T) {
	type args struct {
		dir    string
		suffix string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "this",
			args:    args{dir: "../", suffix: ".go"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readPathsRecursive(tt.args.dir, tt.args.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPathsRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Error("invalid result")
			}
		})
	}
}
