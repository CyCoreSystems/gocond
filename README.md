# gocond [![](https://godoc.org/github.com/CyCoreSystems/gocond?status.svg)](http://godoc.org/github.com/CyCoreSystems/gocond)

Simple syntactic sugar for transacting conditional goroutines

The idea of a conditional goroutine is that a function needs to be called and that function needs to run in the background.  Normally in such a case, you would execute the function inside a goroutine:

```go
   go myFunc()
```

However, there exist a couple common problems here:
   - There may be checks which must be performed which may preclude the goroutine from running
   - There may be synchronous actions which must be performed before the surrounding code is allowed to continue execution

To be sure, there are plenty of solutions to each of these problems.

You could embed the goroutine inside the subroutine:

```go
   func Subroutine(a int) error {
     if a < 1 {
       return errors.New("A must not be less than 1")
     }
     go func() {
       // do stuff
     }()
     return nil
   }

   func main() {
     err := Subroutine(1)
     // do stuff
   }
```

The problem here is that when reading through the top-level logic, it is
unclear that Subroutine runs in the background.

You could provide a channel for signaling:

```go
   func Subroutine(a int, errCh <-chan error) {
     defer close(errCh)
 
     if a < 1 {
       errCh <-errors.New("A must not be less than 1")
     }
 
     go func() {
       // do stuff
     }()
     return
   }
```

This latter path provides a fairly flexible system, with the cost of an
additional temporary goroutine.  What this package does, essentially, is
provide simple wrapper to do exactly this.
