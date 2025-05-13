package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

type config struct {
	problemFilePath    string
	shuffleProblems    bool
	timeLimitInSeconds int
}

type gameState struct {
	problems []problem
	score    int
	done     chan bool
}

func parseProblems(rawProblems [][]string) ([]problem, error) {
	var problems []problem
	for _, rawProblem := range rawProblems {
		if len(rawProblem) < 2 {
			return nil, fmt.Errorf("invalid question format: %v", rawProblem)
		}

		q := problem{
			question: strings.TrimSpace(rawProblem[0]),
			answer:   strings.TrimSpace(rawProblem[1]),
		}
		problems = append(problems, q)
	}
	return problems, nil
}

func loadProblems(config config) ([]problem, error) {
	fmt.Println("Loading problems from:", config.problemFilePath)

	file, err := os.Open(config.problemFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rawQuestions, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return nil, err
	}

	problem, err := parseProblems(rawQuestions)
	if err != nil {
		fmt.Println("Error parsing problems:", err)
		return nil, err
	}

	if !config.shuffleProblems {
		return problem, nil
	}

	rand.Shuffle(len(problem), func(i, j int) {
		problem[i], problem[j] = problem[j], problem[i]
	})

	return problem, nil
}

func askQuestion(questionNumber int, totalQuestions int, problem problem) (bool, error) {
	correctAnswer := problem.answer
	correctAnswerAsInt, err := strconv.Atoi(strings.TrimSpace(correctAnswer))
	if err != nil {
		fmt.Println("Error converting correct answer to int:", err)
		return false, err
	}

	fmt.Printf("Question %d of %d: %s\n", questionNumber, totalQuestions, problem.question)
	var answer string
	fmt.Scanln(&answer)
	answerAsInt, err := strconv.Atoi(answer)
	if err != nil {
		fmt.Println("Error converting answer to int:", err)
		return false, err
	}

	if answerAsInt == correctAnswerAsInt {
		fmt.Println("Correct!")
		return true, nil
	} else {
		fmt.Println("Incorrect! The correct answer is:", problem.answer)
		return false, nil
	}
}

func runQuiz(gameState gameState) {
	for i, question := range gameState.problems {
		result, err := askQuestion(i+1, len(gameState.problems), question)
		if err != nil {
			fmt.Println("Error asking question:", err)
			return
		}

		if result {
			gameState.score++
		}
	}

	gameState.done <- true
}

func main() {
	var config config
	flag.StringVar(&config.problemFilePath, "file", "./assets/default-questions.csv", "Path to the CSV file containing questions")
	flag.BoolVar(&config.shuffleProblems, "shuffle", false, "Shuffle questions before asking")
	flag.IntVar(&config.timeLimitInSeconds, "time", 30, "Time limit for each question in seconds (0 for no limit)")
	flag.Parse()

	println("Welcome to Math Quiz!")
	problems, err := loadProblems(config)
	if err != nil {
		fmt.Println("Error loading questions:", err)
		return
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if config.timeLimitInSeconds == 0 {
		ctx = context.Background()
	} else {
		fmt.Printf("You have %d seconds to answer each question.\n", config.timeLimitInSeconds)
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(config.timeLimitInSeconds)*time.Second)
		defer cancel()
	}

	gameState := gameState{
		problems: problems,
		score:    0,
		done:     make(chan bool),
	}

	go func() { runQuiz(gameState) }()

	select {
	case <-gameState.done:
		fmt.Printf("Quiz complete! You answered %d out of %d questions correctly.\n!", gameState.score, len(gameState.problems))
	case <-ctx.Done():
		fmt.Println("Time's up! You answered", gameState.score, "out of", gameState.score, len(gameState.problems), "questions correctly.")
		fmt.Println("Exiting the quiz.")
	}
}
