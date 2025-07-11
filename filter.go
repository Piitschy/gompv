package mpv

import "fmt"

type Filter struct {
	Operator string
	In       []string
	Out      string
}

func NewFilter(operator string, args ...string) *Filter {
	return &Filter{
		Operator: operator,
		In:       args,
		Out:      "ao",
	}
}

func (f *Filter) SetTarget(target string) {
	f.Out = target
}

func (f *Filter) String() string {
	if len(f.In) == 0 {
		return f.Operator
	}
	args := fmt.Sprintf("[%s]", f.In[0])
	for _, arg := range f.In[1:] {
		args += fmt.Sprintf(" [%s]", arg)
	}
	return fmt.Sprintf("%s %s [%s]", args, f.Operator, f.Out)
}
