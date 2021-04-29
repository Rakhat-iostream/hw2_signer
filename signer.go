package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {

	var wg sync.WaitGroup
	defer wg.Wait()
	in := make(chan interface{})

	for _, work := range jobs {
		wg.Add(1)

		out := make(chan interface{})
		go func(job job, in chan interface{}, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			job(in, out)
		}(work, in, out, &wg)
		in = out
	}
}

func SingleHash(in, out chan interface{}) {

	var wg sync.WaitGroup
	var mx sync.Mutex

	for data := range in {
		wg.Add(1)
		go singleHashWorker(data.(int), out, &wg, &mx)
	}
	wg.Wait()
}

func singleHashWorker(in int, out chan interface{}, wg *sync.WaitGroup, mx *sync.Mutex) {

	defer wg.Done()
	str := strconv.Itoa(in)

	mx.Lock()
	md5Data := DataSignerMd5(str)
	mx.Unlock()

	crc32Chan := make(chan string)
	go func(data string, out chan string) {
		out <- DataSignerCrc32(data)
	}(str, crc32Chan)

	crc32Md5Data := DataSignerCrc32(md5Data)
	crc32Data := <-crc32Chan

	out <- crc32Data + "~" + crc32Md5Data
}

func MultiHash(in, out chan interface{}) {

	const th int = 6
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)
		go multiHashWorker(data.(string), out, th, wg)
	}
	wg.Wait()
}

func multiHashWorker(in string, out chan interface{}, th int, wg *sync.WaitGroup) {
	defer wg.Done()

	var mx sync.Mutex
	var Wg sync.WaitGroup

	concatArr := make([]string, th)

	for i := 0; i < th; i++ {
		Wg.Add(1)
		data := strconv.Itoa(i) + in

		go func(hash []string, index int, data string, Wg *sync.WaitGroup, mx *sync.Mutex) {
			defer Wg.Done()
			data = DataSignerCrc32(data)

			mx.Lock()
			hash[index] = data
			mx.Unlock()
		}(concatArr, i, data, &Wg, &mx)
	}

	Wg.Wait()
	out <- strings.Join(concatArr, "")
}

func CombineResults(in, out chan interface{}) {
	var result []string

	for i := range in {
		result = append(result, i.(string))
	}

	sort.Strings(result)
	out <- strings.Join(result, "_")
}
