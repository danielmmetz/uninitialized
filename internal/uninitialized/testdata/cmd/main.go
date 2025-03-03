package main

import "github.com/danielmmetz/uninitialized/testdata/external"

func main() {
	_ = external.External{ // want `External missing required keys: \[Required\]`
		Optional: true,
	}
}
