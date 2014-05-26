// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"time"
)

// Reporter is a interface to accept events from tests running.
type Reporter interface {
	Start()
	End(groups TestGroups)
	Progress(g *TestGroup, s *Stats)
}

// A TestGroup contains a test group's related data.
type TestGroup struct {
	ID          string
	Description string
	Error       error
	Duration    time.Duration
	Children    TestGroups
}

// TestGroups is the type of slice of *TestGroup
type TestGroups []*TestGroup

// ByID implements Less method of sort.Interface for sorting TestGroups by ID.
type ByID struct{ TestGroups }

// Stats contains statistics of tests running.
type Stats struct {
	Total   int
	Ended   int
	Failed  int
	Pending int
}
