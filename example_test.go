package gocond

import (
	"errors"
	"log"
	"time"
)

func myMaybeGoroutine(a int) error {
	if a < 1 {
		return errors.New("a cannot be less than 1")
	}
	go func() {
		// do stuff
		time.Sleep(5)
	}()
	return nil
}

func ExampleRun() {
	a := 2

	// Explicit local example
	err := <-Run(func() error {
		if a < 1 {
			return errors.New("a cannot be less than 1")
		}

		go func() {
			// do stuff
		}()

		return nil
	})
	if err != nil {
		log.Fatalln("Function 1 failed to run")
	}

	// Compose a simple wrapper on any other function
	err = <-Run(func() error {
		return myMaybeGoroutine(a)
	})
	if err != nil {
		log.Fatalln("Function 2 failed to run")
	}
}
