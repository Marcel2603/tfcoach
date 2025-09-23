package engine_test

import (
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/testutil"
)

func TestEngine_WithStubRule(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"a.tf": `x`}}
	e := engine.New(src)
	e.Register(testutil.AlwaysFlag{Id: "t.id", Message: "m"})
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("want 1, got %d", len(issues))
	}
}

func TestEngine_WithManyStubRule(t *testing.T) {
	src := testutil.MemSource{Files: map[string]string{"a.tf": `x`}}
	e := engine.New(src)
	e.RegisterMany([]engine.Rule{testutil.AlwaysFlag{Id: "t.id", Message: "m"},
		testutil.AlwaysFlag{Id: "t.id", Message: "2", Match: "x"}})
	issues, err := e.Run(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("want 1, got %d", len(issues))
	}
}
