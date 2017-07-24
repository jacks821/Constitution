package main

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"os"
	"log"
	"bufio"
	"strings"
	"bytes"
	"regexp"
	"net/http"
	"net/url"
	"time"
)



func GrabLines(args string) string {
	var text string
	file, err := os.Open(args)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = text + scanner.Text() + " "
	}
	return text
}

func MakeTweets(text string) []string {
	var tweets []string
	words := strings.Split(text, " ")
	var buffer bytes.Buffer
	for i := 0; i < len(words); i++ {
		buffer.WriteString(words[i]+" ")
		if i == (len(words) -1) || (len(buffer.String()) + len(words[i+1])) > 130 {
			tweets = append(tweets, buffer.String())
			buffer.Reset()
		}
	}
	return tweets	
}

func RemoveWhitespace(text string) string {
	regex, _ := regexp.Compile("\\s+")
	return regex.ReplaceAllString(text, " ")
}

func Tweet(tweet string) {
	TWITTER_CONSUMER_KEY := 	os.Getenv("TWITTER_CONSUMER_KEY")
	TWITTER_CONSUMER_SECRET := 	os.Getenv("TWITTER_CONSUMER_SECRET")
	TWITTER_ACCESS_TOKEN := 	os.Getenv("TWITTER_ACCESS_TOKEN")
	TWITTER_ACCESS_SECRET :=	os.Getenv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(TWITTER_CONSUMER_KEY, TWITTER_CONSUMER_SECRET)
	token := oauth1.NewToken(TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_SECRET)	

	httpClient := config.Client(oauth1.NoContext, token)

	body := strings.NewReader(tweet)
	tweet := url.Values{}
	tweet.Add("status", tweet)
	
	updatePath := "https://api.twitter.com/1.1/statuses/update.json?"

	path := fmt.Sprintf("%s%s", updatePath, tweet.Encode())
	req , err := http.NewRequest("POST", path, body)

	fmt.Printf("%v", req)
	resp, _ := httpClient.Do(req)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}

	fmt.Printf("read resp.Body successfully:\n%v\n", string(data))
}






func main() {
	files := []string{"usconstitution.txt", "declaration.txt", "washington.txt"}
	for _, file := range files {
		whole := RemoveWhitespace(GrabLines(file))
		tweets := MakeTweets(whole)
		for _, tweet := range tweets {
			fmt.Printf("%v\n %v --------------\n", tweet, len(tweet))
			time.Sleep(15 * time.Second)
		}
	}
}
