This package contain examples of writing tests with openfeature sdk for golang

Test method: ensure a happens-before relationship between goroutines using latches
to guarantee the following execution order:

1.  Test-one: setup Flag-A
2.  Test-two: setup Flag-B
3.  Test-two: execute code dependent of the Flag-B
4.  Test-one: execute code dependent of the Flag-A


###  Using `memprovider` package 
Tests in this package are expected to fail when executed in parallel as they rely on the global state


### Using `testing` package
Tests in this package are expected to pass when executed in parallel
as they rely on the TestProvider that utilises a custom _goroutine local_ storage
to manage state on a per-test basis.