package taskmaster

// RecoverableFunc is a wrapper for recoverable goruntine
func RecoverableFunc(routine func(), recoverCallback func(interface{})) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				recoverCallback(r)
			}
		}()
		routine()
	}
}
