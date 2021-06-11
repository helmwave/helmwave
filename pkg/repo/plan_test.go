// +build ignore unit

package repo

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
)

func TestGitlabRepo(t *testing.T) {
	const yaml = `
project: sbs
version: 0.7.0


repositories:
  - name: bitnami
    url: https://charts.bitnami.com/bitnami
  - name: cetic
    url: https://cetic.github.io/helm-charts
  - name gitlab 
	url: https://charts.gitlab.io/
  - name: stable 
    url: https://charts.helm.sh/stable

releases:
  - name: gitlab 
    chart: gitlab/gitlab
`
	f, err := helper.CreateFile("helmwave.yml")
	if err != nil {
		t.Error(err)
	}

	_, err = f.WriteString(yaml)
	if err != nil {
		t.Error(err)
	}

	err = f.Close()
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}

}
