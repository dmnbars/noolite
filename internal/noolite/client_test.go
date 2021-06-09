package noolite

import "testing"

func Test_calcCrc(t *testing.T) {
	tests := []struct {
		data []byte
		want byte
	}{
		{
			data: []byte{171, 2, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 172},
			want: 188,
		},
		{
			data: []byte{171, 2, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 0, 0, 72, 0, 172},
			want: 4,
		},
		{
			data: []byte{173, 1, 0, 25, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 174},
			want: 203,
		},
		{
			data: []byte{173, 1, 0, 26, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 174},
			want: 204,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := calcCrc(tt.data); got != tt.want {
				t.Errorf("calcCrc() = %v, want %v", got, tt.want)
			}
		})
	}
}
