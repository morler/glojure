package runtime

import (
	"strings"
	"testing"

	"github.com/glojurelang/glojure/pkg/lang"
)

func TestEvalConvertsEscapingPanicToError(t *testing.T) {
	ns := lang.FindOrCreateNamespace(lang.NewSymbol("runtime.eval-panic-test"))
	lang.InternVar(ns, lang.NewSymbol("panic-value"), lang.FnFunc0(func() interface{} { panic("boom") }), true)
	env := newEnvironment(t.Context(), nil, nil)
	_, err := env.Eval(lang.NewList(lang.NewSymbol("runtime.eval-panic-test/panic-value")))
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("Eval error = %v, want panic converted to error", err)
	}
}

func TestReadEvalReturnsReaderAndEvalErrors(t *testing.T) {
	for _, code := range []string{" )", "missing-symbol-xyz", "(/ 1 0)", "(throw (Exception. \"boom\"))"} {
		got := ReadEval(code)
		if _, ok := got.(error); !ok {
			t.Errorf("ReadEval(%q) = %T %v, want error value", code, got, got)
		}
	}
}
