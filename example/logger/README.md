# Logging Tools

This package contains the logger for the API as well as an example of how to incorporate one better catered to api testing. Note that while it does not use anything from `mimic`, it does contain helpful test tools.

The example api uses [Logrus](https://github.com/sirupsen/logrus), chosen for its prevalence in the community and the fact that it supports hooks. Hooks are important for solving an issue that occurs when testing functionality that logs information while it runs: Upon failure, all logs will print regardless of which test they involved. This is especially confusing when tests are run in parallel, as the logs will interweave.

Imagine you have a package with 20 tests for logic that logs 5 times. On your next build, you find that you have a regression and one of your tests have failed! You now have to decipher which 5 of the 100 logs matter in debugging your code.

### Hooking the logger for test clarity

`test.go` contains setup for a logger that will reroute all its logs to a testing object. This will link the logs to a single test and print them upon failure. It also separates them from any logs from tests running in parallel.

There are some quirks to the implementation based on how test logging and logrus work. Consider your logger if you decide to implement something similar.

#### Logs are added to a list

In `Fire`, nothing is actually logged, but rather appended to a list. This is because `t.Log` will prepend the line that called it-- helpful when called from the test function; not helpful when called from the logging file. It adds enough clutter that this implementation deemed it better to group the logs at the cost of having to call print in the test. There are tricks to keep this from being painful (see test files for an example).

#### No logging level is ignored

This implementation hooks all log levels and does not allow customization. This prevents erroneous logs from being ignored. If we don't even want them when debugging, do we want them at all? However, you may want to allow more customization depending on your circumstances by adding an attribute to your hook struct or something similar.

#### The original logger definition was extracted

The test logger implementation wraps `newInnerLogger`. This will allow any defaults to printing styles to make their way into the injected logger at the expense of a trivial function extraction rather than leaving the definition in `New`.

#### Traditional logging was silenced through trickery

Logrus won't trigger a hook if the logger wouldn't be triggered first, so we can't silence normal printing and still print for the test. If we don't silence logging, however, we get double logs and plenty of clutter. `blackHoleWriter` allows us to enable logging at all levels but not worry about extra prints.