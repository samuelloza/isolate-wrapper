package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/samuelloza/isolate-wrapper/src/application/services"
	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/http_request"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/rabbitmq"
	"github.com/samuelloza/isolate-wrapper/src/infrastructure/testcaseprovider"
)

func main() {
	log.Println("Starting isolate-wrapper listener...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
		log.Println("Using default RABBITMQ_URL:", rabbitURL)
	}

	testCasesPath := os.Getenv("TESTCASES_PATH")
	if testCasesPath == "" {
		testCasesPath = "/home/sam/project/github/isolate-wrapper/test/testcases"
		log.Println("Using default TESTCASES_PATH:", testCasesPath)
	}

	clientResult := os.Getenv("CLIENT_RESULT")
	if clientResult == "" {
		clientResult = "http://localhost:8085/submit"
		log.Println("Using default CLIENT_RESULT:", clientResult)
	}

	numParallelProcess, err := strconv.Atoi(os.Getenv("PARALLEL_PROCESS"))

	if err != nil {
		numParallelProcess = 2
		fmt.Println("Error converting string to integer:", err)
	}

	if numParallelProcess == 0 {
		numParallelProcess = 2
		log.Println("Using default PARALLEL_PROCESS:", numParallelProcess)
	}

	boxPool := services.NewBoxPool(numParallelProcess)
	rmq, err := rabbitmq.NewRabbitService(rabbitURL)
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}
	defer rmq.Close()

	fileSystemProvider := testcaseprovider.NewFileSystemTestCaseProvider(testCasesPath)

	handler := func(msg domain.EvaluationInput) error {
		log.Printf("Received submission ID: %s | Language: %s | Problem: %d",
			msg.ID, msg.Language, msg.ProblemID)

		processor := services.RequestProcessor{
			Input:            msg,
			BoxPool:          boxPool,
			TestCaseProvider: fileSystemProvider,
		}

		result, err := processor.ProcessRequest()
		if err != nil {
			log.Printf("Error processing submission %s: %v", msg.ID, err)
			return err
		}

		log.Printf("Successfully processed submission %s | Passed: %d/%d",
			result.SubmitID, result.TotalPassed, result.TotalCases)

		http_request := http_request.HttpRequest{}
		body, err := http_request.SendRequest(clientResult, result)
		if err != nil {
			log.Printf("Error sending request to server: %v", err)
		}
		log.Printf("Response from server: %s", body)
		return nil
	}

	if err := rmq.Listen("submission_queue", handler); err != nil {
		log.Fatalf("Listener error: %v", err)
	}
}
