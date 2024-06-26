package uniqname

import (
	"fmt"
	"regexp"
	"strings"
)

// Separator is a separator between release name and namespace.
const Separator = "@"

var validateRegexp = regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?")

// UniqName is a unique identificator for release.
type UniqName struct {
	name      string
	namespace string
	context   string
}

var _ fmt.Stringer = UniqName{}

// New returns uniqname for provided release name and namespace.
func New(name, namespace, context string) (UniqName, error) {
	u := UniqName{
		name:      name,
		namespace: namespace,
		context:   context,
	}

	return u, u.Validate()
}

func NewFromString(line string) (UniqName, error) {
	parts := strings.Split(line, Separator)

	var u UniqName
	switch len(parts) {
	case 1:
		u = UniqName{name: parts[0]}
	case 2:
		u = UniqName{name: parts[0], namespace: parts[1]}
	case 3:
		u = UniqName{name: parts[0], namespace: parts[1], context: parts[2]}
	default:
		return UniqName{}, NewValidationError(line)
	}

	return u, u.Validate()
}

// GenerateWithDefaultNamespaceContext parses uniqname out of provided line.
// If there is no namespace in line, default namespace will be used.
func GenerateWithDefaultNamespaceContext(line, namespace, context string) (UniqName, error) {
	// ignoring error here because it will likely to be triggered by empty namespace or context
	u, _ := NewFromString(line)

	if u.namespace == "" {
		u.namespace = namespace
	}

	if u.context == "" {
		u.context = context
	}

	return u, u.Validate()
}

// Equal checks whether uniqnames are equal.
func (n UniqName) Equal(a UniqName) bool {
	return n == a
}

// Validate validates this object.
func (n UniqName) Validate() error {
	if !validateRegexp.MatchString(n.name) {
		return NewValidationError(n.String())
	}

	if !validateRegexp.MatchString(n.namespace) {
		return NewValidationError(n.String())
	}

	if !validateRegexp.MatchString(n.context) {
		return NewValidationError(n.String())
	}

	return nil
}

func (n UniqName) String() string {
	str := n.name

	if n.namespace == "" {
		return str
	}

	str += Separator + n.namespace

	if n.context == "" {
		return str
	}

	str += Separator + n.context

	return str
}

func (n UniqName) Empty() bool {
	return n.name == "" && n.namespace == "" && n.context == ""
}
