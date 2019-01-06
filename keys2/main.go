package main

import (
	. "fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"

	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	. "os"
)

var credential *oauth.Credentials

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*.tmpl")

	r.GET("/", func(c *gin.Context) {
		anaconda := getTwitterApi()
		uri, cred, err := anaconda.AuthorizationURL("http://localhost:8080/twitter/callback")
		if err != nil {
			c.JSON(500, err)
		}
		Println("----------TOKEN----------")
		Println(cred.Token)
		Println("----------Secret----------")
		Println(cred.Secret)
		// credential = cred
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"twitter_url": uri,
		})

		v := url.Values{}
		Println(v)
		v.Set("count", "10")

		tweets, err := anaconda.GetHomeTimeline(v)
		if err != nil {
			panic(err)
		}

		for _, tweet := range tweets {
			Println("tweet: ", tweet.Text)
		}
	})

	r.GET("/twitter/callback", func(c *gin.Context) {
		// anaconda := getTwitterApi()
		// cre, _, err := anaconda.GetCredentials(credential, c.Query("oauth_verifier"))
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	r.Run(":8080")
}

func getTwitterApi() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials("", "", Getenv("TWITTER_CONSUMER_KEY"), Getenv("TWITTER_CONSUMER_SECRET"))
}
