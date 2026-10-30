package runtime

import (
	"fmt"
	"strings"

	"github.com/glojurelang/glojure/pkg/lang"
	"github.com/glojurelang/glojure/pkg/reader"
)

type (
	// ReadEvalOption is an option for ReadEval.
	ReadEvalOption func(*readEvalOptions)

	readEvalOptions struct {
		// env is the environment to use for evaluation. If not set, the
		// global environment is used.
		env lang.Environment
		// filename is the name of the file being read.
		filename string
	}
)

// WithEnv sets the environment to use for evaluation.
func WithEnv(env lang.Environment) ReadEvalOption {
	return func(o *readEvalOptions) {
		o.env = env
	}
}

// WithFilename sets the filename to use for evaluation.
func WithFilename(filename string) ReadEvalOption {
	return func(o *readEvalOptions) {
		o.filename = filename
	}
}

// ReadEval reads and evaluates a string that may contain one or more
// forms in the global environment.
// ReadEval returns the last evaluated value, or an error value if reading or evaluation fails.
// Errors are returned as values to preserve the original single-result API.
func ReadEval(code string, options ...ReadEvalOption) interface{} {
	var opts readEvalOptions
	for _, opt := range options {
		opt(&opts)
	}
	env := opts.env
	if env == nil {
		env = lang.GlobalEnv
	}
	readerOpts := []reader.Option{
		reader.WithGetCurrentNS(func() *lang.Namespace {
			return env.CurrentNamespace()
		}),
	}
	if opts.filename != "" {
		readerOpts = append(readerOpts, reader.WithFilename(opts.filename))
	}

	r := reader.New(strings.NewReader(string(code)), readerOpts...)

	var lastValue interface{}
	for {
		expr, err := r.ReadOne()
		if err == reader.ErrEOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading %v: %w", opts.filename, err)
		}
		lastValue, err = env.Eval(expr)
		if err != nil {
			return fmt.Errorf("error evaluating %v: %w", opts.filename, err)
		}
	}
	return lastValue
}
