package booking

// Booking
type Booking struct {
	BookingID int    `json:"bookingId"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Paid      bool   `json:"paid"`
}
