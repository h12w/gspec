GSpec
=====

[![Build Status](https://travis-ci.org/hailiang/gspec.png?branch=master)](https://travis-ci.org/hailiang/gspec)
[![Coverage Status](https://coveralls.io/repos/hailiang/gspec/badge.png?branch=master)](https://coveralls.io/r/hailiang/gspec?branch=master)
[![GoDoc](https://godoc.org/github.com/hailiang/gspec?status.png)](https://godoc.org/github.com/hailiang/gspec)

GSpec is a concurrent, minimal, extensible and reliable testing framework in Go
that makes it easy to organize and verify the mind model of software.

(under development).

Goals:

* It should be natual to write readable and runnable specifications.
* It should be an enhancement rather than replacement to "go test".
* It should be reliable by robust design and 100% test coverage.
* It should be minimal and extensible.

Design Documents
----------------

[Core](DESIGN.md)

[Expectations](expectation/DESIGN.md)

Usage
-----

    import (
        "github.com/hailiang/gspec"
        "github.com/hailiang/gspec/suite"
        exp "github.com/hailiang/gspec/expect"
    )

    var _ = suite.Add(func(s gspec.S) {
        describe, when, it := s.Alias("describe"), s.Alias("when"), s.Alias("it")
        expect := exp.Alias()

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
