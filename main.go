package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type config struct {
	questionFilePath string
	shuffleQuestions bool
}

func loadQuestions(config config) ([][]string, error) {
	fmt.Println("Loading questions from:", config.questionFilePath)

	file, err := os.Open(config.questionFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	questions, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return nil, err
	}

	if !config.shuffleQuestions {
		return questions, nil
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions, nil
}

func askQuestion(questionNumber int, totalQuestions int, question []string) (bool, error) {
	correctAnswer := question[1]
	correctAnswerAsInt, err := strconv.Atoi(strings.TrimSpace(correctAnswer))
	if err != nil {
		fmt.Println("Error converting correct answer to int:", err)
		return false, err
	}

	fmt.Printf("Question %d of %d: %s\n", questionNumber, totalQuestions, question[0])
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
		fmt.Println("Incorrect! The correct answer is:", question[1])
		return false, nil
	}
}

func main() {
	var config config
	flag.StringVar(&config.questionFilePath, "file", "./assets/default-questions.csv", "Path to the CSV file containing questions")
	flag.BoolVar(&config.shuffleQuestions, "shuffle", false, "Shuffle questions before asking")
	flag.Parse()

	println("Welcome to Math Quiz!")
	questions, err := loadQuestions(config)
	if err != nil {
		fmt.Println("Error loading questions:", err)
		return
	}

	numberOfQuestions := len(questions)
	totalScore := 0

	for i, question := range questions {
		result, err := askQuestion(i+1, numberOfQuestions, question)
		if err != nil {
			fmt.Println("Error asking question:", err)
			return
		}

		if result {
			totalScore++
		}
	}
	fmt.Printf("You answered %d out of %d questions correctly.\n", totalScore, numberOfQuestions)
}
