package bitrix24

type BitrixRestResponseFind struct {
	Result []BitrxiRestResult `json:"result"`
	Total  int                `json:"total"`
}

type BitrixRestResponseUpdate struct {
	Result bool `json:"result"`
}

type BitrixRestResponseAdd struct {
	Result int `json:"result"`
}

type BitrxiRestResult struct {
	ID string `json:"ID"`
}