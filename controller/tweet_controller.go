package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"twitter-clone-api/config/database"
	"twitter-clone-api/middleware"
	"twitter-clone-api/models"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func PostTweet(c *gin.Context) {
	driver, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer driver.Close(c)

	session := driver.NewSession(c, neo4j.SessionConfig{})
	defer session.Close(c)

	var tweet models.Tweet
	jsonData := c.PostForm("data")

	if err := json.Unmarshal([]byte(jsonData), &tweet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		return
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
	fmt.Println("after")
	if len(imageHeaders) > 0 {
		fmt.Println("image embedded")
		createdTweet.Image_urls = uploadTweetImages(c, session, createdTweet.ID, imageHeaders)
	}

	videoHeaders := c.Request.MultipartForm.File["videos[]"]
	if len(videoHeaders) > 0 {
		createdTweet.Video_urls = uploadTweetVideo(c, session, createdTweet.ID, videoHeaders)
	}
	audioHeaders := c.Request.MultipartForm.File["audios[]"]
	if len(audioHeaders) > 0 {
		createdTweet.Audio_urls = uploadTweetAudio(c, session, createdTweet.ID, audioHeaders)
	}
	createdTweet.Content = tweet.Content
	c.JSON(http.StatusOK, createdTweet)
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
