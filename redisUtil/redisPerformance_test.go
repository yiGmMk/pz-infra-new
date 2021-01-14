package redisUtil

import (
	"fmt"
)

var (
	success, failed int
)

type MyObj struct {
	A, B string
}

func worker(start chan bool, index int, key string, ch chan int) {
	<-start

	err := SetObject(key, &MyObj{"123", "456"})
	if err != nil {
		fmt.Println("This is Worker:", index, ", Error", err)
		ch <- 0
		return
	}

	var obj MyObj
	err = GetObject(key, &obj)
	if err != nil {
		fmt.Println("This is Worker:", index, ", Error", err)
		ch <- 0
		return
	}
	fmt.Println("This is Worker:", index, ", Success")
	ch <- 1
}

func main() {
	start := make(chan bool)
	size := 150
	ch := make(chan int, size)
	for i := 1; i <= size; i++ {
		go worker(start, i, string(i), ch)
	}
	close(start)

	recieved := 0
	for {
		v := <-ch
		recieved++
		if v == 0 {
			failed++
		} else {
			success++
		}

		if recieved >= size {
			break
		}
	}
	fmt.Println("success count : ", success, " failed count: ", failed)
}
