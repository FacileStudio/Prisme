package gifs

import "strconv"

type rendition struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type item struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Images struct {
		Preview rendition `json:"fixed_width_downsampled"`
		Small   rendition `json:"fixed_width"`
		Full    rendition `json:"downsized_medium"`
	} `json:"images"`
}

type listing struct {
	Data       []item `json:"data"`
	Pagination struct {
		Offset int `json:"offset"`
		Count  int `json:"count"`
	} `json:"pagination"`
}

func toPage(list listing) Page {
	out := Page{GIFs: make([]GIF, 0, len(list.Data)), Next: list.Pagination.Offset + list.Pagination.Count}
	for _, it := range list.Data {
		preview := pick(it.Images.Preview, it.Images.Small)
		full := pick(it.Images.Full, it.Images.Small)
		if preview.URL == "" || full.URL == "" {
			continue
		}
		out.GIFs = append(out.GIFs, GIF{ID: it.ID, Title: it.Title, Preview: preview, Full: full})
	}
	return out
}

func pick(first, fallback rendition) Image {
	chosen := first
	if !MediaURL(chosen.URL) {
		chosen = fallback
	}
	if !MediaURL(chosen.URL) {
		return Image{}
	}
	width, _ := strconv.Atoi(chosen.Width)
	height, _ := strconv.Atoi(chosen.Height)
	return Image{URL: chosen.URL, Width: width, Height: height}
}
