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

func TestContains(t *testing.T) {
	type args struct {
		a location
		b location
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Exact Match",
			args: args{
				a: "Texas, USA",
				b: "Texas, USA",
			},
			want: true,
		},

		{
			name: "State Match",
			args: args{
				a: "Texas",
				b: "Texas, USA",
			},
			want: true,
		},

		{
			name: "Country Match",
			args: args{
				a: "USA",
				b: "Texas, USA",
			},
			want: true,
		},

		{
			name: "Different Cities in same country",
			args: args{
				a: "Washington, USA",
				b: "Texas, USA",
			},
			want: false,
		},

		{
			name: "Fuzzy Country Match",
			args: args{
				a: "Seattle, Washington, USA ",
				b: "Washington, USA ",
			},
			want: true,
		},

		{
			name: "Fuzzy Country Match",
			args: args{
				a: "Washington, USA ",
				b: "Seattle, Washington, USA ",
			},
			want: true,
		},

		{
			name: "Fuzzy Country Match",
			args: args{
				a: "Washington, USA",
				b: "Seattle, Washington, USA ",
			},
			want: true,
		},

		{
			name: "Invalid separator",
			args: args{
				a: "Seattle, Washington, USA ",
				b: "Washington$USA ",
			},
			want: false,
		},

		{
			name: "Invalid separator",
			args: args{
				a: "Seattle, Washington, USA ",
				b: "Washington.USA ",
			},
			want: false,
		},

		{
			name: "No match",
			args: args{
				a: "Seattle, Washington, USA",
				b: "Hongkong, CN ",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := Contains(tt.args.a, tt.args.b); got != tt.want {
					t.Errorf("Contains() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestGetCityCountry(t *testing.T) {
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
			name:        "Base",
			args:        args{s: "Washington, USA"},
			wantCity:    "Washington",
			wantCountry: "USA",
		},

		{
			name:        "City with space in name",
			args:        args{s: "San Jose, California"},
			wantCity:    "San Jose",
			wantCountry: "California",
		},

		{
			name:        "Base with dot separator",
			args:        args{s: "Washington. USA"},
			wantCity:    "Washington",
			wantCountry: "USA",
		},

		{
			name:        "City with space in name, with dot separator",
			args:        args{s: "San Jose. California"},
			wantCity:    "San Jose",
			wantCountry: "California",
		},

		{
			name:        "City, State, Country format",
			args:        args{s: "Seattle, Washington, USA"},
			wantCity:    "Seattle, Washington",
			wantCountry: "USA",
		},

		{
			name:        "Fictional City, State, Country format",
			args:        args{s: "San Jose, Washington, USA"},
			wantCity:    "San Jose, Washington",
			wantCountry: "USA",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotCity, gotCountry := GetCityCountry(tt.args.s)
				if gotCity != tt.wantCity {
					t.Errorf("GetCityCountry() gotCity = %q, want %q", gotCity, tt.wantCity)
				}
				if gotCountry != tt.wantCountry {
					t.Errorf("GetCityCountry() gotCountry = %q, want %q", gotCountry, tt.wantCountry)
				}
			},
		)
	}
}
