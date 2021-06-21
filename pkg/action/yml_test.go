package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"os"
	"testing"
)

func Test01(t *testing.T) {
	root := "../../action/"
	from := root + "01_helmwave.yml.tpl"
	to := root + "01_helmwave.yml"
	value := "test"

	s := &Yml{
		from,
		to,
	}

	_ = os.Setenv("PROJECT_NAME", value)
	_ = os.Setenv("NAMESPACE", value)

	err := s.Run()

	if err != nil {
		t.Error(err)
	}

	b, err := plan.NewBody(to)
	if err != nil {
		t.Error(err)
	}

	if (value != b.Project) || (value != b.Releases[0].Options.Namespace) {
		t.Error("Failed Test01")
	}

}
