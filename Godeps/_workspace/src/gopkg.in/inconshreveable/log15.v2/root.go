package log15

import (
	"gopkg.in/inconshreveable/log15.v2/term"
	"os"
)

var (
	root          Logger
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	if term.IsTty(os.Stdout.Fd()) {
		StdoutHandler = StreamHandler(os.Stdout, TerminalFormat())
	}

	if term.IsTty(os.Stderr.Fd()) {
		StderrHandler = StreamHandler(os.Stderr, TerminalFormat())
	}

	root = &logger{[]interface{}{}, &swapHandler{handler: StdoutHandler}}
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// Debug is a convenient alias for Root().Debug
func Debug(msg string, ctx ...interface{}) {
	root.Debug(msg, ctx...)
}

// Info is a convenient alias for Root().Info
func Info(msg string, ctx ...interface{}) {
	root.Info(msg, ctx...)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg string, ctx ...interface{}) {
	root.Warn(msg, ctx...)
}

// Error is a convenient alias for Root().Error
func Error(msg string, ctx ...interface{}) {
	root.Error(msg, ctx...)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg string, ctx ...interface{}) {
	root.Crit(msg, ctx...)
}
