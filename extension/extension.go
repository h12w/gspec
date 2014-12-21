// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package extension contains the types for extending the core package.
*/
package extension // import "h12.me/gspec/extension"

import (
	"time"
)

// Reporter is a interface to accept events from tests running.
type Reporter interface {
	Start()                          // Start should be called before all tests start.
	End(group *TestGroup)            // End should be called after all tests end.
	Progress(g *TestGroup, s *Stats) // Progress should be called whenever the statistics change.
}

// A TestGroup contains a test group's related data.
type TestGroup struct {
	ID          string        // ID of the test group, used for focus mode.
	Description string        // Description of the test group.
	Error       error         // Error passed in during execution of the test group.
	Duration    time.Duration // Time duration of running the test group, set only to leaf node (test case).
	Children    TestGroups    // Children nested within the test groups.
}

// TestGroups is the type of slice of *TestGroup
type TestGroups []*TestGroup

// Stats contains statistics of tests running.
type Stats struct {
	Total   int // Total number of test cases.
	Ended   int // Number of ended test cases.
	Failed  int // Number of failed test cases.
	Pending int // Number of pending test cases.
}
