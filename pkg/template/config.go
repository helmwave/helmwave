package template

type Tpl struct {
	From string
	To   string
}

func (t *Tpl) Render() error {
	return Tpl2yml(t.From, t.To, nil)
}
