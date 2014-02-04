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
* It should be reliable itself by careful design and 100% test coverage.
* It should respect the design of builtin "go test" tool so that there will not
  be any integration issues.

Yet Another One?
----------------
Why writing yet another testing framework, given there are already many? (
http://code.google.com/p/go-wiki/wiki/Projects#Testing)

Here is my thought:
* To write software with high quality in both design and implementation, a good
  process of design and test is important, especially for a project with only
  a single developer. A single developer will always be short of time, and it
  means there will not enough time to do manual testing. Without enough test,
  the developer will hesitate to do any improvement to avoid breaking existing
  functionalities.
* "go test" only meets minimal requirement, and left the rest to developers
  themselves.
* When I check all the other testing frameworks, I can barely find any projects
  actually using them, even among the repositories of the framework's authors.
  Since there is not a single mature and dominant testing framework, why not
  using a testing framework written by myself?
* I also checked the source code of these frameworks. It seems that a testing
  framework is not trivial but not hard to implement either.
* It is a good way to learn testing itself.
  e.g. What to test? How much to test? How to test? How to name and orgnaize
  tests?
* It has some challenges and it is fun.

The following sections are organized in the sequence of design decisions, from
major to minor ones.

Extend "go test"
----------------
GSpec should extend rather than replace "go test". It means:
* "go test" command is the only way to run tests.
* "go test" command arguments are respected as much as possible.
* Tests are written in or called from test functions.

Basics
------
A test needs running some code and verifying the result. Further, the code run
by a test can be usually seperated into 4 steps:

* setup a testing context
* perform an action
* verify the result
* teardown the testing context (optional)

Some setup/teardown code might be shared between tests, so there are 4 cases:

* setup before each test
* teardown after each test
* setup run once before all tests
* teardown run once after all tests

Besides tests themselves, there must be some mechanism to gather all the tests
together and run them (test runner).

"go test" meed the requirement above by providing a way to define a test (test
functions with specific signature) and a mechanism to gather tests (specific
function/file naming conventions). It is minimal yet flexible. It can support
concurency at test function level easily. However, "go test" does nothing to
help the developer with setup/teardown, devlopers have to figure out themselves.

Traditional xUnit style testing frameworks implement each step above as virtual
methods in a base class or interface. They do provide a complete solution,
including concurrency. However, They have to introduce unavoidable boilerplate
code of defining derived test classes, and nested testing context is not
straight forward to implement with derived classes.

BDD style frameworks like RSpec use closures to provide an ad hoc and natual
way to specify a test. Pros of this method include:

* A natual short test description rather than a function name
* Setup/teardown code can be nested to form multiple levels of testing contexts

GSpec will try to follow this way, though it is not obvious on how to implement
concurrency at this stage.

Nested Testing Context
----------------------
Most of the tests do not share a testing context, but when they do, the context
should be setup (teardown) before (after) each test to isolate tests from each
other.

Thus, nested testing context is a tree structure. Each leaf represents a test
case and the path from the root node to the leaf represents the setup sequence
of the testing context for the test case. e.g.

    a
        b
            c1
            c2

The execution sequence will be: abc1 abc2.

One way of doing setup is hook function as what RSpec does (before each:, after
:each). An alternative way is to define the setup code directly in the closure
and schedule it to run before each test case. As long as the scheduling logic
does not introduce too much overhead and concurrency can be supported, the
latter should be preferred because of its simplicity. Also, the teardown could
be supported simply by a defer statement. e.g.

    do(func() {
        // setup
        defer func() {
            // teardown
        }()
        do(func() {
            // test case 1
        })
        do(func() {
            // test case 2
        })
    })

It is rare to setup (teardown) a testing context before (after) all tests
without worrying about its state changing. In such rare cases, they can be
handled separatedly at the start (end) of the "go test" testing functions.

Concurrency
-----------
Concurrency is a core feature of Go. "go test" supports concurency at the level
of test functions. With a testing framework with nested testing context by
closures, it is expected to have dozens (or even hundreds) of test cases written
in one test function. On a quad-core CPU, you probably could just split test
cases into four test functions, but CPU could easily get many cores (hundreds or
even thousands of cores) in the foreseeable future, thus it requires one
goroutine per test case to make the most of the hardware.

When the test cases are organized in a tree of closures, there is no way to know
the whole structure without actually running it. So it is not possible to run
all test cases simultaneously, and the process has to be exploratory, start
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

    Run(func(s S) {
        s.do(func() {
            s.do(func() {
                // test case 1
            })
            s.do(func() {
                // test case 2
            })
        })
    })

where variable s of type S contains all the variables needed to control which
test case to run, and function Run is the scheduler that makes sure each test
case run once concurrently.

Panicking
---------

Timeout
-------

Test Description
----------------
"go test": Function test name
string short test description

Matchers
--------

Auto Test
---------

###Nested Testing Context

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
###Mock
* code.google.com/p/gomock/

Reference
---------
http://betterspecs.org/
