/*
Package race_and_sync contains some demo functions of data races scenarios and the solutions for that.

A data race occurs when two goroutines access the same variable concurrently and at least one
of the accesses is a write.

It is necessary to add the -race flag to the go command when testing, building, running and installing.
For detail, see https://golang.org/doc/articles/race_detector.html.
 */
package race_and_sync
