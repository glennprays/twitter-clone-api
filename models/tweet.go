package models

import "time"

type Tweet struct {
	ID         *int64     `json:"id"`
	Content    string     `json:"content"`
	Timestamp  *time.Time `json:"timestamp"`
	Image_urls []*string  `json:"image_urls"`
	Video_urls []*string  `json:"video_urls"`
	Audio_urls []*string  `json:"audio_urls"`
}

type QuoteTweet struct {
	Tweet
	Quoted Tweet `json:"quoted"`
}
