package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"twitter-clone-api/config/database"
	"twitter-clone-api/middleware"
	"twitter-clone-api/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func PostTweetHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	createdTweet, err := postTweet(c, session)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, createdTweet)
	}

}

func uploadTweetImages(c *gin.Context, session neo4j.SessionWithContext, tweetId *int64, imageHeaders []*multipart.FileHeader) []*string {
	query := `
	MATCH (t:Tweet) WHERE id(t) = $tweetId
	CREATE (img:Image {
		filename: $filename,
		timestamp: datetime()
	})
	CREATE (t)-[:EMBEDDED]->(img)
	`
	var imageURLs []*string
	for _, imageHeader := range imageHeaders {
		url := uploadTweetFile(c, session, query, tweetId, "image", imageHeader)
		if url != "" {
			imageURLs = append(imageURLs, &url)
		}
	}
	return imageURLs
}

func postTweet(c *gin.Context, session neo4j.SessionWithContext) (models.Tweet, error) {

	var tweet models.Tweet
	jsonData := c.PostForm("data")

	if err := json.Unmarshal([]byte(jsonData), &tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return models.Tweet{}, err
	}

	username, _, _ := middleware.GetUsernameAndRoleFromCookie(c)

	query := `
	MATCH (u:User { username: $username })

	CREATE (t:Tweet {
	content: $content,
	timestamp: datetime()
	})

	CREATE (u)-[:POSTED]->(t)
	return id(t) as nodeId, t.timestamp as timestamp
`

	result, err := session.Run(c,
		query,
		map[string]any{
			"username": username,
			"content":  tweet.Content,
		})

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return models.Tweet{}, err
	}

	var createdTweet models.Tweet
	if result.Next(c) {
		createdTweet.Content = tweet.Content
		nodeID, ok := result.Record().Get("nodeId")
		if ok {
			nodeID, ok := nodeID.(int64)
			if ok {
				createdTweet.ID = &nodeID
			}
		}
		timestamp, ok := result.Record().Get("timestamp")
		if ok {
			timestamp, ok := timestamp.(time.Time)
			if ok {
				createdTweet.Timestamp = &timestamp
			}
		}

	}
	imageHeaders := c.Request.MultipartForm.File["images[]"]
	if len(imageHeaders) > 0 {
		createdTweet.Image_urls = uploadTweetImages(c, session, createdTweet.ID, imageHeaders)
		fmt.Println("image embedded")
	}

	videoHeaders := c.Request.MultipartForm.File["videos[]"]
	if len(videoHeaders) > 0 {
		createdTweet.Video_urls = uploadTweetVideo(c, session, createdTweet.ID, videoHeaders)
		fmt.Println("video embedded")
	}
	audioHeaders := c.Request.MultipartForm.File["audios[]"]
	if len(audioHeaders) > 0 {
		createdTweet.Audio_urls = uploadTweetAudio(c, session, createdTweet.ID, audioHeaders)
		fmt.Println("audio embedded")
	}

	return createdTweet, nil
}

func uploadTweetVideo(c *gin.Context, session neo4j.SessionWithContext, tweetId *int64, videoHeaders []*multipart.FileHeader) []*string {
	query := `
	MATCH (t:Tweet) WHERE id(t) = $tweetId
	CREATE (vd:Video {
		filename: $filename,
		timestamp: datetime()
	})
	CREATE (t)-[:EMBEDDED]->(vd)
	`
	var videoURLs []*string
	for _, videoHeader := range videoHeaders {
		url := uploadTweetFile(c, session, query, tweetId, "video", videoHeader)
		if url != "" {
			videoURLs = append(videoURLs, &url)
		}
	}
	return videoURLs
}
func uploadTweetAudio(c *gin.Context, session neo4j.SessionWithContext, tweetId *int64, audioHeaders []*multipart.FileHeader) []*string {
	fmt.Println("masuk audio")
	query := `
	MATCH (t:Tweet) WHERE id(t) = $tweetId
	CREATE (au:Audio {
		filename: $filename,
		timestamp: datetime()
	})
	CREATE (t)-[:EMBEDDED]->(au)
	`
	var audioURLs []*string
	for _, audioHeader := range audioHeaders {
		url := uploadTweetFile(c, session, query, tweetId, "audio", audioHeader)
		if url != "" {
			audioURLs = append(audioURLs, &url)
		}
	}
	return audioURLs
}

func uploadTweetFile(c *gin.Context, session neo4j.SessionWithContext, query string, tweetId *int64, fileFormat string, imageHeader *multipart.FileHeader) string {
	fileType := mime.TypeByExtension(filepath.Ext(imageHeader.Filename))
	if !strings.HasPrefix(fileType, (fileFormat + "/")) {
		return ""
	}

	filename, err := uploadFile(imageHeader)
	if err != nil {
		log.Println(err)
		return ""
	}

	url, err := getFileURL(c, filename)
	if err == nil {
		fmt.Println(query)
		_, err = session.Run(c,
			query,
			map[string]interface{}{
				"tweetId":  tweetId,
				"filename": filename,
			})
		if err != nil {
			log.Println(err)
			return ""
		}
		return url
	} else {
		fmt.Println(err)
		return ""
	}
}

func RetweetHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	err = retweet(c, session)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Success Retweeting"})
	}

}

func retweet(c *gin.Context, session neo4j.SessionWithContext) error {
	tweetID := c.Param("tweetID")
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	username, _, _ := middleware.GetUsernameAndRoleFromCookie(c)
	query := `
		MATCH (u:User {username: $username}), (t:Tweet)
		WHERE id(t) = $tweetID
		MERGE (u) -[:RETWEET {timestamp: datetime()}]-> (t)
	`
	result, err := session.Run(c,
		query,
		map[string]interface{}{
			"username": username,
			"tweetID":  tweetIDInt,
		},
	)
	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func QuoteRetweetHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	createdTweet, err := quoteRetweet(c, session)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, createdTweet)
	}

}

func quoteRetweet(c *gin.Context, session neo4j.SessionWithContext) (models.Tweet, error) {

	createdTweet, err := postTweet(c, session)
	if err != nil {
		log.Fatalln(err)
		return models.Tweet{}, err
	}

	tweetID := c.Param("tweetID")
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		log.Fatalln(err)
		return models.Tweet{}, err
	}

	query := `
		MATCH (t1:Tweet), (t2:Tweet)
		WHERE id(t1) = $tweetID1 AND id(t2) = $tweetID2
		MERGE (t2)-[:QUOTATION_RETWEET]->(t1)
	`

	result, err := session.Run(c,
		query,
		map[string]interface{}{
			"tweetID1": tweetIDInt,
			"tweetID2": createdTweet.ID,
		})

	if err != nil {
		return models.Tweet{}, err
	}
	if result.Err() != nil {
		return models.Tweet{}, result.Err()
	}

	return createdTweet, nil

}

func GetTweetsHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	var tweets []models.QuoteTweet

	query := `
	MATCH (n:Tweet)
	OPTIONAL MATCH (n)-[qt:QUOTATION_RETWEET]->(t:Tweet)
	RETURN id(n) as tweetId
	ORDER BY n.timestamp DESC
	`

	result, err := session.Run(c,
		query, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	for result.Next(c) {
		record := result.Record()
		tweetID, ok := record.Get("tweetId")
		if ok {
			tweetIDint, _ := tweetID.(int64)
			tweet, err := getTweet(c, session, tweetIDint)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			tweets = append(tweets, tweet)
		}
	}

	c.JSON(http.StatusOK, tweets)

}

func getTweet(c *gin.Context, session neo4j.SessionWithContext, tweetId int64) (models.QuoteTweet, error) {
	query := `
	MATCH (n:Tweet)
	WHERE id(n) = $tweetId
	OPTIONAL MATCH (n)-[qt:QUOTATION_RETWEET]->(t:Tweet)
	RETURN n.content AS content,
		   n.timestamp AS timestamp,
		   id(n) as tweetId,
		   t.content AS quote,
		   id(t) as quoteId
	`
	result, err := session.Run(c,
		query,
		map[string]interface{}{
			"tweetId": tweetId,
		})
	if err != nil {
		return models.QuoteTweet{}, err
	}

	var tweet models.QuoteTweet
	if result.Next(c) {
		record := result.Record()
		tweetID, _ := record.Get("tweetId")
		content, _ := record.Get("content")
		quote, _ := record.Get("quote")
		quoteID, _ := record.Get("quoteId")
		timestamp, _ := record.Get("timestamp")
		timestampTime := timestamp.(time.Time)

		quoteIDint, ok := quoteID.(int64)
		if ok {

			tweet.Quoted.ID = &quoteIDint
		}
		tweet.Timestamp = &timestampTime
		tweetIDint, ok := tweetID.(int64)
		if ok {

			tweet.ID = &tweetIDint
		}
		tweet.Content = content.(string)
		if quote != nil {
			tweet.Quoted.Content = quote.(string)
		}
		tweet.Image_urls, err = getTweetFiles(c, session, *tweet.ID, "Image")
		if err != nil {
			return models.QuoteTweet{}, err
		}
		tweet.Audio_urls, err = getTweetFiles(c, session, *tweet.ID, "Audio")
		if err != nil {
			return models.QuoteTweet{}, err
		}
		tweet.Video_urls, err = getTweetFiles(c, session, *tweet.ID, "Video")
		if err != nil {
			return models.QuoteTweet{}, err
		}
	}
	return tweet, nil
}

func getTweetFiles(c *gin.Context, session neo4j.SessionWithContext, tweetId int64, filetype string) ([]*string, error) {
	query := `
	MATCH (n:Tweet)-[:EMBEDDED]->(file:` + filetype + `) where id(n)=$tweetId RETURN file.filename as filename
	`

	result, err := session.Run(c,
		query,
		map[string]interface{}{
			"tweetId": tweetId,
		})
	if err != nil {
		return nil, err
	}
	var fileUrls []*string
	if result.Next(c) {
		record := result.Record()
		filename, _ := record.Get("filename")
		filenamestr, ok := filename.(string)
		if ok {
			url, err := getFileURL(c, filenamestr)
			if err != nil {
				return nil, err
			} else {
				fileUrls = append(fileUrls, &url)

			}
		}

	}
	return fileUrls, nil
}

func ReplyHandler(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	tweetID := c.Param("tweetID")
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	createdTweet, err := postTweet(c, session)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	query := `
		MATCH (t1:Tweet), (t2:Tweet)
		WHERE id(t1) = $tweetID1 AND id(t2) = $tweetID2
		CREATE (t1)-[:REPLY { timestamp: datetime() }]->(t2)
	`

	result, err := session.Run(c,
		query,
		map[string]interface{}{
			"tweetID1": createdTweet.ID,
			"tweetID2": tweetIDInt,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if result.Err() != nil {
		log.Println(result.Err())
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success Reply"})

}

func GetRepliesHandler(c *gin.Context) {
	fmt.Println("replies")
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	tweetID := c.Param("tweetID")
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var replies []models.QuoteTweet

	query := `
	MATCH (n:Tweet)-[:REPLY]->(t:Tweet)  where id(t) = $tweetId 
	RETURN id(n) as tweetId
	`

	result, err := session.Run(c, query,
		map[string]interface{}{
			"tweetId": tweetIDInt,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	for result.Next(c) {
		record := result.Record()
		tweetID, _ := record.Get("tweetId")
		tweetIDint, ok := tweetID.(int64)
		if ok {
			fmt.Println(tweetIDint)
			reply, err := getTweet(c, session, tweetIDint)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}

			replies = append(replies, reply)
		}
	}

	c.JSON(http.StatusOK, replies)

}

func GetTweetHandler(c *gin.Context) {
	fmt.Println("get tweet")
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	tweetID := c.Param("tweetID")
	tweetIDInt, err := strconv.ParseInt(tweetID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	fmt.Println(tweetIDInt)

	var tweet models.QuoteTweet

	tweet, err = getTweet(c, session, tweetIDInt)
	fmt.Println(tweet.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if tweet.Content == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.JSON(http.StatusOK, tweet)
	}
}
