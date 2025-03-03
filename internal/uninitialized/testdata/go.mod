module github.com/danielmmetz/uninitialized/testdata

go 1.17

require (
	github.com/danielmmetz/uninitialized/testdata/external v0.0.0
)

replace github.com/danielmmetz/uninitialized/testdata/external => ./src/github.com/danielmmetz/uninitialized/testdata/external
