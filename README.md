GSpec: a productive Go test framework
=====================================

[![Build Status](https://travis-ci.org/hailiang/gspec.png?branch=master)](https://travis-ci.org/hailiang/gspec)
[![Coverage Status](https://coveralls.io/repos/hailiang/gspec/badge.png?branch=master)](https://coveralls.io/r/hailiang/gspec?branch=master)
[![GoDoc](https://godoc.org/github.com/hailiang/gspec?status.png)](https://godoc.org/github.com/hailiang/gspec)

GSpec is an *expressive, reliable, concurrent and extensible* Go test framework
that makes it productive to organize and verify the mind model of software.

* *Expressive*: a complete runnable specification can be organized via both BDD
                and table driven styles.
* *Reliable*:   the implementation has minimal footprint and is tested with 100%
                coverage.
* *Concurrent*: test cases can be executed concurrently or sequentially.
* *Extensible*: customizable BDD cue words, expectations and test reporters.
* *Compatible*: "go test" is sufficient but not mandatory to run GSpec tests.

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
```bash
go get -t -u github.com/hailiang/gspec
go test github.com/hailiang/gspec/...
```

###Write tests with GSpec
According to the [convention of Go](http://golang.org/doc/code.html#Testing),
write GSpec tests in file xxx_test.go to test code in xxx.go.
```go
import (
	"fmt"

	// GSpec follows modular design.

	// core implements core alogrithms of nested test groups within 500 lines of code.
	"github.com/hailiang/gspec/core"
	// expectation contains extensible expectation (assertion) helpers.
	exp "github.com/hailiang/gspec/expectation"
	// suite gathers top level test functions (core.TestFunc) and run them.
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

Write the following go test function for only once in any test file within the
package (e.g. all_test.go).

```go
import (
	"testing"
	"github.com/hailiang/gspec/suite"
)

// Defined only once within a package.
func TestAll(t *testing.T) {
	suite.Test(t)
}
```

###Run tests with "go test"
Run all the tests concurrently (sequencially by default) and display errors.
```bash
go test -concurrent
```
Run all the tests and view the complete specification.
```bash
go test -v
```
Run only a failing test case (even it is an entry in the driven table):
```bash
go test -focus 1/1
```

Understand GSpec
----------------
###Test organization
GSpec tests are defined in a top level function of signature core.TestFunc.
```go
type TestFunc func(S)
```
S is an interface that provides methods for defining nested test groups and
reporting test errors.
```go
type S interface {
	Alias(name string) DescFunc
	Fail(err error)
	FailNow(err error)
}

type DescFunc func(description string, f func())
```
Within a TestFunc, an alias function of signature core.Desc needs to be defined
for the cue word of BDD style test. e.g.
```go
	describe := s.Alias("describe")
```
Then the "describe" function can be used to define a test group.
```go
	describe("website login", func() {
	})
```
GSpec will concatenate the cue word and the description argument, so the complete
description of the test group becomes: "describe website login".

Those DescFuncs can be nested, forming a tree of nested test groups. Each leaf
test group corresponds to a test case. To run a specific test case, GSpec
executes from the top level TestFunc down to the leaf test group, and the test
groups not on path are ignored. GSpec will guarantee that each test case is
executed only once.

###Test error
Within the closure of a test group, tests can fail and the test error is reported
to GSpec via S.Fail or S.FailNow method. The only difference between FailNow and
Fail is that FailNow stops the execution of the test case immediately while Fail
continues after reporting the error. Fail and FailNow only affect the currently
running test case.

core knows nothing about the internals of error objects, and only transfers
them to reporters, but there are two exceptions: extension.PanicError and
extension.PendingError:

1. core captures a panicking error, put it in a PanicError object and transfer
   it the same as other errors.
2. When a DescFunc is called with a nil test closure, it is treated as a pending
   test case, and a PendingError is passed to the test reporter.

###Test execution
To actually run the tests, a core.Controller object is needed. Controller.Start
is responsible for starting top level test functions.
```go
func (c *Controller) Start(path Path, concurrent bool, funcs ...TestFunc) error
```
The path parameter is used to specify a path within the tree of nested test
groups. An empty path means the top level of test group should be executed,
including all its descendants.

suite package provides a convenient way to gather and run TestFuncs. suite.Add
adds a TestFunc to a global slice and suite.Test runs all the gathered tests.
Other parameters like path and concurrent are provided by command-line flag.

###Test report
Test results are reported via extension.Reporter interface by the core.
```go
type Reporter interface {
	Start()
	End(groups TestGroups)
	Progress(g *TestGroup, s *Stats)
}

type TestGroups []*TestGroup

type TestGroup struct {
	ID          string
	Description string
	Error       error
	Duration    time.Duration
	Children    TestGroups
}

type Stats struct {
	Total   int
	Ended   int
	Failed  int
	Pending int
}
```
Start gets called before all tests started and End gets called after all tests
end. The complete and final test result are passed to a reporter as the groups
parameter of method End. Progress method is used to report the progress during
the test execution.

core does not contain an implementation of a reporter. Multiple external
reporters can be provided when constructing a new Controller via
core.NewController. These reporters are notified one by one.
```go
func NewController(reporters ...ext.Reporter) *Controller
```

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
2. the error package is responsible for implementing the details of an error, e.g.
   the type of the error, file, line number and the stack trace.
3. the expectation package implements expectation helpers. It reports expecation
   errors to Fail or FailNow method of interface core.S. core receives and hand
   errors over to reporters without knowing their exact types. expectation package
   can be replaced by any package with an error reporting function of the same
   signature.
4. the reporter package contains all the builtin test reporters that implement
   extension.Reporter. A reporter gets notifications about the progress of test
   running and gets a complete specification of all the nested test groups,
   including test errors.
5. the suite package integrates all other packages together, providing a quick
   way of test gathering, running and reporting.

Extend GSpec
------------
###Error
Good error message is mandatory to productive testing.

###Expectation

###Reporter

Hack GSpec
----------
###Test
[Design of GSpec](DESIGN.md)


