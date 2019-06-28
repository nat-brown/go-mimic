# Example server tested with go-mimic

This is a simple server for testing go-mimic and demonstrating how to incorporate it into functional or server-wide integration testing.

In the `dependencies` folder, you'll find another example using a different framework. All "dependencies" are hosted on port `3000` as different routes on the same server to simplify things.

---
### Setup

You'll need a working Go environment. Installation instructions are [here](https://golang.org/doc/install) and workspace instructions can be found [here](https://golang.org/doc/code.html).

While go-mimic aims to only use the standard library, this demonstration has no such constraints. This project will use `dep` until acceptance of `Modules` is wider-spread for ease of understanding. You can install and read how to run it [here](https://github.com/golang/dep).

The API doesn't currently require any special databases to run, so the simple command

    go run main.go

from the example directory will start it up on port `8080`. 

But, we're all here to see some tests, so taking a look at the tests throughout the code is much more rewarding. You should also try pulling down the code and running your own tests to see how things look.

To run the tests, use the command

    go test

from a directory with a test file. To run tests against all directories including and branching from your current one, run

    go test ./...
