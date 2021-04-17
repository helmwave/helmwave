// +build ignore unit

package release

import (
	"reflect"
	"testing"
)

func Test_normalizeTagList(t *testing.T) {
	type args struct {
		tags []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "idempotent",
			args: args{
				tags: []string{"1", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "sort",
			args: args{
				tags: []string{"3", "2", "1"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "trim",
			args: args{
				tags: []string{" 1", "2 ", " 3 "},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeTagList(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeTagList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkTagInclusion(t *testing.T) {
	type args struct {
		targetTags  []string
		releaseTags []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no target tags",
			args: args{
				targetTags:  []string{},
				releaseTags: []string{"1"},
			},
			want: false,
		},
		{
			name: "no release tags",
			args: args{
				targetTags:  []string{"1"},
				releaseTags: []string{},
			},
			want: false,
		},
		{
			name: "equal tags",
			args: args{
				targetTags:  []string{"1"},
				releaseTags: []string{"1"},
			},
			want: true,
		},
		{
			name: "contains",
			args: args{
				targetTags:  []string{"1", "2", "3"},
				releaseTags: []string{"1", "4", "20"},
			},
			want: true,
		},
		{
			name: "doesn't contain",
			args: args{
				targetTags:  []string{"1", "2", "3"},
				releaseTags: []string{"4", "5", "6"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkTagInclusion(tt.args.targetTags, tt.args.releaseTags); got != tt.want {
				t.Errorf("checkTagInclusion() = %v, want %v", got, tt.want)
			}
		})
	}
}
