package main

import (
	"context"
	"errors"
	"sync"
)

func ParallelMap[T any, R any](
	ctx context.Context,
	workers int,
	in []T,
	fn func(context.Context, T) (R, error),
) ([]R, error) {
	if len(in) == 0 {
		return nil, nil
	}
	
	if workers <= 0 {
		return nil, errors.New("number of workers cannot be less then 1")
	}
	
	inputChan := make(chan int, len(in))
	go func() {
		for i := range in {
			inputChan <- i
		}
		close(inputChan)
	}()

	ctx, cansel := context.WithCancel(ctx)
	defer cansel()
	
	var wg sync.WaitGroup
	var err error
	var once sync.Once

	wg.Add(workers)
	out := make([]R, len(in))

	for range workers {
		go func() {
			defer wg.Done()
			for index := range inputChan {
				select{
				case <-ctx.Done():
					once.Do(func() {
						err = ctx.Err()
					})
				default:
					var internalErr error
					out[index], internalErr = fn(ctx, in[index])
					if internalErr != nil {
						once.Do(
							func() {
								err = internalErr
							},
						)
						cansel()
					}
				}
			}
		}()
	}
	wg.Wait()
	return out, err
}
