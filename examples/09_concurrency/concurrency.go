package concurrencydemo

import (
	"sync"
	"time"
)

func SquareAll(nums []int) []int {
	results := make([]int, len(nums))

	var wg sync.WaitGroup
	for index, value := range nums {
		wg.Add(1)
		go func(index, value int) {
			defer wg.Done()
			results[index] = value * value
		}(index, value)
	}

	wg.Wait()
	return results
}

func RaceMessages(delays map[string]time.Duration) string {
	results := make(chan string, len(delays))

	for label, delay := range delays {
		go func(label string, delay time.Duration) {
			time.Sleep(delay)
			results <- label
		}(label, delay)
	}

	return <-results
}

func WorkerPool(workers int, inputs []int) []int {
	if workers < 1 {
		workers = 1
	}

	type job struct {
		index int
		value int
	}

	type result struct {
		index int
		value int
	}

	jobs := make(chan job)
	results := make(chan result, len(inputs))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- result{
					index: job.index,
					value: job.value * job.value,
				}
			}
		}()
	}

	go func() {
		for index, value := range inputs {
			jobs <- job{index: index, value: value}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	output := make([]int, len(inputs))
	for result := range results {
		output[result.index] = result.value
	}

	return output
}
