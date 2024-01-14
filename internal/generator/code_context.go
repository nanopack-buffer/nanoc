package generator

type CodeContext struct {
	vars map[string]bool
}

func (c CodeContext) IsVariableInScope(varName string) bool {
	_, ok := c.vars[varName]
	return ok
}

func (c CodeContext) AddVariableToScope(varName string) {
	c.vars[varName] = true
}

func (c CodeContext) RemoveVariableFromScope(varName string) {
	delete(c.vars, varName)
}

func (c CodeContext) ClearVariablesFromScope() {
	c.vars = map[string]bool{}
}

func (c CodeContext) NewLoopVar() string {
	i := 0
	v := loopVar[i]
	for c.IsVariableInScope(v) {
		i += 1
		v = loopVar[i]
	}
	c.AddVariableToScope(v)
	return v
}
