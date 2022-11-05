package main

import "testing"

func TestDiv(t *testing.T) {
	type args struct {
		v1 int
		v2 int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"f1", args{v1: 10, v2: 2}, 5},
		{"f2", args{v1: 8, v2: 1}, 8},
		{"f3", args{v1: 1, v2: 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Div(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}
