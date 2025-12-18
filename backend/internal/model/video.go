package model

// VideoItem represents a video in search results
type VideoItem struct {
	VodName    string       `json:"vod_name"`
	VodPic     string       `json:"vod_pic"`
	VodRemarks string       `json:"vod_remarks,omitempty"`
	TypeName   string       `json:"type_name,omitempty"`
	Sources    []SourceInfo `json:"sources"`
}

// SourceInfo represents source information for a video
type SourceInfo struct {
	SourceCode string `json:"source_code"`
	SourceName string `json:"source_name"`
	VodID      int    `json:"vod_id"`
}

// VideoDetail represents video detail information
type VideoDetail struct {
	VodName     string   `json:"vod_name"`
	VodPic      string   `json:"vod_pic"`
	VodContent  string   `json:"vod_content,omitempty"`
	VodYear     string   `json:"vod_year,omitempty"`
	VodArea     string   `json:"vod_area,omitempty"`
	VodDirector string   `json:"vod_director,omitempty"`
	VodActor    string   `json:"vod_actor,omitempty"`
	Episodes    []string `json:"episodes"`
}
