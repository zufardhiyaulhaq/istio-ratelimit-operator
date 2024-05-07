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
				"foo-1.16",
				"foo-1.17",
				"foo-1.18",
				"foo-1.19",
				"foo-1.20",
				"foo-1.21",
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
		{
			name: "build envoy filter 1.15",
			args: args{
				base: "foo",
				versions: []string{
					"1.15",
				},
			},
			want: []string{
				"foo-1.15",
			},
		},
		{
			name: "build envoy filter 1.16",
			args: args{
				base: "foo",
				versions: []string{
					"1.16",
				},
			},
			want: []string{
				"foo-1.16",
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
