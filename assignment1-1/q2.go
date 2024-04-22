package cos418_hw1_1

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`

	sum := 0
	for num := range nums {
		sum += num
	}
	fmt.Println(sum)
	out <- sum
	close(out)
}

type ch chan int

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	r, err := os.Open(fileName)
	checkError(err)
	defer r.Close()

	arr, err := readInts(r)
	checkError(err)

	channels := make([]ch, num)
	outIntermediate := make([]ch, num)
	BUFFER_SIZE := len(arr) + 1

	for i := 0; i < num; i++ {
		channels[i] = make(chan int, BUFFER_SIZE)
		outIntermediate[i] = make(chan int)
	}

	// Map phase
	for i := 0; i < len(arr); i++ {
		random := rand1toN(num)
		channels[random] <- arr[i] // Distribute value to n channels
	}

	for i := 0; i < num; i++ {
		close(channels[i]) // Need to close so range in sumWorker can end or range  will block forever
	}

	for i := 0; i < num; i++ {
		go sumWorker(channels[i], outIntermediate[i])
	}

	// Reduce phase
	total := 0
	for _, out := range outIntermediate {
		total += <-out
	}
	fmt.Println(total)

	return total
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}

func rand1toN(maxx int) int {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(maxx)
	return n
}
