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
* Natual:     Both BDD style and table driven style supported. Just use the one that fits your test scenario.
* Reliabile:  the design is minimal and orthogonal; the code is tested under 100% coverage.
* Separable:  the expectation (assertion) package can be used alone.
* Extensible: fully customizable expectations and test reporters.
* Compatible: "go test" is enough.

Design Documents
----------------

[Core](DESIGN.md)

[Expectations](expectation/DESIGN.md)

Usage
-----

    import (
        exp "github.com/hailiang/gspec/expect"
        "github.com/hailiang/gspec"
        "github.com/hailiang/gspec/suite"
        "testing"
    )

    var _ = suite.Add(func(s gspec.S) {
        describe, when, it := s.Alias("describe"), s.Alias("when"), s.Alias("it")
        expect := exp.Alias(s.Fail)

        describe("an integer", func() {
            i := 2
            when("it is incremented by 1", func() {
                i++
                it("should has a value of original value plus 1", func() {
                    expect(i).Equals(3)
                })
            })
            // more scenarios here.
        })

        // more tests here.
    })

    func TestAll(t *testing.T) {
        suite.Run(t, false)
    }
