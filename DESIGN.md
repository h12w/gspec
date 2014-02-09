Design of GSepc
===============

Introduction
------------
GSpec is a Go testing framework that makes it easy to organize and verify the
mind model of software.

Design goals:

* It should be natual to write readable and runnable specifications.
* It should be an extension rather than replacement to "go test".
* It should be reliable by robust design and 100% test coverage.
* It should be minimal and extensible.

Yet Another One?
----------------
Why writing yet another testing framework, given there are already many? (
http://code.google.com/p/go-wiki/wiki/Projects#Testing)

* Automatic test is important, especially for a single developer who lacks time
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

###Extend "go test"
GSpec should not break but extend existing features that "go test" provides:

* "go test" command is the only way to run tests.
* Besides concurency at the level of test functions, GSpec should support
  concurency at the level of each expectation.
* GSpec should be able to split tests into multiple test functions/files.
* GSpec should be able to print helpful and easily readable messages when a test
  case fails.

###Test Case
A test case needs running some code and verifying the result. Further, the code
run by a test case can be usually seperated into 4 steps:

* setup a test context
* perform an action
* verify the result
* teardown the test context (optional)

Some setup/teardown code might be shared between tests, so there are 4 cases:

* setup before each test
* teardown after each test
* setup run once before all tests
* teardown run once after all tests

"go test" meet the requirement above by providing a way to define a test (test
functions with specific signature). It is minimal yet flexible. It can support
concurency at test function level easily. However, "go test" does nothing to
help the developer with setup/teardown, devlopers have to figure out themselves.

Traditional xUnit style testing frameworks implement each step above as virtual
methods in a base class or interface. They do provide a complete solution,
including concurrency. However, They have to introduce unavoidable boilerplate
code of defining derived test classes, and nested test group is not straight
forward to implement with derived classes.

BDD style frameworks like RSpec use closures to provide an ad hoc and natual
way to specify test cases. Pros of this method include:

* A natual short test description rather than a function name
* Setup/teardown code can be nested to form multiple levels of test groups
* The tests form a readable and runnable specification

GSpec should try to follow this way, though it is not obvious on how to
implement concurrency at this stage (RSpec does not support concurrency itself).

###Nested Test Group
Most of the test cases do not share a test context, but when they do, the
context should be setup (teardown) before (after) each test to isolate test
cases from each other.

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

A limitation of lacking a dedicated setup method is: when the setup panicks,
GSpec can only give the panicking postion but cannot tell it is a leaf node or
it still has children. However, it is a tolerable tradeoff.

It is rare to setup (teardown) a test context before (after) all tests without
worrying about its state changing. In such rare cases, they can be handled
separatedly at the start (end) of the "go test" testing functions.

###Concurrency
Concurrency is a core feature of Go. "go test" supports concurency at the level
of test functions. With a testing framework with nested test group by closures,
it is expected to have dozens (or even hundreds) of test cases written
in one test function. On a quad-core CPU, you probably could just split test
cases into four test functions, but CPU could easily get many cores (hundreds or
even thousands of cores) in the foreseeable future, thus it requires one
goroutine per test case to make the most of the hardware.

When the test cases are organized in a tree of closures, there is no way to know
the whole structure without actually running it. The process has to be
exploratory, starting goroutines on the fly.

To support concurency, test cases should be completed isolated from each other.
No variables should be shared without careful synchronization.

It seems that the variables defined in closures are shared between test cases,
but actually they are fully isolated, because each test case runs from the root
level. Variables are allocated on call stack of its own gouroutine.

The critical point is the state variables related to test scheduling. There has
to be a way to guide the test along a path from the root down to a certain leaf.
The path has to be stored somewhere and should not be shared between test cases.
Thus, scheduling related variables have to be passed into the goroutine function
(RootFunc) as an argument (g of interface type G). e.g.

    scheduler.Start(func(g G) {
        g.Group(func() {
            g.Group(func() {
                // test case 1
            })
            g.Group(func() {
                // test case 2
            })
        })
    })

where the RootFunc is defined as:

    type RootFunc func(g G)

and variable g contains all the variables needed to control which test case to
run, and the scheduler makes sure each test case run once concurrently
(sequentially).

RESTRICTION:
* To support concurency, there must be one context variable of type G per
  goroutine. So the Group method cannot be simply defined as a global function,
  and a top-level function of type RootFunc is mandatary for writing tests.
  GSpec has to exchange some simplicty for concurency (This could be compensated
  by defining aliases for the Group method).
* For loop should not be defined in any places excpet the closure of leaf node.

###Test Gathering
"go test" gathers test functions with specific function/file naming conventions.
GSpec should also be able to gather tests across test functions/files.

Unlike the group context G, the scheduler can be shared by all goroutines as
long as carefully locked. So RootFuncs can be defined anywhere and are gathered
by a single scheduler in one "go test" function. This requires Scheduler.Start
accepts multiple RootFuncs.

RESTRICTION:
There is no finalization hook provided by "go test", so there is no way to know
when all test functions finish. The only possible way is to output the result
for each test gathering (Scheduler.Start). So RootFuncs should be better
gathered in a single "go test" function.

###Test Specification
GSpec should be able to generate a structured, readable plain text specification
from the tests written. There should be a way to define and collect information
for each level of test group.

####Alias
GSpec should be able to provide a convenient way to use customized alias names
for the nested group function. e.g. describe, context, it, specify and example.

GSpec should be able to tell what alias name is used for a certain test group.

####Description
GSpec should be able to assotiate a description to each level of test group.

Besides, it also should be able to allow types inserted between test description
so that during refactoring, these types can be found and modified too and won't
get out of sync. Go does not support first class type, so it could be
implemented by variables with zero value. (OPTIONAL)

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
customized reporter. The listener should call a reporter's method when needed.

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

###Failure
When test code panicks, "go test" prints the error and terminates immediately.
GSpec should respect this design.

When an expectation (assertion) fails, GSpec should record the error and
continue running other test cases. t.Fail should be called to notify "go test"
there are failures.

GSpec itself could have internal bugs and fail, when it happens, GSpec should
print the error message and terminate immediately. (fail fast)

###Timeout

###Matchers

###Table-driven Testing
There are 2 ways:
* Define the for-loop in leaf node
* Define the for-loop out side of RootFunc
  A for-loop arround a higher order function that accepts the test variable and
  returns a RootFunc.

###Focus Mode
####Support metadata for each test group?

####Filter by meta data including regular expressions

###Benchmark & coverage
What "go test" provides are good enough. Just don't break them.

###Options
Options of GSpec should able to set hard coded or via CLI flags. Flag should
have higher priority than hard coded value so that can be changed at runtime.

###Mock

###Auto Test

Usage Guidelines
----------------
* Write abstract specifications in text; write concrete examples in code.
* Use given-when-then pattern to describe the context, input and expected
  behavior of the software.
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

###Matcher
* github.com/onsi/gomega
* launchpad.net/gocheck
* github.com/stretchr/testify

###Mock
* code.google.com/p/gomock

Reference
---------
http://betterspecs.org/
