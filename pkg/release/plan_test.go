// +build ignore unit

package release

import (
	"helm.sh/helm/v3/pkg/action"
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
			want: true,
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
			name: "contains partially",
			args: args{
				targetTags:  []string{"1", "2", "3"},
				releaseTags: []string{"1", "4", "20"},
			},
			want: false,
		},
		{
			name: "contains completely",
			args: args{
				targetTags:  []string{"1", "4"},
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

func TestPlan(t *testing.T) {
	type args struct {
		tags               []string
		enableDependencies bool
	}

	releases := []*Config{
		{
			Name: "release1",
			Tags: []string{"1", "3"},
			Options: action.Upgrade{
				Namespace: "ns",
			},
			DependsOn: []string{"release2@ns"},
		},
		{
			Name: "release2",
			Tags: []string{"2", "3"},
			Options: action.Upgrade{
				Namespace: "ns",
			},
		},
	}
	tests := []struct {
		name     string
		args     args
		wantPlan []*Config
	}{
		{
			name: "empty tags",
			args: args{
				tags: []string{},
			},
			wantPlan: releases,
		},
		{
			name: "tag filter",
			args: args{
				tags:               releases[0].Tags,
				enableDependencies: false,
			},
			wantPlan: []*Config{releases[0]},
		},
		{
			name: "global tag (check release duplication)",
			args: args{
				tags:               []string{"3"},
				enableDependencies: true,
			},
			wantPlan: releases,
		},
		{
			name: "multiple tags",
			args: args{
				tags:               []string{"1", "2"},
				enableDependencies: false,
			},
			wantPlan: []*Config{},
		},
		{
			name: "multiple tags",
			args: args{
				tags:               []string{"1", "3"},
				enableDependencies: false,
			},
			wantPlan: []*Config{releases[0]},
		},
		{
			name: "nonexistent tag",
			args: args{
				tags: []string{"1231231"},
			},
			wantPlan: []*Config{},
		},
		{
			name: "add dependency",
			args: args{
				tags:               releases[0].Tags,
				enableDependencies: true,
			},
			wantPlan: releases,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPlan := Plan(tt.args.tags, releases)
			if len(gotPlan) == 0 && len(tt.wantPlan) == 0 {
				return
			}
			if !reflect.DeepEqual(gotPlan, tt.wantPlan) {
				t.Errorf("Plan() = %v, want %v", gotPlan, tt.wantPlan)
			}
		})
	}
}
