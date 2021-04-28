package main

import "sync"

// сюда писать код
func ExecutePipeline(jobs ... job) {
	var wg sync.WaitGroup
	defer wg.Wait()
	in := make(chan interface{}, 1)

	for _, job := range jobs {
		wg.Add(1)
		out := make(chan interface{} ,1)
		go func(in chan interface{}, out chan interface{}, job, *sync.WaitGroup){}(in, out, job, &wg)
	}


}


func SingleHash(){

}

func MultiHash(){

}

func CombineResults(){

}