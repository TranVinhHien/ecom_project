package services

import "time"

type Media struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	FileName      string    `bson:"fileName" json:"fileName"`
	FilePath      string    `bson:"filePath" json:"filePath"`
	FileType      string    `bson:"fileType" json:"fileType"`                     // ví dụ: image/png, video/mp4
	MediaType     string    `bson:"mediaType" json:"mediaType"`                   // "image" hoặc "video"
	Size          int64     `bson:"size" json:"size"`                             // dung lượng (bytes)
	Duration      int64     `bson:"duration,omitempty" json:"duration,omitempty"` // thời lượng video (giây), nếu là ảnh thì bỏ trống
	Width         int       `bson:"width,omitempty" json:"width,omitempty"`
	Height        int       `bson:"height,omitempty" json:"height,omitempty"`
	ThumbnailPath string    `bson:"thumbnailPath,omitempty" json:"thumbnailPath,omitempty"`
	Tags          []string  `bson:"tags,omitempty" json:"tags,omitempty"`
	Description   string    `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	CreateBy      string    `bson:"createBy" json:"createBy"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
	UploadedBy    string    `bson:"uploadedBy" json:"uploadedBy"` // userId
}
