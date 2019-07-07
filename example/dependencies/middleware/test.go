package middleware

import "github.com/gobuffalo/buffalo"

// TestContextKey is the id used for distinguishing which
// instances to use for a given test.
const TestContextKey = "TestID"

// ForTesting inserts headers to help track tests.
// Note the header used here should not be tampered with
// outside of test scenarios.
func ForTesting(next buffalo.Handler) buffalo.Handler {
	return func(ctx buffalo.Context) error {
		testID := ctx.Request().Header.Get("Test-ID")
		if testID != "" {
			ctx.Set(TestContextKey, testID)
		}
		return next(ctx)
	}
}
