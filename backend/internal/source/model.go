package source

// MacCMSResponse is the MacCMS v10 API response structure
type MacCMSResponse struct {
	Code  int        `json:"code"`
	Msg   string     `json:"msg"`
	List  []RawVideo `json:"list"`
	Total int        `json:"total"`
}

// RawVideo is the raw video data from API
type RawVideo struct {
	VodID       int    `json:"vod_id"`
	VodName     string `json:"vod_name"`
	VodPic      string `json:"vod_pic"`
	VodRemarks  string `json:"vod_remarks"`
	TypeName    string `json:"type_name"`
	VodPlayURL  string `json:"vod_play_url"`
	VodContent  string `json:"vod_content"`
	VodYear     string `json:"vod_year"`
	VodArea     string `json:"vod_area"`
	VodDirector string `json:"vod_director"`
	VodActor    string `json:"vod_actor"`

	// Injected fields
	SourceCode string `json:"-"`
	SourceName string `json:"-"`
}
