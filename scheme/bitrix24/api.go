package bitrix24

type BitrixRestResponse struct {
	Result []BitrxiRestResult `json:"result"`
	Total  int                `json:"total"`
}

type BitrixRestResponseUpdate struct {
	Result bool `json:"result"`
}

type BitrxiRestResult struct {
	ID string `json:"ID"`
}