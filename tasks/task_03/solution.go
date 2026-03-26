package main

import "fmt"

import "errors"

func fizzBuzz(n int) (string, error) {
	if n < 0 {
		return "", errors.New("negative number")
	}
	if n%15 == 0 {
		return "FizzBuzz", nil
	}
	if n%3 == 0 {
		return "Fizz", nil
	}
	if n%5 == 0 {
		return "Buzz", nil
	}
	return fmt.Sprint(n), nil
}
