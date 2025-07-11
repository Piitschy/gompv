package gompv

const base = "mpv"

type Command struct {
	Flags map[string]string
	Args  []string
}

func (c Command) String() string {
	cmd := base
	for flag, value := range c.Flags {
		if value == "true" {
			cmd += " --" + flag
		} else {
			cmd += " --" + flag + "=" + value
		}
	}
	for _, arg := range c.Args {
		cmd += " " + arg
	}
	return cmd
}

func (c Command) Slice() []string {
	cmd := []string{base}
	for flag, value := range c.Flags {
		if value == "true" {
			cmd = append(cmd, "--"+flag)
		} else {
			cmd = append(cmd, "--"+flag+"="+value)
		}
	}
	cmd = append(cmd, c.Args...)
	return cmd
}

func NewCommand() *Command {
	return &Command{
		Flags: make(map[string]string),
		Args:  []string{},
	}
}

func (c *Command) AddFlag(flag, value string) {
	if value != "" {
		c.Flags[flag] = value
	} else {
		c.Flags[flag] = "true"
	}
}

func (c *Command) AddArg(arg string) {
	c.Args = append(c.Args, arg)
}
