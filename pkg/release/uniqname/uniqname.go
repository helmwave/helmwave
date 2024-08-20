package uniqname

import (
	"fmt"
	"regexp"
	"strings"
)

// Separator is a separator between release Name and Namespace.
const Separator = "@"

var (
	NamespaceRegexp   = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
	KubecontextRegexp = regexp.MustCompile("^[a-z0-9]([-_:/a-z0-9]*[a-z0-9])?$")
	ReleaseRegexp     = regexp.MustCompile("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$")
)

// UniqName is a unique identificator for release.
type UniqName struct {
	Name      string
	Namespace string
	Context   string
}

var _ fmt.Stringer = UniqName{}

// New returns uniqname for provided release Name and Namespace.
func New(name, namespace, context string) (UniqName, error) {
	u := UniqName{
		Name:      name,
		Namespace: namespace,
		Context:   context,
	}

	return u, u.Validate()
}

func NewFromString(line string) (UniqName, error) {
	parts := strings.Split(line, Separator)

	var u UniqName
	switch len(parts) {
	case 1:
		u = UniqName{Name: parts[0]}
	case 2:
		u = UniqName{Name: parts[0], Namespace: parts[1]}
	case 3:
		u = UniqName{Name: parts[0], Namespace: parts[1], Context: parts[2]}
	default:
		return UniqName{}, u.Error(line)
	}

	return u, u.Validate()
}

// Equal checks whether uniqnames are equal.
func (n UniqName) Equal(a UniqName) bool {
	return n == a
}

// Validate validates this object.
func (n UniqName) Validate() error {
	if !ReleaseRegexp.MatchString(n.Name) {
		return n.Error(n.Name)
	}

	if !NamespaceRegexp.MatchString(n.Namespace) {
		return n.Error(n.Namespace)
	}

	if n.Context != "" && !KubecontextRegexp.MatchString(n.Context) {
		return n.Error(n.Context)
	}

	return nil
}

func (n UniqName) String() string {
	str := n.Name

	if n.Namespace == "" {
		return str
	}

	str += Separator + n.Namespace

	if n.Context == "" {
		return str
	}

	str += Separator + n.Context

	return str
}

func (n UniqName) Empty() bool {
	return n.Name == "" && n.Namespace == "" && n.Context == ""
}
