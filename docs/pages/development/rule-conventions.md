# Writing rules

Rules implement:

```go
type Rule interface {
  ID() string
  META() engine.RuleMeta
  Apply(file string, f *hcl.File) []engine.Issue
}
```

Parse with `hclsyntax.ParseConfig`, iterate `body.Blocks`, or use the optional AST walker for nested/expr-heavy checks.
Return issues with precise `hcl.Range`. Keep rules single-purpose and fast.

The ID follows this patter `package.name`
