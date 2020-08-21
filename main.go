package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
var m = sync.RWMutex{}

func main() {

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		wg.Add(2)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			if b, ok := queryCache(id); ok {
				fmt.Printf("from cache %v\n", b.ID)
				//fmt.Println(b)
			}
		}(id, wg)
		go func(id int, wg *sync.WaitGroup) {
			defer wg.Done()
			if b, ok := queryDatabase(id); ok {
				fmt.Printf("from database %v\n", b.ID)
				//fmt.Println(b)
			}
		}(id, wg)
		//time.Sleep(150 * time.Millisecond)
	}

	wg.Wait()
	//time.Sleep(2 * time.Second)
}

func queryCache(id int) (Book, bool) {
	m.RLock()
	b, ok := cache[id]
	m.RUnlock()
	fmt.Println(ok, id)
	return b, ok
}

func queryDatabase(id int) (Book, bool) {
	time.Sleep(150 * time.Millisecond)
	for _, b := range books {
		if b.ID == id {
			m.Lock()
			cache[id] = b
			m.Unlock()
			return b, true
		}
	}

	return Book{}, false
}
