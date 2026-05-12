package main

import (
	"fmt"
	"kvs"
	"runtime"
	"time"
)

const targetUnit = 1024 * 1024

var c = int64(0)
var l = int64(0)
var mem runtime.MemStats
var out interface{}

func main() {
	var now time.Time
	store := kvs.NewStore(targetUnit * 16)
	var ss [targetUnit * 16]string
	var err error
	var s string

	targetSize := targetUnit

	runtime.ReadMemStats(&mem)
	fmt.Println(
		"mem (GiB)", float64(mem.Alloc/1024/1024)/1024,
		float64(mem.TotalAlloc/1024/1024)/1024, float64(mem.Sys/1024/1024)/1024)

	now = time.Now()
	fmt.Println()
	fmt.Println("creating store members pre-allocated: ", targetSize)
	for i := 0; i < targetSize; i++ {
		s, err = store.SetRndKey(struct{}{})
		if err != nil {
			panic(err)
		}
		ss[i] = s
		c++
	}
	fmt.Println("avg-time", time.Since(now)/time.Duration(targetSize))
	runtime.ReadMemStats(&mem)
	fmt.Println(
		"mem (GiB)", float64(mem.Alloc/1024/1024)/1024,
		float64(mem.TotalAlloc/1024/1024)/1024, float64(mem.Sys/1024/1024)/1024)

	now = time.Now()
	fmt.Println()
	fmt.Println("reading store members: ", targetSize)
	for i := 0; i < targetSize; i++ {
		out, err = store.Get(ss[i])
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("avg-time", time.Since(now)/time.Duration(targetSize))
	runtime.ReadMemStats(&mem)
	fmt.Println(
		"mem (GiB)", float64(mem.Alloc/1024/1024)/1024,
		float64(mem.TotalAlloc/1024/1024)/1024, float64(mem.Sys/1024/1024)/1024)

	targetSize *= 16
	now = time.Now()
	fmt.Println()
	fmt.Println("increasing size: ", 16, "times")
	for i := 0; i < targetSize; i++ {
		s, err = store.SetRndKey(struct{}{})
		if err != nil {
			panic(err)
		}
		ss[i] = s
	}
	fmt.Println("avg-time", time.Since(now)/time.Duration(targetSize))
	runtime.ReadMemStats(&mem)
	fmt.Println(
		"mem (GiB)", float64(mem.Alloc/1024/1024)/1024,
		float64(mem.TotalAlloc/1024/1024)/1024, float64(mem.Sys/1024/1024)/1024)

	now = time.Now()
	fmt.Println()
	fmt.Println("reading store members: ", targetSize)
	for i := 0; i < targetSize; i++ {
		out, err = store.Get(ss[i])
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("avg-time", time.Since(now)/time.Duration(targetSize))
	runtime.ReadMemStats(&mem)
	fmt.Println(
		"mem (GiB)", float64(mem.Alloc/1024/1024)/1024,
		float64(mem.TotalAlloc/1024/1024)/1024, float64(mem.Sys/1024/1024)/1024)

	//sec++
	//fmt.Println("sec", sec, "rate: ", atomic.LoadInt64(&c), "total: ", atomic.LoadInt64(&l))
}
