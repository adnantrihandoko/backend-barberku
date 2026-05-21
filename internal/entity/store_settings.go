package entity

type StoreSettings struct {
	OpenHour     int `json:"open_hour"`
	CloseHour    int `json:"close_hour"`
	MaxQueueSize int `json:"max_queue_size"`
}
