package models

type ImageInfo struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Size       string  `json:"size"`
	Layers     int    `json:"layers"`
}
