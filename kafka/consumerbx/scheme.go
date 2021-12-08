package consumerbx

type BitrixRestResponse struct {
	Result []BitrxiRestResult `json:"result"`
	Total  int                `json:"total"`
}

type BitrxiRestResult struct {
	ID string `json:"ID"`
}

type Config struct {
	Brokers string `json:"brokers"`
	Topic   string `json:"topic"`
	GroupId string `json:"groupId"`
	RESTUrl string `json:"restUrl"`
}
