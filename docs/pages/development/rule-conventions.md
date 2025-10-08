# Writing rules

Rules implement:

```go
type Rule interface {
  ID() string
  META() RuleMeta
  Apply(file string, f *hcl.File) []Issue
  Finish() []Issue
}
```

Parse with `hclsyntax.ParseConfig`, iterate `body.Blocks`, or use the optional AST walker for nested/expr-heavy checks.
Return issues with precise `hcl.Range`. Keep rules single-purpose and fast.

Depending on what the rule needs to assert, you may report issues for each file independently (in `Apply`) or collect
information and report after all files have been checked (in `Finish`).

The ID follows this pattern: `package.name`
