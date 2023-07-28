package module4

import (
	"fmt"
	"time"
)

type LightColor int64

const (
	Red LightColor = iota
	Green
	Yellow
)

func (c LightColor) String() string {
	switch c {
	case Red:
		return "red"
	case Green:
		return "green"
	case Yellow:
		return "yellow"
	}
	return "unknown"
}

func Lesson_2() {
	var state = make(chan string)
	var msg = map[string]int{
		"red":    15,
		"green":  30,
		"yellow": 3,
	}

	var light string
	var timer int

	go func() {
		for {
			state <- "red"
			time.Sleep(time.Second * 15)
			state <- "green"
			time.Sleep(time.Second * 30)
			state <- "yellow"
			time.Sleep(time.Second * 3)
		}
	}()

	for {
		light = <-state
		timer = msg[light]
		for {
			time.Sleep(time.Second * 1)
			fmt.Printf("Light: %s - Time: %d\n", light, timer)
			timer--

			if timer == 0 {
				break
			}
		}
	}

}
