Design of GSpec
===============

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Introduction](#introduction)
- [Rationale](#rationale)
- [Test organization](#test-organization)
  - [A overview of test case](#a-overview-of-test-case)
  - [Nested test group](#nested-test-group)
  - [Shared setup/teardown](#shared-setupteardown)
  - [Table-driven test](#table-driven-test)
  - [Test specification](#test-specification)
    - [Alias](#alias)
    - [Description](#description)
  - [Test case gathering](#test-case-gathering)
- [Test composition](#test-composition)
  - [Expectation](#expectation)
    - [Design goals](#design-goals)
    - [Error Message](#error-message)
    - [List of expectations](#list-of-expectations)
  - [Test Double](#test-double)
- [Test execution](#test-execution)
  - [Concurrency](#concurrency)
  - [Focus mode](#focus-mode)
  - [Test time](#test-time)
    - [Timeout](#timeout)
    - [Find slow tests](#find-slow-tests)
  - [Benchmark & coverage](#benchmark-&-coverage)
  - [Options](#options)
    - ["go test" command](#go-test-command)
  - [Auto Test](#auto-test)
- [Test report](#test-report)
  - [Collector](#collector)
  - [Reporter](#reporter)
    - [Helpful messages](#helpful-messages)
    - [Test case failure](#test-case-failure)
    - [Fatal error](#fatal-error)
    - [Builtin reporter](#builtin-reporter)
- [Usage Guidelines](#usage-guidelines)
- [Reference](#reference)
  - [Existing Go Test Frameworks](#existing-go-test-frameworks)
    - [xUnit Style](#xunit-style)
    - [BDD Style](#bdd-style)
    - [Expectations (assertions)](#expectations-assertions)
    - [Mock](#mock)
  - [Links](#links)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Introduction
------------
GSpec is a concurrent, minimal, extensible and reliable unit test framework
in Go that makes it easy to organize and verify the mind model of software.

Rationale
---------
Writing a software is the process of building a model in one's mind and
synthesize programming patterns into a new software, making sure that the
software behaves as expected as the mind model. When the model of software is
too complicated to fit into one's mind all at once, the divide-and-conquer
strategy can be used. The model is divided into multiple parts and the
interactions among those parts are kept minimal, so that each part of the
software can be written and verified separately.

This is the basic idea about unit test.

The builtin test framework for Go (gotest for short, including "go test" command
and "testing" package) is a lightweight test framework that solves most of the
fundamental problems about unit test, including:
* Test organization
    - a test case is a function with specific signature.
    - a test function can be put in any files with name *_test.go.
* Test composition
    - method for error reporting
    - method for terminating or skipping a test case
    - method for measuring test time
* Test execution
    - "go test" command to run tests
    - Run test cases concurrently
    - Timeout detection
    - test coverage analysis
    - blocking analysis
    - run chosen test cases (focus mode by regular expressions matching test
      function names)
* Test report
    - console output

Gotest is superior in the sense that it solves major problems in daily test
in a good way, and contains no unnecessary features. However, it lacks some
features that are also important in unit test:
* Test organization
    - gotest does not help much beyond test functions scattered in multiple test
      files.
    - gotest does not help reducing duplicate setup and teardown code.
* Test composition
    - gotest lacks expectation (assertion) helpers that reduce redundant code.
* Test execution
    - gotest encourages [table-driven test](https://code.google.com/p/go-wiki/wiki/TableDrivenTests),
      but it is not clear how to select and run a single test case in the table.
* Test report
    - gotest emphasize the importance of [good error message](http://golang.org/doc/faq#testing_framework),
      but does not provide any tools to achieve that.
    - gotest does not provide alternative way besides console output.
    - gotest does not provide a clean message when test panics.

All the shortcomings of gotest can be remedied by a minimal framework that
provide the missing features while keeping the solution provided by gotest
intact. GSpec framework should achieve the following goals:
* Test organization
    - A natual way to organize test cases into a complete specification.
    - Common test logic can be shared easily.
* Test composition
    - Extensible expectation helpers.
* Test execution
    - All gotest features should not be broken.
    - Table driven test (with focus mode).
* Test report
    - Extensible test reporters that provide informative and helpful error
      messages.
    - Handle panics gracefully.
* GSpec should be reliable itself.
    - minimal and modular design.
    - 100% test coverage.

Why writing yet another test framework, given there are already many? (
http://code.google.com/p/go-wiki/wiki/Projects#Testing)

* Automatic test is important, especially for a small team that lacks resources
  to test manually.
* "go test" only meets minimal requirement, leaving a gap to fill in a way or
  another.
* None of the existing frameworks fully satisfy all the goals above.
* It is a good way to learn test itself.
* It is fun.

The following sections are organized in the sequence of design decisions, from
major features to minor ones.

Test organization
-----------------

###A overview of test case
A test case needs running some code and verifying the result. Further, the code
run by a test case can be usually seperated into 3 steps:

* setup a test context
* perform an action and verify the result
* teardown the test context (optional)

Sometimes setup/teardown is shared among test cases, so there are 4 situations:

* setup before each test
* teardown after each test
* setup run once before all tests
* teardown run once after all tests

Sometimes both setup/teardown and test logic can be shared among test cases,
and the differences between test cases can be represented with data only and
organized in a table, thus table-driven tests can be applied.

###Nested test group
GSpec uses nested test groups to organize test cases.

Nested test group is a tree structure. Each leaf represents a test case and the
path from the root node to the leaf represents the setup sequence of the test
context for the test case. Reversely, the path from the leaf to the root
represents the teardown sequence of the test context for the test case. e.g.

    a
        b
            c1
            c2

The setup sequence should be: abc1 abc2, and the teardown sequence should be:
c1ba c2ba. (The relative order between c1 and c2 is not important)

###Shared setup/teardown
One way of doing setup/teardown is by hook functions as what RSpec does (before,
after method with argument :each). An alternative way is to define the setup
code directly in the closure and schedule it to run before each test case. As
long as the scheduling logic does not introduce too much overhead and
concurrency can be supported, the latter should be preferred because of its
simplicity. Also, the teardown could be supported simply by a defer statement.
e.g.

    group(func() {
        // setup
        defer func() {
            // teardown
        }()
        group(func() {
            // test case 1
        })
        group(func() {
            // test case 2
        })
    })

A limitation of lacking a dedicated setup method is: when the setup panics,
GSpec can only give the panicking postion but cannot tell it is a leaf node or
it still has children. However, it is a tolerable tradeoff.

It is rare to setup (teardown) a test context before (after) all tests without
worrying about its state changing. In such rare cases, they can be handled
separatedly at the start (end) of the "go test" testing functions.

IMPLEMENTATION:
To implement the tree traversing logic, the path is composed of unique function
IDs, which could be implemented as either the relative position from the root
test group or the address of the function.

It seems that the function address approach is simpler, however, there are two
major flaws of this approach:
1. the function address changes with any code modifications, making it hard to
   implement a useful "focus mode".
2. the function address is always the same within a loop, making it hard to
   support table-driven test.

So the former way is chosen.

###Table-driven test
GSpec should allow table driven tests and group functions nested in arbitrary
ways. e.g.

    group(func() {
        for i, testCase := range table {
            group(func() {
                // test case i
            })
        }
    })

IMPLEMENTATION:
The main challenge of table driven test is: the same closure could be run
multiple times within the loop, so each different run of the same closure should
have different function ID.

###Test specification
GSpec should be able to generate a structured, readable plain text specification
from the tests written. There should be a way to define and collect information
for each level of test group.

Another benefit of plain text specification is that it eliminates the need to
define a unique function name for each test function.

####Alias
GSpec should be able to provide a convenient way to use customized alias names
for the nested group function. e.g. describe, context, it, specify and example.

GSpec should be able to tell what alias name is used for a certain test group.

####Description
GSpec should be able to assotiate a description to each level of test group.

Besides, it also should be able to allow types inserted between test description
so that during refactoring, these types can be found and modified too and won't
get out of sync. Go does not support first class type, so it could be
implemented by variables with zero value. (OPTIONAL, TODO)

####Pending mode
GSpec should allow a test case with only description but no test closure, a
placeholder to delay the test implementation.

No test closure can be easily implemented with a nil closure, which will trigger
a special error sent to the test reporter so that it can be marked as "pending".

###Test case gathering
"go test" gathers test functions with specific function/file naming conventions.
GSpec should also be able to gather tests across test functions/files.

Unlike the group context S, the controller can be shared by all goroutines as
long as carefully locked. So RootFuncs can be defined anywhere and are gathered
by a single controller instance in one "go test" function. This requires
Controller.Start accepts multiple RootFuncs.

RESTRICTION:
There is no finalization hook provided by "go test", so there is no way to know
when all test functions finish. The only possible way is to output the result
for each test gathering. So RootFuncs should be gathered only all at once in one
package.

Test composition
----------------

###Expectation

####Design goals
* Provide structured and helpful error messages.
* Extensible.
* Can be used alone with "go test".
* Fluent interface.
* Do not panic.

####Error Message
An error message should tell the developer clearly:
* Where the failure is
* What is expected
* What is actually got

Should a customizable description supported? No!
* In a BDD-style framework, such descriptions already included in the "it" test
  group. e.g. it("should xxx", func() { expect ... }). So providing such
  optional arguments will be a duplication.
* When used directly with "go test", such messages can be easily written as a
  comment above the expectation line of code, and printed out when the test
  goroutine fails.

An error message should be provided in the form of an error object.
* Plain string message can be easily obtained by method Error().
* Detail structure can still be hold within the object, so that an aggresive 
  reporter can use such information to provide an enhanced visualization such
  as colorful highlighting.

####List of expectations
This is a complete list of expectations, "+" means already supported:
* +To: general expectation
* +Panic: panic and panic with an object
* +Equal: deep equal
*  Is: shallow equal
* +IsType: if a interface is of a type
*  Order comparision: >, <, >=, <=, Within(delta).Of(value)
*  Composition: Not, And, Or
*  True/False: can be supported by Is(true), Is(false)
*  Nil: can be supported by Is(nil)
*  Empty: empty for container types
*  Error: can be supported by Equal
*  Implements: can be supported by IsType
*  String
    - +HasPrefix
    - +HasSuffix
    - +Contains
    - ContainsAny, ContainsRune, EqualFold, Match...
*  Collection: HasLen, All, Any, Exact, InOrder, InPartialOrder ...

###Test Double
(TODO)

Test execution
--------------

###Concurrency
Concurrency is a core feature of Go. "go test" supports concurency at the level
of test functions. With a test framework with nested test group by closures,
it is expected to have dozens (or even hundreds) of test cases written
in one test function. On a quad-core CPU, you probably could just split test
cases into four test functions, but CPU could easily get many cores (dozens or
even hundreds of cores) in the foreseeable future, thus it requires one
goroutine per test case to make the most of the hardware.

To support concurency, test cases should be completed isolated from each other.
No variables should be shared without careful synchronization. It is an illusion
that the variables defined in closures are shared between test cases. Actually,
each test case runs from the root level and variables are allocated on call
stack of its own gouroutine.

When the test cases are organized in a tree of closures, there is no way to know
the whole structure without actually running it. The process has to be
exploratory, starting goroutines on the fly and guiding the test case along a
path from the root down to a certain leaf. The path has to be stored somewhere
and should not be shared between test cases. Thus, scheduling related variables
have to be passed into the goroutine function (RootFunc) as an argument
(s of interface type S). So the group function now becomes the Group method of
S, e.g.

    controller.Start(func(s S) {
        s.Group(func() {
            s.Group(func() {
                // test case 1
            })
            s.Group(func() {
                // test case 2
            })
        })
    })

where the RootFunc is defined as:

    type RootFunc func(s S)

and variable s contains all the variables needed to control which test case to
run, and the controller makes sure each test case run once concurrently
(or sequentially).

RESTRICTION:
* To support concurency, there must be one context variable of type S per
  goroutine. So the Group method cannot be simply defined as a global function,
  and a top-level function of type RootFunc is mandatary for writing tests.
  GSpec has to exchange some simplicty for concurency (This could be compensated
  by defining aliases for the Group method).

###Focus mode
Focus mode allows you specify a single test case and run it. It is especially
useful when you are trying to fix a failed test case and want to print some logs
or attach a debugger.

Focus mode should be implemented as a command line option of "go test", without
affecting the test code. Like this:
    go test -focus <test case id>

As mentioned in the [Shared setup/teardown](#shared-setupteardown) section, the
test case id should not be the address of the closures, but the path of the
serial numbers in the order of execution for each closure, because you do not
want to change the command line to run the same test case every time. The
address of the function will easily change even if a single line is added or
removed, while the path of execution is much more stable.

###Test time
####Timeout
It will be awkward if one of the test cases hangs without knowing which one.
"go test" uses two steps to address this problem:
1. In the testing package, a Timer is used to limit the test time, and panics
   when the time exeeds the limit (pkg/testing/testing.go).
2. If the Timer does not get the chance to run because the code under tests
   saturate all the CPU resouces within the test process, "go test" will detects
   the problem, signal the test process, wait 5 seconds for dumping stack trace,
   and kill the test process (cmd/go/test.go).

Though not friendly enough, it is a comprehensive method. GSpec will do nothing
about timeout for now, unless an obvious better way is found.

####Find slow tests
It is important to keep unit tests run fast so that they can be run as often as
possible during the development, thus, it is important to find slow tests and
improve them.

It can be implemented by measuring each run of a test case and store the
time duration. Test reporter is responsible to analyze the result and find slow
test cases.

###Benchmark & coverage
What "go test" provides are good enough. Just don't break them.

###Options
Options of GSpec should able to set hard coded or via CLI flags. Flag should
have higher priority than hard coded value so that can be changed at runtime.

####"go test" command
* "go test" command is the only way to run GSpec tests.
* GSpec should try to be consistent with test flags that "go test" have.
    - "-outputdir": place output files in the specified directory.
    - "-parallel": Allow parallel execution of test functions that call
      t.Parallel.
    - "-v": Verbose output.

###Auto Test
(TODO)

Test report
-----------

###Collector
A collector is an internal object embedded in the test controller that collects
the outputs from the tests, reconstructing the tree structure of nested test
groups with their descriptions and results.

The implementation of a collector must assume being called concurrently out of
order. To reconstruct a tree structure, it would be easy if a parent node is
always collected before its children, so a collector should collect before the
start of each test group. To collect the result of each test case, a collector
also needs to collect at the end of each test group.

###Reporter
A reporter is responsible for reporting the test progress and display the test
result. It is defined as an public interface so that GSpec is able to use any
customized reporter.

####Helpful messages
A reporter should provide messages that help the developer to find and fix
issues.

* Complete: the information provided should be complete to minimize further
  investigation as much as possible.
* Clean: unrelavent noise should be reduced as much as possible.
* Structured: the information could be provided as a structured object instead
  of a string to provide more capability of customized reportring.

####Test case failure
When an expectation (assertion) fails, It should *not* panic. Instead, GSpec
should provide a method (S.Fail) to record the error and continue running other
test cases. t.Fail should be called to notify "go test" there are failures.

This also means that there should be only one expectation for a test case.

Expectations should be provided as a separate package (gspec/expect).

NOTE: the S.Fail method should allow to be called from another goroutine.

####Fatal error
When test code panics, "go test" prints the error and terminates immediately.
GSpec should respect this design:
* Simple way is just to repanicking in goroutine, which causes the process
  terminates.
* For a cleaner error message, print a better message and call os.Exit.

GSpec itself could have internal bugs and fail, when it happens, GSpec should
panic immediately. (fail fast)

####Builtin reporter
GSpec should contain a builtin reporter:
* It should be as simple as possible, advanced reporters should be implemented
  in an extention package (gspec/reporter?).
* It is a plain text reporter, without colors.
* In default mode, it should display nothing but the failure information.
* In verbose mode, it should display an additional progress line (.F*).

RESTRICTION:
* If there is a single test controller in a single "go test" function, then a
single reporter instance that serves it needs no locking. (Recommended)
needed.
* If there are multiple test controllers, there should be a single reporeter
instance to seve them all and it needs to be locked.
* If there are multiple test controllers and one reporter per controller, then
"go test" functions should not be run concurrently.
* Otherwise, test cases could be printed interwaved in console.

Usage Guidelines
----------------
* Write tests that matter.
* Test simple abstractions rather than complex details.
    * Avoid mocking when possible.
* Write abstract specifications in text; write concrete examples in code.
* One expectation per test case.

Reference
---------
###Existing Go Test Frameworks

####xUnit Style
* launchpad.net/gocheck
* github.com/stretchr/testify

####BDD Style
* smartystreets.github.io/goconvey
* github.com/onsi/ginkgo
* github.com/franela/goblin
* github.com/orfjackal/gospec
* github.com/stesla/gospecify
* github.com/azer/mao
* github.com/pranavraja/zen (forked from mao)

####Expectations (assertions)
* github.com/onsi/gomega
* launchpad.net/gocheck
* github.com/stretchr/testify/assert

####Mock
* code.google.com/p/gomock
* https://github.com/qur/withmock (gomock companion)
* github.com/stretchr/testify/mock
* github.com/jvshahid/mock4go

###Links
* http://betterspecs.org/
* http://doctoc.herokuapp.com/
