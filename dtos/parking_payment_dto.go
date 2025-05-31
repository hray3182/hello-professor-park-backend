package dtos

// ParkingPaymentPayload defines the JSON structure for paying a parking record.
type ParkingPaymentPayload struct {
	PaymentMethod string  `json:"paymentMethod" binding:"required" example:"MobilePay"`
	AmountPaid    float64 `json:"amountPaid" binding:"required" example:"50.00"`
	// 可選，如果前端有來自支付閘道的參考ID或備註
	PaymentReference string `json:"paymentReference,omitempty" example:"TXN_REF_123XYZ"`
}
