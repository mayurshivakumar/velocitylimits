package main

import (
	"bufio"
	"encoding/json"
	"os"
	"velocitylimits/models"

	"velocitylimits/cache"
	"velocitylimits/service"

	"velocitylimits/config"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	config := config.ParseConfig()
	cache := cache.NewCache()
	errGroup := errgroup.Group{}

	// go routine to read the file
	requestC, getRequest := GetRequest(config)
	errGroup.Go(getRequest)
	// go routine to attempt load and validate
	responseC, attemptLoadF := AttemptLoad(config, requestC, cache)
	go attemptLoadF()
	// go routine to write the response back to file
	responderF := Responder(config, responseC)
	errGroup.Go(responderF)

	if err := errGroup.Wait(); err != nil {
		logrus.Panicf("error. closing wait group: %v", err)
	}
}

// GetRequest reads the file and converts to request
func GetRequest(config *config.Configurations) (<-chan *models.Request, func() error) {
	requestC := make(chan *models.Request)
	parser := func() error {
		inputFile, err := OpenFile(config)
		if err != nil {
			return err
		}
		defer inputFile.Close()
		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			request, err := models.NewRequest(scanner.Text())
			if err != nil {
				logrus.Error(err)
				return err
			}
			// add the request to the request channel
			requestC <- request
			// error reading file
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		// close the channel
		close(requestC)
		return nil
	}
	return requestC, parser
}

// AttemptLoad reads the request, validates, attempts to load and writes the response back
func AttemptLoad(config *config.Configurations, requestC <-chan *models.Request, cache *cache.Cache) (<-chan *models.Response, func()) {
	responseC := make(chan *models.Response)
	attemptLoader := func() {

		for request := range requestC {
			// attempt to load
			response := service.AttemptLoad(request, config, cache)
			// adds the response to the response channel
			responseC <- response
		}
		// close the response channel
		close(responseC)
	}
	return responseC, attemptLoader
}

// Responder writes the response back to the file
func Responder(config *config.Configurations, responseC <-chan *models.Response) func() error {
	responder := func() error {
		outputFile := CreateFile(config)
		defer outputFile.Close()
		writer := bufio.NewWriter(outputFile)

		for response := range responseC {
			resBytes, err := json.Marshal(response)
			if err != nil {
				logrus.Errorf("Error marshalling json:%v", err)
				return err
			}
			// write to file
			if _, err = writer.WriteString(string(resBytes) + "\n"); err != nil {
				logrus.Errorf("Error writing to file file:%v", err)
				return err
			}

		}
		writer.Flush()
		return nil
	}

	return responder
}

func OpenFile(config *config.Configurations) (*os.File, error) {
	// TODO: Path for the file needs to be handled better
	input, err := os.Open("../" + config.VelocityLimit.InputFile)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func CreateFile(config *config.Configurations) *os.File {
	// TODO: Path for the file needs to be handled better
	output, err := os.Create("../" + config.VelocityLimit.OutputFile)
	if err != nil {
		logrus.Panicf("Unable to open file: %s", err)
	}
	return output
}
