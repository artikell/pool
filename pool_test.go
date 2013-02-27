/* Copyright (c) 2013, Stefan Talpalaru <stefan.talpalaru@od-eon.com>, Odeon Consulting Group Pte Ltd <od-eon.com>
 * All rights reserved. */

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package pool

import (
	"math"
	"runtime"
	"testing"
)

func work(args ...interface{}) interface{} {
	x := args[0].(float64)
	j := 0.
	for i := 1.0; i < 10000000; i++ {
		j += math.Sqrt(i)
	}
	return x*x + j
}

func processResults(t *testing.T, results []*Job) (sum float64) {
	for _, job := range results {
		if job.Result == nil {
			t.Error("got error:", job.Err)
		} else {
			sum += job.Result.(float64)
		}
	}
	return
}

func TestCorrectness(t *testing.T) {
	num_jobs := float64(50)
	runtime.GOMAXPROCS(5) // number of OS threads

	// without the pool
	reference := float64(0)
	for i := float64(0); i < num_jobs; i++ {
		reference += work(i).(float64)
	}

	// 1 worker, add before running
	mypool := NewPool(1)
	for i := float64(0); i < num_jobs; i++ {
		mypool.Add(work, i)
	}
	mypool.Run()
	mypool.Wait()
	if processResults(t, mypool.Results()) != reference {
		t.Error("1 worker, add before running")
	}

	// 1 worker, run before adding
	mypool = NewPool(1)
	mypool.Run()
	for i := float64(0); i < num_jobs; i++ {
		mypool.Add(work, i)
	}
	mypool.Wait()
	if processResults(t, mypool.Results()) != reference {
		t.Error("1 worker, run before adding")
	}

	// 10 workers, add before running
	mypool = NewPool(10)
	for i := float64(0); i < num_jobs; i++ {
		mypool.Add(work, i)
	}
	mypool.Run()
	mypool.Wait()
	if processResults(t, mypool.Results()) != reference {
		t.Error("10 workers, add before running")
	}

	// 10 workers, run before adding
	mypool = NewPool(10)
	mypool.Run()
	for i := float64(0); i < num_jobs; i++ {
		mypool.Add(work, i)
	}
	mypool.Wait()
	if processResults(t, mypool.Results()) != reference {
		t.Error("10 workers, run before adding")
	}
}