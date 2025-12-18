/** Source info */
export interface SourceInfo {
	source_code: string;
	source_name: string;
	vod_id: number;
}

/** Video item in search results */
export interface VideoResult {
	vod_name: string;
	vod_pic: string;
	vod_remarks?: string;
	type_name?: string;
	sources: SourceInfo[];
}

/** Video detail */
export interface VideoDetail {
	vod_name: string;
	vod_pic: string;
	vod_content?: string;
	vod_year?: string;
	vod_area?: string;
	vod_director?: string;
	vod_actor?: string;
	episodes: string[];
}

/** Search API response */
export interface SearchResponse {
	code: number;
	msg?: string;
	list: VideoResult[];
}

/** Detail API response */
export interface DetailResponse {
	code: number;
	msg?: string;
	data: VideoDetail;
}

/** Source list item */
export interface SourceListItem {
	code: string;
	name: string;
}

/** Sources API response */
export interface SourcesResponse {
	code: number;
	msg?: string;
	list: SourceListItem[];
}
