package pow

import (
	"testing"
)

func Test_isValid(t *testing.T) {
	type args struct {
		zeroCount int
		b         []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"1",
			args{
				zeroCount: 1,
				b:         []byte{0x1F, 0xFF, 0xFF},
			},
			true,
		},
		{
			"2",
			args{
				zeroCount: 1,
				b:         []byte{0xFF, 0xFF, 0xFF},
			},
			false,
		},
		{
			"2",
			args{
				zeroCount: 4,
				b:         []byte{0x0F, 0xFF, 0xFF},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValid(tt.args.zeroCount, tt.args.b); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
