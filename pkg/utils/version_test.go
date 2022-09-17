package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestBuildEnvoyFilterNamesAllVersion(t *testing.T) {
	type args struct {
		base string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "build envoy filter names",
			args: args{
				base: "foo",
			},
			want: []string{
				"foo-1.4",
				"foo-1.5",
				"foo-1.6",
				"foo-1.7",
				"foo-1.8",
				"foo-1.9",
				"foo-1.10",
				"foo-1.11",
				"foo-1.12",
				"foo-1.13",
				"foo-1.14",
				"foo-1.15",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildEnvoyFilterNamesAllVersion(tt.args.base)

			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildEnvoyFilterNamesAllVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildEnvoyFilterNames(t *testing.T) {
	type args struct {
		base     string
		versions []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "build envoy filter names",
			args: args{
				base: "foo",
				versions: []string{
					"1.4",
					"1.5",
				},
			},
			want: []string{
				"foo-1.4",
				"foo-1.5",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildEnvoyFilterNames(tt.args.base, tt.args.versions)

			sort.Strings(got)
			sort.Strings(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildEnvoyFilterNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
