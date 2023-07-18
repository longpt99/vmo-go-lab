package module4

import (
	"fmt"
	"net"
	"sync"
)

// Implement a port scanner that checks the status of multiple ports (1-65535)
// on a given host concurrently using goroutines.
func Lesson_3() {
	host := "localhost"

	var wg sync.WaitGroup

	// Loop through all the ports and scan them concurrently
	for port := 1; port <= 65535; port++ {
		// Add a new goroutine to the wait group
		wg.Add(1)

		// Scan the port in a separate goroutine
		go func(p int) {
			// Attempt to connect to the port
			address := fmt.Sprintf("%s:%d", host, p)
			conn, err := net.Dial("tcp", address)

			// If the connection is successful, print a message and close the connection
			if err == nil {
				fmt.Printf("Port %d is open\n", p)
				conn.Close()
			}

			// Signal the wait group that the goroutine has finished
			wg.Done()
		}(port)
	}

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}
