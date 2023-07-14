package models

import "time"

type Tweet struct {
	ID        *int64     `json:"id"`
	Content   string     `json:"content"`
	Timestamp *time.Time `json:"timestamp"`
	Image_url *string    `json:"image_url"`
	Video_url *string    `json:"video_url"`
	Audio_url *string    `json:"audio_url"`
}
