package xtime

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "27-05-1967",
			args:    args{s: "27-05-1967"},
			want:    time.Date(1967, 5, 27, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},

		{
			name:    "32-05-1967",
			args:    args{s: "32-05-1967"},
			want:    time.Date(1967, 5, 32, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},

		{
			name:    "32-05-1967-909",
			args:    args{s: "32-05-1967-909"},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := Parse(tt.args.s)
				if (err != nil) != tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
