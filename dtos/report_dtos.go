package dtos

// TotalParkingCountResponse defines the structure for total parking count response
type TotalParkingCountResponse struct {
	TotalCount int64 `json:"total_count"`
}

// TotalRevenueResponse defines the structure for total revenue response
type TotalRevenueResponse struct {
	TotalRevenue float64 `json:"total_revenue"`
	Currency     string  `json:"currency"` // e.g., "TWD", "USD"
}

// ImageAttachmentRateResponse defines the structure for image attachment rate
type ImageAttachmentRateResponse struct {
	TotalEntries     int64   `json:"total_entries"`
	EntriesWithImage int64   `json:"entries_with_image"`
	AttachmentRate   float64 `json:"attachment_rate"` // Value between 0.0 and 1.0
}

// AvailableSpotsResponse defines the structure for available parking spots response
type AvailableSpotsResponse struct {
	TotalCapacity  int   `json:"total_capacity"`
	OccupiedSpots  int64 `json:"occupied_spots"`
	AvailableSpots int64 `json:"available_spots"`
}
