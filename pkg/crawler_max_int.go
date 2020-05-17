package pkg

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var urls = []string{
	"http://jack-random.herokuapp.com/number/7",
	"http://jack-random.herokuapp.com/number/9",
	"http://jack-random.herokuapp.com/number/65",
	"http://jack-random.herokuapp.com/number/34",
	"http://jack-random.herokuapp.com/number/16",
	"http://jack-random.herokuapp.com/number/77",
}

func usingWaitGroup() int {
	var wg sync.WaitGroup
	var myMutex sync.Mutex
	var maxNumber = math.MinInt64

	for _, url := range urls {
		wg.Add(1)

		go func(url string, wg *sync.WaitGroup) {
			num, _ := get(url)
			myMutex.Lock()
			if *num > maxNumber {
				maxNumber = *num
			}
			myMutex.Unlock()
			wg.Done()
		}(url, &wg)
	}

	wg.Wait()

	return maxNumber
}

func usingBufferChannel() int {
	myIntChannel := make(chan int, len(urls))
	var maxNumber = math.MinInt64

	for _, url := range urls {
		go func(url string, myIntChannel chan int) {
			num, _ := get(url)
			myIntChannel <- *num
		}(url, myIntChannel)
	}

	for i := 0; i < len(urls); i++ {
		number := <-myIntChannel
		if maxNumber < number {
			maxNumber = number
		}
	}

	return maxNumber
}

func usingChannel() int {
	myIntChannel := make(chan int)
	var maxNumber = math.MinInt64

	for _, url := range urls {
		go func(url string, myIntChannel chan int) {
			num, _ := get(url)
			myIntChannel <- *num
		}(url, myIntChannel)
	}

	for i := 0; i < len(urls); i++ {
		number := <-myIntChannel
		if maxNumber < number {
			maxNumber = number
		}
	}

	// cannot use range like this as we have to close the channel
	//for number := range myIntChannel {
	//	if maxNumber < number {
	//		maxNumber = number
	//	}
	//}

	return maxNumber
}

func get(url string) (*int, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			fmt.Printf("Cannot close response: %+v", err)
		}
	}()

	// if response.StatusCode == 200
	dataInBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result, err := strconv.Atoi(string(dataInBytes))
	return &result, err
}
