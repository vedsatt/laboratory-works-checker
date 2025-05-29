package checker

func (c *Checker) splitedTests(compilerCh chan struct {
	msg string
	err error
}) (string, error) {
	out := <-compilerCh
	if out.err != nil {
		return "", out.err
	}
	if out.msg != "" {
		return out.msg, nil
	}

	return "", nil
}
