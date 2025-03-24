<h1 align="center">
  go-golden
</h1>

<p align="center">
  <strong>
    Yet another Go package for working with <code>*.golden</code> test files,
    with a focus on simplicity.
  </strong>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/jimeh/go-golden"><img src="https://img.shields.io/badge/%E2%80%8B-reference-387b97.svg?logo=go&logoColor=white" alt="Go Reference"></a>
  <a href="https://github.com/jimeh/go-golden/actions"><img src="https://img.shields.io/github/actions/workflow/status/jimeh/go-golden/ci.yml?logo=github" alt="Actions Status"></a>
  <a href="https://codeclimate.com/github/jimeh/go-golden"><img src="https://img.shields.io/codeclimate/coverage/jimeh/go-golden.svg?logo=code%20climate" alt="Coverage"></a>
  <a href="https://github.com/jimeh/go-golden/issues"><img src="https://img.shields.io/github/issues-raw/jimeh/go-golden.svg?style=flat&logo=github&logoColor=white" alt="GitHub issues"></a>
  <a href="https://github.com/jimeh/go-golden/pulls"><img src="https://img.shields.io/github/issues-pr-raw/jimeh/go-golden.svg?style=flat&logo=github&logoColor=white" alt="GitHub pull requests"></a>
  <a href="https://github.com/jimeh/go-golden/blob/master/LICENSE"><img src="https://img.shields.io/github/license/jimeh/go-golden.svg?style=flat" alt="License Status"></a>
</p>

## Import

```go
import "github.com/jimeh/go-golden"
```

## Usage

```go
func TestExampleMyStruct(t *testing.T) {
    got, err := json.Marshal(&MyStruct{Foo: "Bar"})
    require.NoError(t, err)

    want := golden.Do(t, got)

    assert.Equal(t, want, got)
}
```

The above example will read/write to:

- `testdata/TestExampleMyStruct.golden`

The call to `golden.Do()` is equivalent to:

```go
if golden.Update() {
    golden.Set(t, got)
}
want := golden.Get(t)
```

To update the golden file (have `golden.Update()` return `true`), simply set the
`GOLDEN_UPDATE` environment variable to one of `1`, `y`, `t`, `yes`, `on`, or
`true` when running tests.

## Documentation

Please see the
[Go Reference](https://pkg.go.dev/github.com/jimeh/go-golden#section-documentation)
for documentation and examples.

## License

[MIT](https://github.com/jimeh/go-golden/blob/master/LICENSE)
