// Copyright 2011 The Golint Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"sync"
)

// Represents a line-based linter (which may or may not hold state).
//
// As a general rule, LineLinter should only be used for linters which do /not/
// need to hold state - those linters that must hold state would generally
// benefit from i.e. parsing capabilities. Additionally, LineLinters are not
// reset at the beginning of every file, so a LineLinter would have no way to
// tell where one file begins and another ends.
type LineLinter interface {
	Linter
	RunLint(chan Line, chan Lint, *sync.WaitGroup)
}

// A simple, generic line linter. Any other LineLinter should be theoretically
// equivalent to a specialization of this, although it may be desirable not to
// use this as a base (for readability and other reasons).
type SimpleLineLinter struct {
	LinterDesc
	lintFunc func(string) (bool, string)
}

func (sl SimpleLineLinter) Desc() LinterDesc {
	return sl.LinterDesc
}

func (sl SimpleLineLinter) RunLint(
		text chan Line,
		lints chan Lint,
		wg *sync.WaitGroup) {
	wg.Add(1)
	for line := range text {
		if bad, issue := sl.lintFunc(line.line); bad {
			lints <- LineLint{sl, line.Location, issue}
		}
	}
	wg.Done()
}

// A line-based linter using regular expressions.
//
// The linter returns one piece of lint for every line matching the regular
// expression - multiple matches in a single line are ignored.
type RegexLinter struct {
	LinterDesc
	// the regular expression to check for
	Regex string
}

func (rl RegexLinter) Desc() LinterDesc {
	return rl.LinterDesc
}

func (rl RegexLinter) String() string {
	return fmt.Sprintf("%s (%s)", rl.LinterDesc.String(), rl.Regex)
}

func (rl RegexLinter) RunLint(
		text chan Line,
		lints chan Lint,
		wg *sync.WaitGroup) {
	wg.Add(1)
	for line := range text {
		if matches, _ := regexp.Match(rl.Regex, []byte(line.line))
				matches {
			lints <- LineLint{rl, line.Location, ""}
		}
	}
	wg.Done()
}

type LineLint struct {
	linter Linter
	Location
	issue  string
}

func (lint LineLint) String() string {
	if len(lint.issue) == 0 {
		return fmt.Sprintf("%s at %s",
			lint.linter.String(), lint.Location)
	}
	return fmt.Sprintf("%s at %s: %s",
		lint.linter.String(), lint.Location, lint.issue)
}

