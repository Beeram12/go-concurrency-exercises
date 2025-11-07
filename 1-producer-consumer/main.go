//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	So here now I think we need to use goroutine and channels
	one go routine for producer and one go routine for consumer

*/
/*
	Things to remember
	1.I used buffered channels buffered channels are size defined doesnt require synchrounus communication
	  that means reciever end doesn't need to be ready
	2.While creating channels make sure what is the type you are creating intially from the previous code I created
	  channel of []*Tweets which is wrong then I changed it back to *Tweet which is obvious because I am passing only
	  values not the slices into channel
	3.In consumer function make sure go routine call is made from main not from the function because the go routine
	  will be created and then the programm will end
	4. Dont froget to use wg and defer in go routines while handling

*/

func producer(stream Stream) <-chan *Tweet {

	tweetsChan := make(chan *Tweet, len(stream.tweets))
	go func() {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				close(tweetsChan)
				return
			}
			tweetsChan <- tweet
		}
	}()
	return tweetsChan
}

/*
	ORIGINAL CODE:

	func producer(stream Stream) (tweets []*Tweet) {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				return tweets
			}

			tweets = append(tweets, tweet)
		}
	}

	func consumer(tweets []*Tweet) {
		for _, t := range tweets {
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		}
	}
*/

func consumer(tweets <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}

}

func main() {
	start := time.Now()
	stream := GetMockStream()

	// Producer
	tweets := producer(stream)

	// Consumer
	var wg sync.WaitGroup
	wg.Add(1)
	go consumer(tweets, &wg)
	wg.Wait()
	fmt.Printf("Process took %s\n", time.Since(start))
}
