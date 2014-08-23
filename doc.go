// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gspec is an expressive, reliable, concurrent and extensible Go test framework
that makes it productive to organize and verify the mind model of software.

	- Expressive: a complete runnable specification can be organized via both BDD and table driven styles.
	- Reliable:   the implementation has minimal footprint and is tested with 100% coverage.
	- Concurrent: test cases can be executed concurrently or sequentially.
	- Extensible: customizable BDD cue words, expectations and test reporters.
	- Compatible: "go test" is sufficient but not mandatory to run GSpec tests.

GSpec is very modular and sub packages have minimal or no dependance on each
other. The top package "gspec" integrates all other sub packages and provide a
quick way of test gathering, executing and reporting.
*/
package gspec
