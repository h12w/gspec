GSpec: a productive Go test framework
=====================================

[![wercker status](https://app.wercker.com/status/d5ed30f0d03d4a5210f65659ead13888/s "wercker status")](https://app.wercker.com/project/bykey/d5ed30f0d03d4a5210f65659ead13888)
[![GoDoc](https://godoc.org/h12.me/gspec?status.png)](https://godoc.org/h12.me/gspec)

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
- [Understand GSpec](#understand-gspec)
  - [Test organization](#test-organization)
  - [Test error](#test-error)
  - [Expectation](#expectation)
  - [Test execution](#test-execution)
  - [Test report](#test-report)
- [Extend GSpec](#extend-gspec)
  - [Expectation](#expectation-1)
  - [Reporter](#reporter)
- [Hack GSpec](#hack-gspec)
  - [Design document](#design-document)
  - [Package organization](#package-organization)
  - [Test](#test)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Quick start
-----------

###Get GSpec
```bash
go get -u -f h12.me/gspec
go test h12.me/gspec/...
```

###Write tests with GSpec
According to the [convention of Go](http://golang.org/doc/code.html#Testing),
write GSpec tests in file xxx_test.go to test code in xxx.go.
```go
import (
	"fmt"

	"h12.me/gspec"
)

// Only one gspec.Add is needed for each xxx_test.go file.
var _ = gspec.Add(func(s gspec.S) {
	// BDD cue word is customizible.
	describe, given, when, it := s.Alias("describe"), s.Alias("given"), s.Alias("when"), s.Alias("it")
	// expectation cue word is customizible too.
	expect := gspec.Expect(s.FailNow)

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
	"h12.me/gspec"
)

// Defined only once within a package.
func TestAll(t *testing.T) {
	gspec.Test(t)
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
executes from the top level TestFunc down to the leaf test group, ignoring test
groups that are not on path. GSpec will guarantee that each test case is
executed only once.

###Test error
Good error message is very important to productive testing. To achieve it, text
error messages are not good enough, an error object is passed directly to allow
the test reporter determine how to render it. The test error object is simply an
object satisfying the error interface. Each error object has an Error method
that can returns a string message, which can be the fallback method for
rendering.

A test error object is passed to core via S.Fail or S.FailNow method. The
differences between Fail and FailNow are:

1. FailNow stops the execution of the test case immediately  but cannot be called
   from another goroutine spawned by the user.
2. Fail continues after reporting the error, and can be called from any
   goroutine.

Note that both Fail and FailNow only record the first error and ingoring the
later ones.

core does not care about the specific type of error objects, except two cases:
extension.PanicError and extension.PendingError:

1. core captures a panicking error, wrap it in a PanicError object and report it
   the same way as other errors.
2. When a DescFunc is called with a nil test closure, it is treated as a pending
   test case, and a PendingError is passed to the test reporter.

###Expectation
Usually there is no need to call S.Fail or S.FailNow directly, because the
expectation package will handle it.

First an alias function of signature expecation.ExpectFunc needs to be defined
for the cue word of the expecation. e.g.
```go
expect := exp.Alias(s.FailNow)
```
It does rot have to be named as "expect", any valid Go variable name is
possible. Usually s.FailNow should be used, unless you want to test the
expecation within another goroutine.

An ExpectFunc accepts the actual value and returns an expectation.Actual object.
The Actual object has a general method "To" to check against an expected value
with a specific type of expectation.Checker.
```go
type ExpectFunc func(actual interface{}) *Actual

func (a *Actual) To(check Checker, expected interface{})

type Checker func(actual, expected interface{}, skip int) error
```

Fluent methods for builtin checkers are defined directly in the Actual object to
allow more succinct code, e.g.
```go
expect(i).Equal(2)
```

###Test execution
To actually run the tests, a core.Controller object is needed. Controller.Start
is responsible for starting top level test functions.
```go
func (c *Controller) Start(path Path, concurrent bool, funcs ...TestFunc) error
```
The path parameter is used to specify a path within the tree of nested test
groups. An empty path means the top level of test group should be executed,
including all its descendants.

gspec package provides a convenient way to gather and run TestFuncs. gspec.Add
adds a TestFunc to a global slice and gspec.Test runs all the gathered tests.
Other parameters like path and concurrent are provided by command-line flags.

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
Reporter.Start gets called before all tests started and Reporter.End gets called
after all tests end. The complete and final test result are passed to a reporter
as the groups parameter of Reporter.End. Reporter.Progress method is used to
report the progress during the test execution.

core does not contain an implementation of a reporter. Multiple external
reporters can be provided when constructing a new Controller via
core.NewController. These reporters will be notified one by one.
```go
func NewController(reporters ...ext.Reporter) *Controller
```

Extend GSpec
------------
###Expectation
To create a customized expectation of your own, just write a checker function of
signature expecation.Checker.

In the checker function, the actual and the expected value are compared in
specific way. If the expectation passes, just return nil, otherwise, an error
object should be returned.

The error package is intended to make it easier to write customized error types.
error.ExpectError has already defined the basic format of an expecation error,
including file and line number, and error.CompareError defines the basic format
for comparing two values.

###Reporter
Currently GSpec has a text-based reporter defined in the reporter package. The
interface is clearly define in the extension package, and it should not be hard
to write a reporter of your own.

Hack GSpec
----------
It is welcome to make any improvements to GSpec itself. Here are some information
that might help with it.

###Design document
GSpec has a comprehensive [design document](DESIGN.md), including the rationales
of every design decisions.

###Package organization
The subpackages are organized with minimal coupling.
```
extension   <-
core        <- extension
error       <-
expectation <- error
reporter    <- extension, error
gspec       <- core, exntension, reporter
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
5. the gspec package integrates all other packages together, providing a quick
   way of test gathering, running and reporting.


###Test
GSpec is thouroughly checked and tested inlcuding:

1. go vet
2. golint
3. go test -race
4. go test -cover

There is a bash script check.sh will do all the items above automatically.
