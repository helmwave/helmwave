package monitor

func (c *config) Validate() error {
	if c.Name() == "" {
		return ErrNameEmpty
	}

	if c.TotalTimeout < c.IterationTimeout {
		return ErrLowTotalTimeout
	}

	if c.Interval == 0 {
		return ErrLowInterval
	}

	err := c.subConfig.Validate()
	if err != nil {
		return NewSubMonitorError(err)
	}

	return nil
}
