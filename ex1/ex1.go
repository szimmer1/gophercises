package ex1

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const DEFAULT_FILENAME string = "problems.csv"
const DEFAULT_LIMIT int = 30

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Runner struct {
	Total    int
	Correct  int
	Scanner  *bufio.Scanner
	Reader   *bufio.Reader
	doneChan chan bool
}

func NewRunner(scanner *bufio.Scanner) Runner {
	stdioreader := bufio.NewReader(os.Stdin)
	return Runner{
		Scanner:  scanner,
		Reader:   stdioreader,
		doneChan: make(chan bool),
	}
}

func (r *Runner) DoLine(s string) {
	// increment problem count
	r.Total += 1

	// parse this line
	split := strings.Split(s, ",")
	expressionStr, sumStr := split[0], split[1]

	// ask & block for input
	fmt.Printf("Problem #%d: %s = ", r.Total, expressionStr)
	input, err := r.Reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Trim(input, " \n\t")

	// determine if answer correct
	if sumStr == input {
		r.Correct += 1
	}
}

func (r *Runner) RunRand() {
	S := make([]string, 0)
	for r.Scanner.Scan() {
		// flip a coin to stack or unstack/print
		if rand.Intn(2) == 1 {
			// stack
			S = append(S, r.Scanner.Text())
		} else {
			// print
			r.DoLine(r.Scanner.Text())
		}
	}
	if err := r.Scanner.Err(); err != nil {
		panic(err)
	}
	// pop off the stack until empty
	for i := 0; i < len(S); i++ {
		r.DoLine(S[len(S)-i-1])
	}
	r.doneChan <- true
}

func (r *Runner) Run() {
	for r.Scanner.Scan() {
		r.DoLine(r.Scanner.Text())
	}
	if err := r.Scanner.Err(); err != nil {
		panic(err)
	}
	r.doneChan <- true
}

func Main() {
	// --csv <filename>, defaults to problems.csv
	filename := flag.String("csv", DEFAULT_FILENAME, "Filename for the problem set CSV")
	// --limit <time-limit>, default to 30s
	limit := flag.Int("limit", DEFAULT_LIMIT, "Time limit")
	flag.Parse()

	// handle file opening/closing
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// start runner async
	runner := NewRunner(bufio.NewScanner(f))
	go runner.RunRand()

	// start timer
	timer := time.NewTimer(time.Second * time.Duration(*limit))

	select {
	case <-runner.doneChan:
		// print results
		fmt.Printf("You scored %d out of %d.\n", runner.Correct, runner.Total)
	case <-timer.C:
		fmt.Println("Time limit exceeded.")
	}
}
