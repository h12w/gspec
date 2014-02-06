Design of GSepc
===============

Introduction
------------
GSpec is a Go testing framework that makes it easy to organize and verify the
mind model of software.

Design goals:

* It should provide a natual way to write readable and runnable specifications.
* It should be able to run tests concurrently so that the tests are expected to
  run faster and faster on future hardware.
* It should be minimal yet extensible. Core features are provided by composition
  of orthogonal parts and advanced features are provided by external extensions.
* It should be reliable itself by robust design and 100% test coverage.
* It should respect the design of builtin go-test tool so that there will not
  be any integration issues.

Yet Another One?
----------------
Why writing yet another testing framework, given there are already many? (
http://code.google.com/p/go-wiki/wiki/Projects#Testing)

* Automatic test is important, especially for a single developer who lacks time
  to test manually.
* go-test only meets minimal requirement, it leaves a gap to fill in a way or
  another.
* None of the existing frameworks fully satisfy all the goals above.
* It is a good way to learn testing itself. e.g. What to test? How much to test?
  How to test? How to name and orgnaize tests?
* Writing a testing framework is not trivial but not hard either, and it is fun.

Features
--------

The following sections are organized in the sequence of design decisions, from
major features to minor ones.

###Extend go-test
GSpec should extend rather than replace go-test. It means:
* go-test command is the only way to run tests.
* go-test command arguments are respected as much as possible.
* Tests are written in or called from test functions.

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

Besides tests themselves, there must be some mechanism to gather all the test
cases together and run them (test runner).

go-test meed the requirement above by providing a way to define a test (test
functions with specific signature) and a mechanism to gather tests (specific
function/file naming conventions). It is minimal yet flexible. It can support
concurency at test function level easily. However, go-test does nothing to
help the developer with setup/teardown, devlopers have to figure out themselves.

Traditional xUnit style testing frameworks implement each step above as virtual
methods in a base class or interface. They do provide a complete solution,
including concurrency. However, They have to introduce unavoidable boilerplate
code of defining derived test classes, and nested test group is not straight
forward to implement with derived classes.

BDD style frameworks like RSpec use closures to provide an ad hoc and natual
way to specify test cases. Pros of this method include:

* A natual short test description rather than a function name
* Setup/teardown code can be nested to form multiple levels of test group

GSpec will try to follow this way, though it is not obvious on how to implement
concurrency at this stage.

###Nested Test Group
Most of the test cases do not share a test context, but when they do, the
context should be setup (teardown) before (after) each test to isolate tests
from each other.

Thus, nested test group is a tree structure. Each leaf represents a test case
and the path from the root node to the leaf represents the setup sequence of the
test context for the test case. Reversely, the path from the leaf to the root
represents the teardown sequence of the test context for the test case. e.g.

    a
        b
            c1
            c2

The setup sequence will be: abc1 abc2, and the teardown sequence will be: c1ba
c2ba. (The relative order between c1 and c2 is not important)

One way of doing setup/teardown is hook function as what RSpec does (before,
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
there is no way to tell it is a leaf node or it still has children. It is a
tolerable tradeoff.

It is rare to setup (teardown) a test context before (after) all tests without
worrying about its state changing. In such rare cases, they can be handled
separatedly at the start (end) of the go-test testing functions.

###Concurrency
Concurrency is a core feature of Go. go-test supports concurency at the level
of test functions. With a testing framework with nested test group by
closures, it is expected to have dozens (or even hundreds) of test cases written
in one test function. On a quad-core CPU, you probably could just split test
cases into four test functions, but CPU could easily get many cores (hundreds or
even thousands of cores) in the foreseeable future, thus it requires one
goroutine per test case to make the most of the hardware.

When the test cases are organized in a tree of closures, there is no way to know
the whole structure without actually running it. So it is not possible to run
all test cases simultaneously, and the process has to be exploratory, starting
goroutines on the fly.

To support concurency, test cases should be completed isolated from each other.
No variables should be shared without careful synchronization.

It seems that the variables defined in closures are shared between test cases,
but actually they are fully isolated, because each test case runs from the root
level. Variables are allocated on call stack of its own gouroutine.

The critical point is the state variables related to test scheduling. There has
to be a way to guide the test along a path from the root down to a certain leaf.
The path has to be stored somewhare but the path should not be shared between
test cases. Thus, scheduling related variables have to be passed into the
goroutine function as an argument. e.g.

    runner.Run(func(g *G) {
        g.Group(func() {
            g.Group(func() {
                // test case 1
            })
            g.Group(func() {
                // test case 2
            })
        })
    })

where variable s of type S contains all the variables needed to control which
test case to run, and runner is the scheduler that makes sure each test case run
once concurrently (or sequentially).

###Test Specification
GSpec should be able to generate a structured, readable plain text specification
from the tests written. There should be a way to define and collect information
for each level of test group.

Note that a specification should be generated for each go-test function
separately, because there is no direct way to call "generate" after all go-test
functions return. (A timeout mechanism in a special test function may help, but
it just does not worth it)

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
A listener is an interface that collects the outputs from the tests, including
the structure of nested test groups, test descriptions and the results of test
running.

The implementation of a listener must assume being called concurrently out of
order. To reconstruct a tree structure, it will be easy if a parent node is
always collected before its children, so a listener should collect before the
start of each test group. To collect the result of each test case, a listener
also needs to collect at the end of each test group.

###Reporting
####Report progress
A progresser shows the progress of testing.

####Report final result
A reporter generates a complete report of all tests.

###Failure
GSpec should isolate each test case so that a fail on one test case does not
affect another, when any of the test case fails, t.Fail should be called.

GSpec should call t.FailNow when an internal error occurs, e.g. Formatter error.

###Timeout

###Matchers

###Focus Mode

###Benchmark

###Options
Options of GSpec should able to set hard coded or via CLI flags. Flag should
have higher priority than hard coded value so that can be changed at runtime.

###Mock

###Auto Test

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
