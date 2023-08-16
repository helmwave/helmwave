package release

import "regexp"

func (rel *config) Validate() error {
	if rel.Name() == "" {
		return ErrNameEmpty
	}

	if rel.Namespace() == "" {
		rel.Logger().Warnf("namespace is empty. I will use the namespace of your k8s context.")
	}

	if !validateNS(rel.Namespace()) {
		return InvalidNamespaceError{Namespace: rel.Namespace()}
	}

	if err := rel.Uniq().Validate(); err != nil {
		return err
	}

	return nil
}

func validateNS(ns string) bool {
	r := regexp.MustCompile("[a-z0-9]([-a-z0-9]*[a-z0-9])?")

	return r.MatchString(ns)
}
