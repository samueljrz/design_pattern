package main

import (
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type piFunc func(int) float64

func decoratorLog(fun piFunc, logger *log.Logger) piFunc {
	return func(n int) float64 {
		fn := func(n int) (result float64) {
			defer func(t time.Time) {
				logger.Printf("took=%v, n=%v, result=%v", time.Since(t), n, result)
			}(time.Now())

			return fun(n)
		}
		return fn(n)
	}
}

func decoratorCache(fun piFunc, cache *sync.Map) piFunc {
	return func(n int) float64 {
		fn := func(n int) float64 {
			res, ok := cache.Load(n)
			if ok {
				return res.(float64)
			}
			res = fun(n)
			cache.Store(n, res)
			return res.(float64)
		}
		return fn(n)
	}
}

func Pi(n int) float64 {
	ch := make(chan float64)

	for i := 0; i < n; i++ {
		go func(ch chan float64, k float64) {
			ch <- 4 * math.Pow(-1, k) / (2*k + 1)
		}(ch, float64(i))
	}

	x := 0.0
	for i := 0; i < n; i++ {
		x += <-ch
	}
	return x
}

func main() {
	funCache := decoratorCache(Pi, &sync.Map{})
	fun := decoratorLog(funCache, log.New(os.Stdout, "test ", 1))

	fun(100000)
	fun(20000)
	fun(20000)

	// funCache(1000)
	// funCache(1000)
	// funCache(1000000)
	// funCache(1000000)

	//fmt.Println(Pi(1000))
	//fmt.Println(Pi(1000000))
}
