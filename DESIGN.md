Design of GSpec
===============

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](http://doctoc.herokuapp.com/)*

- [Introduction](#introduction)
- [Yet Another One?](#yet-another-one)
- [Features](#features)
  - [Enhancing "go test"](#enhancing-go-test)
    - [Concurrency](#concurrency)
    - [Table driven tests](#table-driven-tests)
    - [Test organization](#test-organization)
    - ["go test" command](#go-test-command)
  - [Test Case](#test-case)
    - [Shared setup/teardown (BDD style)](#shared-setupteardown-bdd-style)
    - [Shared test (Table driven style)](#shared-test-table-driven-style)
  - [Nested Test Group](#nested-test-group)
  - [Table-driven Testing](#table-driven-testing)
  - [Concurrency](#concurrency-1)
  - [Test case organization](#test-case-organization)
  - [Test Specification](#test-specification)
    - [Alias](#alias)
    - [Description](#description)
    - [Listener](#listener)
  - [Reporter](#reporter)
    - [Helpful messages](#helpful-messages)
    - [Test case failure](#test-case-failure)
    - [Fatal error](#fatal-error)
    - [Builtin reporter](#builtin-reporter)
  - [Focus Mode](#focus-mode)
  - [Test Time](#test-time)
    - [Timeout](#timeout)
    - [Find slow tests](#find-slow-tests)
  - [Benchmark & coverage](#benchmark-&-coverage)
  - [Options](#options)
  - [Test Double](#test-double)
  - [Auto Test](#auto-test)
- [Usage Guidelines](#usage-guidelines)
- [Existing Go Testing Frameworks](#existing-go-testing-frameworks)
  - [xUnit Style](#xunit-style)
  - [BDD Style](#bdd-style)
  - [Expectations (assertions)](#expectations-assertions)
  - [Mock](#mock)
- [Reference](#reference)

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

This is the basic idea about unit testing.

The builtin test framework for Go (gotest for short, including "go test" command
and "testing" package) is a lightweight test framework that solves most of the
fundamental problems about unit testing, including:
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
    - run chosen test cases (regular expressions matching test function names)

Gotest is superior in the sense that it solves major problems in daily testing
in a good way, and contains no unnecessary features. However, it lacks some
features that are also important in unit testing:
* gotest does not help much with test organization beyond test functions
  scattered in multiple test files.
* gotest does not provide enough facilities to help reducing redundancy in
  testing code.
* gotest emphasize the importance of (good error message)[http://golang.org/doc/faq#testing_framework],
  but does not provide any tools to achieve that.

All the short shortcomings of gotest can be remedied by a minimal framework that
provide the missing features while keeping the solution provided by gotest
intact. The framework (GSpec) should provide the following features:
* Test cases can be organized in nested test groups so that the whole test suite
  of a package can form a complete specification.
    - Common setup/teardown/test logic can be shared among test cases via test
      groups.
    - Good error messages are achieved by attaching informative descriptions to
      each level of test groups.
* Expectation (assertion) helpers are provided to reduce the redundancy in test
  code.
* GSpec should be modular and extensible.
* No existing solutions of gotest are broken by GSpec. Especially,
    - organize test cases in multiple files
    - "go test" command
    - Concurrency
    - Timeout
    - Ability to run chosen test case (focus mode)
* GSpec should be reliable by robust design and 100% test coverage.

Why writing yet another test framework, given there are already many? (
http://code.google.com/p/go-wiki/wiki/Projects#Testing)

* Automatic test is important, especially for a small team that lacks resources
  to test manually.
* "go test" only meets minimal requirement, leaving a gap to fill in a way or
  another.
* None of the existing frameworks fully satisfy all the goals above.
* It is a good way to learn testing itself.
* It is fun.

Features
--------

The following sections are organized in the sequence of design decisions, from
major features to minor ones.

###Test Case
A test case needs running some code and verifying the result. Further, the code
run by a test case can be usually seperated into 3 steps:

* setup a test context
* perform an action and verify the result
* teardown the test context (optional)

####Shared setup/teardown (BDD style)
Sometimes setup/teardown is shared between test cases, so there are 4 situations:

* setup before each test
* teardown after each test
* setup run once before all tests
* teardown run once after all tests

"go test" does nothing to help the developer with shared setup/teardown code.

xUnit style test frameworks implement each step above as virtual methods in a
base class or interface. They do provide a complete solution, including possible
support for concurrency. However, They have to introduce unavoidable boilerplate
code of defining derived test classes, and nested test group is not straight
forward to implement with derived classes.

BDD style frameworks like RSpec use closures to provide an ad hoc and natual
way to specify test cases. Pros of this method include:

* A natual short test description rather than a function name
* Setup/teardown code can be nested to form multiple levels of test groups
* The tests form a readable and runnable specification

GSpec will choose BDD style as its basic form of test cases.

####Shared test (Table driven style)
Sometimes both setup/teardown and test logic can be shared between test cases,
and the differences between test cases can be represented with data only and
organized in a table, thus table-driven tests can be applied.

GSpec should also support table driven tests, allowing table driven tests to be
embedded in BDD style test group, or vice versa.

###Nested Test Group
When a testing context is shared between test cases, the context should be setup
(teardown) before (after) each test to isolate test cases from each other.

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
   implement a useful "focus mode". So the former
2. the function address is always the same within a loop, making it hard to
   support table-driven testing.

So the former way is chosen.

###Table-driven Testing
[Table driven test](https://code.google.com/p/go-wiki/wiki/TableDrivenTests) is
the recommended way when possible.

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
(s of interface type S). e.g.

    scheduler.Start(func(s S) {
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
run, and the scheduler makes sure each test case run once concurrently
(or sequentially).

RESTRICTION:
* To support concurency, there must be one context variable of type S per
  goroutine. So the Group method cannot be simply defined as a global function,
  and a top-level function of type RootFunc is mandatary for writing tests.
  GSpec has to exchange some simplicty for concurency (This could be compensated
  by defining aliases for the Group method).

###Test case organization
"go test" gathers test functions with specific function/file naming conventions.
GSpec should also be able to gather tests across test functions/files.

Unlike the group context S, the scheduler can be shared by all goroutines as
long as carefully locked. So RootFuncs can be defined anywhere and are gathered
by a single scheduler instance in one "go test" function. This requires
Scheduler.Start accepts multiple RootFuncs.

RESTRICTION:
There is no finalization hook provided by "go test", so there is no way to know
when all test functions finish. The only possible way is to output the result
for each test gathering. So RootFuncs should be gathered only all at once in one
package.

###Test Specification
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

####Listener
A listener is an internal object embedded in the test scheduler that collects
the outputs from the tests, reconstructing the tree structure of nested test
groups with their descriptions and results.

The implementation of a listener must assume being called concurrently out of
order. To reconstruct a tree structure, it would be easy if a parent node is
always collected before its children, so a listener should collect before the
start of each test group. To collect the result of each test case, a listener
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
* If there is a single test scheduler in a single "go test" function, then a
single reporter instance that serves it needs no locking. (Recommended)
needed.
* If there are multiple test schedulers, there should be a single reporeter
instance to seve them all and it needs to be locked.
* If there are multiple test schedulers and one reporter per scheduler, then
"go test" functions should not be run concurrently.
* Otherwise, test cases could be printed interwaved in console.

###Focus Mode
(TODO)

###Test Time
(TODO)
####Timeout
####Find slow tests

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


###Test Double
(TODO)

###Auto Test
(TODO)

Usage Guidelines
----------------
* Write tests that matter.
* Test simple abstractions rather than complex details.
    * Avoid mocking when possible.
* Write abstract specifications in text; write concrete examples in code.
* One expectation per test case.

Existing Go Testing Frameworks
------------------------------
###xUnit Style
* launchpad.net/gocheck
* github.com/stretchr/testify

###BDD Style
* smartystreets.github.io/goconvey
* github.com/onsi/ginkgo
* github.com/franela/goblin
* github.com/orfjackal/gospec
* github.com/stesla/gospecify
* github.com/azer/mao
* github.com/pranavraja/zen (forked from mao)

###Expectations (assertions)
* github.com/onsi/gomega
* launchpad.net/gocheck
* github.com/stretchr/testify/assert

###Mock
* code.google.com/p/gomock
* https://github.com/qur/withmock (gomock companion)
* github.com/stretchr/testify/mock
* github.com/jvshahid/mock4go

Reference
---------
http://betterspecs.org/
http://doctoc.herokuapp.com/
