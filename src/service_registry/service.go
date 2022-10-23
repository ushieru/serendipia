package serviceregistry

type Service struct {
	Name      string `json:"name"`
	Ip        string `json:"ip"`
	Protocol  string `json:"protocol"`
	Port      string `json:"port"`
	Timestamp int64  `json:"timestamp"`
}
