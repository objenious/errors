# errors for go 1.13 and beyond [![Travis-CI](https://travis-ci.org/objenious/errors.svg)](https://travis-ci.org/objenious/errors) [![AppVeyor](https://ci.appveyor.com/api/projects/status/b98mptawhudj53ep/branch/master?svg=true)](https://ci.appveyor.com/project/objenious/errors/branch/master) [![GoDoc](https://godoc.org/github.com/objenious/errors?status.svg)](http://godoc.org/github.com/objenious/errors) [![Report card](https://goreportcard.com/badge/github.com/objenious/errors)](https://goreportcard.com/report/github.com/objenious/errors) [![Sourcegraph](https://sourcegraph.com/github.com/objenious/errors/-/badge.svg)](https://sourcegraph.com/github.com/objenious/errors?badge)

Package errors provides simple error handling primitives, compatible with Go 1.13 error wrapping.

`go get github.com/objenious/errors`

The traditional error handling idiom in Go is roughly akin to
```go
if err != nil {
        return err
}
```
which applied recursively up the call stack results in error reports without context or debugging information. The errors package allows programmers to add context to the failure path in their code in a way that does not destroy the original value of the error.

## Adding context to an error

The errors.Wrap function returns a new error that adds context to the original error. For example
```go
_, err := ioutil.ReadAll(r)
if err != nil {
        return errors.Wrap(err, "read failed")
}
```

[Read the package documentation for more information](https://godoc.org/github.com/objenious/errors).


## License

BSD-2-Clause
