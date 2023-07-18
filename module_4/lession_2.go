package module4

import (
	"log"
	"sync"
	"time"
)

func Lesson_2() {
	var wg sync.WaitGroup
	var msg = map[string]int{
		"Red":    15,
		"Green":  30,
		"Yellow": 3,
	}

	for key, val := range msg {
		wg.Add(1)
		go trafficLight(key, val, &wg)
	}
	wg.Wait()

}

func trafficLight(msg string, second int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(second) * time.Second)
	var times = 0
	for {
		log.Printf("Traffic Light Color: %s\n", msg)
		time.Sleep(1 * time.Second)
		times++

		if times == second {
			break
		}
	}

}
