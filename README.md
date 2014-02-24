GSpec
=====

[![Build Status](https://travis-ci.org/hailiang/gspec.png?branch=master)](https://travis-ci.org/hailiang/gspec)
[![Coverage Status](https://coveralls.io/repos/hailiang/gspec/badge.png?branch=master)](https://coveralls.io/r/hailiang/gspec?branch=master)
[![GoDoc](https://godoc.org/github.com/hailiang/gspec?status.png)](https://godoc.org/github.com/hailiang/gspec)

GSpec is a concurrent, minimal, extensible and reliable testing framework in Go
that makes it easy to organize and verify the mind model of software. It
supports both BDD style and table driven testing.

(under development).

Highlights:

* Concurrent: one goroutine per test case (sequential running also supported).
* Natual:     BDD and table driven style are integrated natually. Use either one or both to fit your test scenario.
* Reliabile:  the design is minimal and orthogonal; the code is tested under 100% coverage.
* Extensible: Customizable BDD cue words, expectations and test reporters.
* Separable:  subpackages are organized with minimal coupling.
* Compatible: "go test" is enough to run GSpec tests (However, it does not depend on "testing" package).
* Succinct:   the core implementation is less than 500 lines of code.

Design
------

[Core](DESIGN.md)

[Expectations](expectation/DESIGN.md)

Examples
--------
###Concurrent

```go
import (
	"testing"

	"github.com/hailiang/gspec/core"
	exp "github.com/hailiang/gspec/expectation"
	"github.com/hailiang/gspec/suite"
)

var _ = suite.Add(func(s core.S) {
	describe, given, when, it := s.Alias("describe"), s.Alias("given"), s.Alias("when"), s.Alias("it")
	expect := exp.Alias(s.Fail)

	describe("an integer i", func() {
		i := 2
		given("another integer j", func() {
			j := 3
			when("j is added to i", func() {
				i += j
				it("should become the sum of original i and j", func() {
					expect(i).Equal(5)
				})
			})
		        // more scenarios here.
		})

		// more scenarios here.
	})

	// more tests here.
})

func TestAll(t *testing.T) {
	suite.Run(t, false)
}
```

Separable
---------
The subpackages are organized with minimal coupling.

1. core and expectation does not know each other. 
2. core and reporter communicate via the interface defined in extension.
3. core receives and transfers errors to reporter without knowing their exact types.

```
extension   <- 
core        <- extension
errors      <- 
expectation <- errors
reporter    <- extension, errors
suite       <- core, exntension, reporter
```


