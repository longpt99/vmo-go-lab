package module4

import (
	"log"
	"sync"
	"time"
)

func Lesson_1() {
	var wg sync.WaitGroup
	var msg = []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5"}

	for i, val := range msg {
		wg.Add(1)
		go genMessage(val, 5-i, &wg)
	}
	wg.Wait()

}

func genMessage(msg string, second int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(second) * time.Second)
	log.Printf("Message: %s\n", msg)
	log.Printf("Time to sleep: %d\n", second)
}
