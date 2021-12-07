package consumer

type Config struct {
	Brokers string `json:"brokers"`
	Topic   string `json:"topic"`
	GroupId string `json:"groupId"`
	DBMS    string `json:"dbms"`
	CRMDataline string `json:"crm-dataline"`
}
