// Package gocond provides some simple syntactic sugar for transacting conditional goroutines
//
// The idea of a conditional goroutine is that a function needs to be called and that function needs to run in the background.  Normally in such a case, you would execute the function inside a goroutine:
//
//    go myFunc()
//
// However, there exist a couple common problems here:
//   - There may be checks which must be performed which may preclude the goroutine from running
//   - There may be synchronous actions which must be performed before the surrounding code is allowed to continue execution
//
// To be sure, there are plenty of solutions to each of these problems.
//
// You could embed the goroutine inside the subroutine:
//
//    func Subroutine(a int) error {
//      if a < 1 {
//        return errors.New("A must not be less than 1")
//      }
//      go func() {
//        // do stuff
//      }()
//      return nil
//    }
//
//    func main() {
//      err := Subroutine(1)
//      // do stuff
//    }
//
// The problem here is that when reading through the top-level logic, it is
// unclear that Subroutine runs in the background.
//
// You could provide a channel for signaling:
//
//    func Subroutine(a int, errCh <-chan error) {
//      defer close(errCh)
//
//      if a < 1 {
//        errCh <-errors.New("A must not be less than 1")
//      }
//
//      go func() {
//        // do stuff
//      }()
//      return
//    }
//
//
// This latter path provides a fairly flexible system, with the cost of an
// additional temporary goroutine.  What this package does, essentially, is
// provide simple wrapper to do exactly this.
//
package gocond

import "context"

// Conditional is a function which returns an error.  Semantically, it is
// implied to fork in the background if and when its preconditions are
// satisfied.  Otherwise, it will return an error.
type Conditional func() error

// ConditionalCtx is a function which takes context and returns an error.  Semantically, it is implied that the ConditionalCtx will execute in the background if and when its internal preconditions are satisfied.  Otherwise, it returns an error.
type ConditionalCtx func(context.Context) error

// Run executes a function which implicitly runs in the background after checking and processing all of its preconditions.
func Run(f func() error) <-chan error {
	ctx := context.Background()
	return RunCtx(ctx, func(context.Context) error {
		return f()
	})
}

// RunCtx executes a backgrounding function using the supplied context
func RunCtx(ctx context.Context, f func(context.Context) error) <-chan error {
	ret := make(chan error)
	go func() {
		select {
		case <-ctx.Done():
			return
		case ret <- f(ctx):
			return
		}
	}()
	return ret
}
