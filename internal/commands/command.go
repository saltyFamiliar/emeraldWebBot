package commands

type Command struct {
	ParamString string
	Function    func(string) string
	Result      string
	ErrorMsg    string
}

func (cmd *Command) Execute() {
	cmd.Result = cmd.Function(cmd.ParamString)
}
