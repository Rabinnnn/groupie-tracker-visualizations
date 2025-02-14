package location

import "testing"

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name        string
		args        args
		wantCity    string
		wantCountry string
	}{
		{
			name:        "philadelphia-usa",
			args:        args{s: "philadelphia-usa"},
			wantCity:    "philadelphia",
			wantCountry: "usa",
		},

		{
			name:        "north_carolina-usa",
			args:        args{s: "north_carolina-usa"},
			wantCity:    "north carolina",
			wantCountry: "usa",
		},

		{
			name:        "dunedin-new_zealand",
			args:        args{s: "dunedin-new_zealand"},
			wantCity:    "dunedin",
			wantCountry: "new zealand",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotCity, gotCountry := Parse(tt.args.s)
				if gotCity != tt.wantCity {
					t.Errorf("Parse() gotCity = %v, want %v", gotCity, tt.wantCity)
				}
				if gotCountry != tt.wantCountry {
					t.Errorf("Parse() gotCountry = %v, want %v", gotCountry, tt.wantCountry)
				}
			},
		)
	}
}
