GSpec: a productive Go test framework
=====================================

[![Build Status](https://travis-ci.org/hailiang/gspec.png?branch=master)](https://travis-ci.org/hailiang/gspec)
[![Coverage Status](https://coveralls.io/repos/hailiang/gspec/badge.png?branch=master)](https://coveralls.io/r/hailiang/gspec?branch=master)
[![GoDoc](https://godoc.org/github.com/hailiang/gspec?status.png)](https://godoc.org/github.com/hailiang/gspec)

GSpec is a *concurrent, minimal, extensible and reliable* test framework in Go
that makes it productive to organize and verify the mind model of software.

Highlights:

* *Expressive*: a complete running specification can be organized via both BDD
                and table driven styles.
* *Reliabile*:  the implementation has minimal footprint and is tested with 100%
                coverage.
* *Concurrent*: run test cases concurrently or sequentially.
* *Extensible*: customizable BDD cue words, expectations and test reporters.
* *Compatible*: "go test" is enough to run GSpec tests.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Quick start](#quick-start)
  - [Get GSpec](#get-gspec)
  - [Write tests with GSpec](#write-tests-with-gspec)
  - [Run tests with "go test"](#run-tests-with-go-test)
- [Extend GSpec](#extend-gspec)
    - [Package organization](#package-organization)
  - [Error](#error)
  - [Expectation](#expectation)
  - [Reporter](#reporter)
- [Hack GSpec](#hack-gspec)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Quick start
-----------

###Get GSpec
```
go get -t -u github.com/hailiang/gspec
go test github.com/hailiang/gspec/...
```

###Write tests with GSpec
As Go's convention, write GSpec tests in file xxx_test.go to test code in xxx.go.
```go
// GSpec follows modular design.
import (
	"fmt"

	// core implements core alogrithms of nested test groups with less than 500 lines of code.
	"github.com/hailiang/gspec/core"
	// expectation contains extensible expectation (assertion) helpers.
	exp "github.com/hailiang/gspec/expectation"
	// suite gathers test functions and run them.
	"github.com/hailiang/gspec/suite"
)

// Only one suite.Add is needed for each xxx_test.go file.
var _ = suite.Add(func(s core.S) {
	// BDD cue word is customizible.
	describe, given, when, it := s.Alias("describe"), s.Alias("given"), s.Alias("when"), s.Alias("it")
	// expectation cue word is customizible too.
	expect := exp.Alias(s.FailNow)

	// A BDD example.
	describe("an integer i", func() {
		// setup
		i := 2
		defer func() {
			// teardown (if any)
		}()
		given("another integer j", func() {
			j := 3
			when("j is added to i", func() {
				i += j
				it("should become the sum of original i and j", func() {
					expect(i).Equal(5) // a passing case
				})
			})
			when("j is minused from i", func() {
				i -= j
				it("should become the difference of j minus i", func() {
					expect(i).Equal(4) // a failing case
				})
			})
			when("j is multiplied to i", nil) // a pending case
		})
	})

	// A table-driven example.
	testcase := s.Alias("testcase")
	describe("integer summation", func() {
		for _, c := range []struct{ i, j, sum int }{
			{1, 2, 3}, // a passing case
			{1, 1, 0}, // a failing case
		} {
			testcase(fmt.Sprintf(`%d + %d = %d`, c.i, c.j, c.sum), func() {
				expect(c.i + c.j).Equal(c.sum)
			})
		}
	})
})
```

Write the following go test function in any test file (e.g. all_test.go).

```go
import (
	"testing"
	"github.com/hailiang/gspec/suite"
)

// Defined only once within a package.
func TestAll(t *testing.T) {
	suite.Run(t)
}
```

###Run tests with "go test"
Run all the tests concurrently (sequencially by default) and display errors.
```
go test -concurrent
```
Run all the tests and view the complete specification.
```
go test -v
```
Run only a failing test case (even an entry in the driven table):
```
go test -focus 1/1
```

Extend GSpec
------------
###Package organization
The subpackages are organized with minimal coupling.
```
extension   <- 
core        <- extension
error       <- 
expectation <- error
reporter    <- extension, error
suite       <- core, exntension, reporter
```
1. the core package implements core algorithms of test organization and execution,
   but nothing else. It is extensible through the types defined in the extension
   package.
2. the expectation package implements expectation helpers. It reports expecation
   errors to Fail or FailNow method of interface core.S. core receives and hand
   errors over to reporters without knowing their exact types. expectation package
   can be replaced by any package with an error reporting function of the same
   signature.
3. the error package is responsible for implementing the details of an error, e.g.
   the type of the error, file, line number and the stack trace.
4. the reporter package contains all the builtin test reporters that implement
   extension.Reporter. A reporter gets notifications about the progress of test
   running and gets a complete specification of all the nested test groups,
   including test errors.
5. the suite package integrates all other packages together, providing a quick
   way of test gathering, running and reporting.

###Error
Good error message is mandatory to productive testing.

###Expectation

###Reporter

Hack GSpec
----------

[Design of GSpec](DESIGN.md)


