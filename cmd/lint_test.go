package cmd

var exitCode int

func fakeExit(code int) {
	exitCode = code
}
