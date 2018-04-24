# deque

[![Build Status](https://travis-ci.org/gammazero/deque.svg)](https://travis-ci.org/gammazero/deque)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/deque)](https://goreportcard.com/report/github.com/gammazero/deque)
[![codecov](https://codecov.io/gh/gammazero/deque/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/deque)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/gammazero/deque/blob/master/LICENSE)

A fast ring-buffer deque (double-ended queue) implementation.

[![GoDoc](https://godoc.org/github.com/gammazero/deque?status.png)](https://godoc.org/github.com/gammazero/deque)

This deque implementation automatically re-sizes by powers of two, growing when additional capacity is needed and shrinking when only a quarter of the capacity is used.  This allows bitwise arithmetic for all calculations.

The ring-buffer implementation significantly improves memory and time performance with fewer GC pauses, compared to implementations based on slices and linked lists.

For maximum speed, this deque implementation leaves concurrency safety up to the application to provide, however is best for the application if needed at all.

## Installation

```
$ go get github.com/gammazero/deque
```

## Example

```go
package main

import (
    "fmt"
    "github.com/gammazero/deque"
)

func main() {
    var q deque.Deque
    q.PushBack("foo")
    q.PushBack("bar")
    q.PushBack("baz")

    fmt.Println(q.Len())   // Prints: 3
    fmt.Println(q.Front()) // Prints: foo
    fmt.Println(q.Back())  // Prints: baz

    q.PopFront() // remove "foo"
    q.PopBack()  // remove "baz"

    q.PushFront("hello")
    q.PushBack("world")

    // Print: hello bar world
    for i := 0; i < q.Len(); i++ {
        fmt.Print(q.PeekAt(i), " ")
    }
    fmt.Println()
}
```
