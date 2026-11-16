package gifs

import "github.com/FacileStudio/Prisme/apps/api/internal/gifs"

// Image is one rendition of a GIF, served through Prisme.
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// GIF is one result: a preview for the grid and the rendition that is sent.
// Both URLs point at GET /api/gifs/media, never at GIPHY.
type GIF struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Preview Image  `json:"preview"`
	Full    Image  `json:"full"`
}

// PageResponse is one page of GIFs and the offset of the next page.
type PageResponse struct {
	GIFs []GIF `json:"gifs"`
	Next int   `json:"next"`
}

func image(img gifs.Image) Image {
	return Image{URL: MediaPath(img.URL), Width: img.Width, Height: img.Height}
}
