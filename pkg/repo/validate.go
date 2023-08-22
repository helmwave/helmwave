package repo

import "net/url"

func (c *config) Validate() error {
	if c.Name() == "" {
		return ErrNameEmpty
	}

	if c.URL() == "" {
		return ErrURLEmpty
	}

	if _, err := url.Parse(c.URL()); err != nil {
		return NewInvalidURLError(c.URL())
	}

	return nil
}
