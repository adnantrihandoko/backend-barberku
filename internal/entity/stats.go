package entity

type DailyStats struct {
	TotalServed       int     `json:"total_served"`
	TotalCanceled     int     `json:"total_canceled"`
	AvgWaitTimeMin    float64 `json:"avg_wait_time_min"`
	AvgServiceTimeMin float64 `json:"avg_service_time_min"`
	TotalRevenue      float64 `json:"total_revenue"`
}
