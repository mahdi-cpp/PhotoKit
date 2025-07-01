package models

// UIImage An object that manages image data in your app.
type UIImage struct {
	Named       string    `json:"named"`
	Format      string    `json:"format"`
	Orientation int       `json:"orientation"`
	AspectRatio float32   `json:"aspectRatio"`
	Size        CGSize    `json:"size"`
	VideoInfo   VideoInfo `json:"videoInfo"`
	VideoFormat string    `json:"videoFormat"`
}

type CGSize struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type VideoInfo struct {
	IsVideo       bool   `json:"isVideo"`
	VideoDuration int    `json:"videoDuration"`
	HasSubtitle   bool   `json:"hasSubtitle"`
	VideoFormat   string `json:"videoFormat"`
}
