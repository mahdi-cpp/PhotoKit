package model

// UIImage An object that manages image data in your app.
type UIImage struct {
	Name        string    `json:"name"`
	FileType    string    `json:"fileType"`
	Orientation int       `json:"orientation"`
	AspectRatio float32   `json:"aspectRatio"`
	Size        CGSize    `json:"size"`
	VideoInfo   VideoInfo `json:"videoInfo"`
}

type CGSize struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type VideoInfo struct {
	IsVideo         bool   `json:"isVideo"`
	VideoDuration   int    `json:"videoDuration"`
	HasSubtitle     bool   `json:"hasSubtitle"`
	VideoFormat     string `json:"videoFormat"`
	HasVideoControl bool   `json:"hasVideoControl"`
}
