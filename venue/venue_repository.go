// venue/venue_repository.go
package venue

import (
	"database/sql"
	"errors"
	"log"

	"github.com/JkD004/playarena-backend/db"
)

// ----------------------------------------------------
// CREATE VENUE
// ----------------------------------------------------
func CreateVenue(v *Venue) error {

	query := `
		INSERT INTO venues (
			owner_id, status, name, sport_category,
			description, address, price_per_hour,
			created_at, opening_time, closing_time,
			lunch_start_time, lunch_end_time
		)
		VALUES (?, 'pending', ?, ?, ?, ?, ?, NOW(), ?, ?, ?, ?)
	`

	var lunchStart, lunchEnd sql.NullString
	if v.LunchStart != "" {
		lunchStart = sql.NullString{String: v.LunchStart, Valid: true}
	}
	if v.LunchEnd != "" {
		lunchEnd = sql.NullString{String: v.LunchEnd, Valid: true}
	}

	res, err := db.DB.Exec(query,
		v.OwnerID, v.Name, v.SportCategory,
		v.Description, v.Address, v.PricePerHour,
		v.OpeningTime, v.ClosingTime,
		lunchStart, lunchEnd,
	)

	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	v.ID = id
	return nil
}

// ----------------------------------------------------
// FETCH VENUES BY STATUS
// ----------------------------------------------------
func FindVenuesByStatus(status string) ([]Venue, error) {

	query := `
		SELECT 
			id, owner_id, status, name, sport_category,
			description, address, price_per_hour,
			created_at, opening_time, closing_time,
			lunch_start_time, lunch_end_time
		FROM venues
		WHERE status = ?
	`

	rows, err := db.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []Venue

	for rows.Next() {
		v, err := scanVenue(rows)
		if err == nil {
			venues = append(venues, *v)
		}
	}

	return venues, nil
}

// ALL APPROVED VENUES
func FindApprovedVenues() ([]Venue, error) {
	return FindVenuesByStatus("approved")
}

// ----------------------------------------------------
// FIND VENUE BY ID (only approved)
// ----------------------------------------------------
func FindApprovedVenueByID(id int64) (*Venue, error) {

	query := `
		SELECT 
			id, owner_id, status, name, sport_category,
			description, address, price_per_hour,
			created_at, opening_time, closing_time,
			lunch_start_time, lunch_end_time
		FROM venues
		WHERE id = ? AND status = 'approved'
	`

	res, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		return scanVenue(res)
	}

	return nil, sql.ErrNoRows
}

// ----------------------------------------------------
// UPDATE VENUE STATUS
// ----------------------------------------------------
func UpdateVenueStatusInDB(tx *sql.Tx, id int64, status string) error {

	query := `UPDATE venues SET status = ? WHERE id = ?`

	_, err := tx.Exec(query, status, id)
	return err
}

// ----------------------------------------------------
// FETCH OWNER'S VENUES
// ----------------------------------------------------
func FindVenuesByOwnerID(ownerID int64) ([]Venue, error) {

	query := `
		SELECT 
			id, owner_id, status, name, sport_category,
			description, address, price_per_hour,
			created_at, opening_time, closing_time,
			lunch_start_time, lunch_end_time
		FROM venues
		WHERE owner_id = ?
	`

	rows, err := db.DB.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []Venue

	for rows.Next() {
		v, err := scanVenue(rows)
		if err == nil {
			venues = append(venues, *v)
		}
	}

	return venues, nil
}

// ----------------------------------------------------
// PHOTO FUNCTIONS
// ----------------------------------------------------
func GetPhotosByVenueID(venueID int64) ([]VenuePhoto, error) {
	query := "SELECT id, venue_id, image_url, created_at FROM venue_photos WHERE venue_id = ?"

	rows, err := db.DB.Query(query, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []VenuePhoto

	for rows.Next() {
		var p VenuePhoto
		err := rows.Scan(&p.ID, &p.VenueID, &p.ImageURL, &p.CreatedAt)
		if err == nil {
			photos = append(photos, p)
		}
	}

	return photos, nil
}

func DeletePhoto(id int64) error {
	_, err := db.DB.Exec("DELETE FROM venue_photos WHERE id = ?", id)
	return err
}

func GetVenueIDByPhotoID(photoID int64) (int64, error) {
	var venueID int64
	err := db.DB.QueryRow("SELECT venue_id FROM venue_photos WHERE id = ?", photoID).Scan(&venueID)
	return venueID, err
}

// ----------------------------------------------------
// REVIEWS
// ----------------------------------------------------
func CreateReview(review *Review) error {
	query := `
		INSERT INTO reviews (venue_id, user_id, rating, comment)
		VALUES (?, ?, ?, ?)
	`
	_, err := db.DB.Exec(query, review.VenueID, review.UserID, review.Rating, review.Comment)
	return err
}

func GetReviewsByVenueID(venueID int64) ([]Review, error) {

	query := `
		SELECT r.id, r.venue_id, r.user_id,
			   u.first_name, u.last_name,
			   r.rating, r.comment, r.created_at
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.venue_id = ?
		ORDER BY r.created_at DESC
	`

	rows, err := db.DB.Query(query, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review

	for rows.Next() {
		var r Review
		err := rows.Scan(
			&r.ID, &r.VenueID, &r.UserID,
			&r.UserFirst, &r.UserLast,
			&r.Rating, &r.Comment, &r.CreatedAt,
		)
		if err == nil {
			reviews = append(reviews, r)
		}
	}

	return reviews, nil
}

// ----------------------------------------------------
// UPDATE VENUE DETAILS
// ----------------------------------------------------
func UpdateVenueDetails(v *Venue) error {

	query := `
		UPDATE venues SET 
			name=?, sport_category=?, description=?, address=?, price_per_hour=?,
			opening_time=?, closing_time=?, lunch_start_time=?, lunch_end_time=?
		WHERE id=?
	`

	_, err := db.DB.Exec(query,
		v.Name, v.SportCategory, v.Description, v.Address, v.PricePerHour,
		v.OpeningTime, v.ClosingTime, v.LunchStart, v.LunchEnd,
		v.ID,
	)

	return err
}

// ----------------------------------------------------
// SCAN VENUE (IMPORTANT: must match database column order)
// ----------------------------------------------------
func scanVenue(rows *sql.Rows) (*Venue, error) {

	var v Venue
	var desc, addr, lunchStart, lunchEnd sql.NullString
	var price sql.NullFloat64
	var created sql.NullTime

	err := rows.Scan(
		&v.ID, &v.OwnerID, &v.Status, &v.Name, &v.SportCategory,
		&desc, &addr, &price,
		&created, &v.OpeningTime, &v.ClosingTime,
		&lunchStart, &lunchEnd,
	)

	if err != nil {
		return nil, err
	}

	v.Description = desc.String
	v.Address = addr.String
	v.PricePerHour = price.Float64
	v.LunchStart = lunchStart.String
	v.LunchEnd = lunchEnd.String
	if created.Valid {
		v.CreatedAt = created.Time
	}

	return &v, nil
}
