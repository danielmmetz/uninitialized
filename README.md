# Uninitialized
Uninitialized: check for uninitialized but required struct fields within composite literals.

### Why?

I've found myself wanting required fields within structs as sometimes a zero value is either unworkable or impractical.
While constructors can help to mitigate this need, constructors can suffer from the same issue--a new field is added to the
type being constructed, but one has forgotten to update the constructor too.

#### Caveats

  * This linter is not yet battle tested.
  * There was a [similar proposal for Go 2](https://github.com/golang/go/issues/28348). It was rejected by the core team.

### Install
```bash
go get -u github.com/danielmmetz/uninitialized
```

### Use

Annotate struct fields with a tag: `required:"true"`.

Example:
```go
type Foo struct {
    Member Bar `required:"true"`
}
```

Composite literals of `Foo` that do not explicitly set the member field `Member` will then be flagged.

Example:
```go
func main() {
    _ = Foo{}  // `Foo missing required keys: [Member]`
    _ = Foo{Member: Bar{}}  // OK }
```

### Run

```bash
uninitialized: check for uninitialized but required struct fields within composite literals

Usage: uninitialized [-flag] [package]
```

Example:
```bash
uninitialized ./...
```
