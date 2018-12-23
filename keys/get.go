package main

import (
	. "fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
)

func main() {
	api := GetTwitterApi()

	v := url.Values{}
	v.Set("count", "10")

	tweets, err := api.GetHomeTimeline(v)
	if err != nil {
		panic(err)
	}

	for _, tweet := range tweets {
		Println("tweet: ", tweet.Text)
	}

	text := "Hello world from golang"
	tweet, err := api.PostTweet(text, nil)
	if err != nil {
		panic(err)
	}

	Print(tweet.Text)

}

func GetTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(Getenv("TWITTER_ACCESS_TOKEN"), Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	return api
}
