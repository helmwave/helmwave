package yml

func (yaml *Body) Plan() (plan *Body) {
	plan.Project = yaml.Project
	plan.Version = yaml.Version

	return plan
}
