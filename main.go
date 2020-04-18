//flag package, csv package, os package, channels and go routines
//for timer and time package with timer

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	//setting up being able to input their own csv file and setting default to problems.csv
	csvFile := flag.String("csvFile", "problems.csv", "csv file with questions,answer")
	flag.Parse()

	//set up timer with default 30 seconds
	timeLimit := flag.Int("limit", 3, "time limit to answer quiz in seconds")
	flag.Parse()

	//opening the file and checking for errors
	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Printf("Error: Failed to open file %s", *csvFile)
		os.Exit(1)
	}

	//reading the file and checking for errors
	r := csv.NewReader(file)
	record, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	//making the file into a slice of structs
	problems := makeStruct(record)
	//fmt.Println(problems)

	//creating a timer after setup is ready
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	//displaying questions, adding to correct score if user answered correctly
	correct := 0
	for i, p := range problems {
		//automatically give each question after an answer
		fmt.Printf("#%d: %s = \n", i+1, p.q)

		//setup answer channel to know if they gave an answer
		answerCh := make(chan string)

		//setup anonymous function to have them answer so it doesn't block the timer.
		//if answer is given before time is up send it over the channel
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		//if time limit has reached automatically stop
		//if still in time log if the answer was correct
		select {
		case <-timer.C:
			fmt.Printf("\nQuiz Complete. You scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	//printing out final score if they completed before time was out
	fmt.Printf("Quiz Complete. You scored %d out of %d.\n", correct, len(problems))
}

type problem struct {
	q string
	a string
}

//making the file into a slice of structs
func makeStruct(record [][]string) []problem {
	probAns := make([]problem, len(record))
	for i, rec := range record {
		probAns[i] = problem{
			q: rec[0],
			a: strings.TrimSpace(rec[1]),
		}
	}
	return probAns
}
