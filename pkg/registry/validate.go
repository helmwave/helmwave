package registry

func (c *config) Validate() error {
	if c.Host() == "" {
		return ErrNameEmpty
	}

	return nil
}
