package main

import (
	"testing"
)

func Test_detectRepository(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name      string
		args      args
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name: "filled",
			args: args{
				owner: "a",
				repo:  "b",
			},
			wantOwner: "a",
			wantRepo:  "b",
			wantErr:   false,
		},
		{
			name: "no owner",
			args: args{
				owner: "",
				repo:  "b",
			},
			wantOwner: "storz",
			wantRepo:  "b",
			wantErr:   false,
		},
		{
			name: "no repo",
			args: args{
				owner: "a",
				repo:  "",
			},
			wantOwner: "a",
			wantRepo:  "prdiff",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := detectRepository(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantOwner {
				t.Errorf("detectRepository() got = %v, wantOwner %v", got, tt.wantOwner)
			}
			if got1 != tt.wantRepo {
				t.Errorf("detectRepository() got1 = %v, wantOwner %v", got1, tt.wantRepo)
			}
		})
	}
}

func Test_parseArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name     string
		args     args
		wantBase string
		wantHead string
	}{
		{
			name: "no argument",
			args: args{
				args: []string{},
			},
			wantBase: "",
			wantHead: "",
		},
		{
			name: "one argument",
			args: args{
				args: []string{"a"},
			},
			wantBase: "",
			wantHead: "a",
		},
		{
			name: "one argument with sep",
			args: args{
				args: []string{"a...b"},
			},
			wantBase: "a",
			wantHead: "b",
		},
		{
			name: "two argument",
			args: args{
				args: []string{"a", "b"},
			},
			wantBase: "a",
			wantHead: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseArgs(tt.args.args)
			if got != tt.wantBase {
				t.Errorf("parseArgs() got = %v, wantBase %v", got, tt.wantBase)
			}
			if got1 != tt.wantHead {
				t.Errorf("parseArgs() got1 = %v, wantBase %v", got1, tt.wantHead)
			}
		})
	}
}
