// booking/booking_model.go
package booking

import "time"

type Booking struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	VenueID       int64     `json:"venue_id"`
	VenueName     string    `json:"venue_name"`      // <-- NEW
	SportCategory string    `json:"sport_category"`  // <-- NEW
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// ... (keep CreateBookingRequest struct)

// Add a struct for the request body, as users won't send everything
type CreateBookingRequest struct {
	VenueID   int64     `json:"venue_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	// Price will be calculated on the backend
}

// AdminBookingView includes venue name and user name
type AdminBookingView struct {
	BookingID     int64     `json:"booking_id"`
	VenueID       int64     `json:"venue_id"`
	VenueName     string    `json:"venue_name"`
	SportCategory string    `json:"sport_category"`
	UserID        int64     `json:"user_id"`
	UserFirstName string    `json:"user_first_name"` // <-- ADD THIS
	UserLastName  string    `json:"user_last_name"`  // <-- ADD THIS
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
}

// OwnerStats defines the data for the owner's dashboard
type OwnerStats struct {
	TotalBookings int64   `json:"total_bookings"`
	TotalRevenue  float64 `json:"total_revenue"`
	PopularTime   string  `json:"popular_time"`
}

// VenueStats defines the stats for a single venue
type VenueStats struct {
	VenueID       int64   `json:"venue_id"`
	VenueName     string  `json:"venue_name"`
	SportCategory string  `json:"sport_category"`
	TotalBookings int64   `json:"total_bookings"`
	TotalRevenue  float64 `json:"total_revenue"`
}// booking/booking_model.go

type BookedSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

