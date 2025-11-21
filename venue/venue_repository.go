package venue

import (
	"database/sql"
	"log"

	"github.com/JkD004/playarena-backend/db"
)

// ---------------------------------------------------------------
// CREATE VENUE
// ---------------------------------------------------------------
func CreateVenue(v *Venue) error {
	query := `
		INSERT INTO venues (
			owner_id, name, sport_category, description, address, price_per_hour,
			opening_time, closing_time, lunch_start_time, lunch_end_time, status
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`

	var lunchStart, lunchEnd sql.NullString
	if v.LunchStart != "" {
		lunchStart = sql.NullString{String: v.LunchStart, Valid: true}
	}
	if v.LunchEnd != "" {
		lunchEnd = sql.NullString{String: v.LunchEnd, Valid: true}
	}

	result, err := db.DB.Exec(query,
		v.OwnerID, v.Name, v.SportCategory, v.Description, v.Address, v.PricePerHour,
		v.OpeningTime, v.ClosingTime, lunchStart, lunchEnd,
	)

	if err != nil {
		log.Println("Error inserting venue:", err)
		return err
	}

	id, _ := result.LastInsertId()
	v.ID = id
	v.Status = "pending"
	return nil
}

// ---------------------------------------------------------------
// GET VENUE OWNER (Used by Upload/Delete/Update logic)
// ---------------------------------------------------------------
func GetVenueOwner(venueID int64) (int64, error) {
	var ownerID int64
	query := `SELECT owner_id FROM venues WHERE id = ?`

	err := db.DB.QueryRow(query, venueID).Scan(&ownerID)
	if err != nil {
		log.Println("Error getting venue owner:", err)
		return 0, err
	}
	return ownerID, nil
}

// ---------------------------------------------------------------
// CHECK IF USER OWNS VENUE
// ---------------------------------------------------------------
func IsVenueOwner(venueID int64, ownerID int64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM venues WHERE id = ? AND owner_id = ?`

	err := db.DB.QueryRow(query, venueID, ownerID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ---------------------------------------------------------------
// GET ALL APPROVED VENUES
// ---------------------------------------------------------------
func FindApprovedVenues() ([]Venue, error) {
	return FindVenuesByStatus("approved")
}

// ---------------------------------------------------------------
// GET VENUES BY STATUS (admin)
// ---------------------------------------------------------------
func FindVenuesByStatus(status string) ([]Venue, error) {
	query := `
		SELECT id, owner_id, status, name, sport_category, description, address,
		       price_per_hour, opening_time, closing_time, lunch_start_time,
		       lunch_end_time, created_at
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
// GetAllVenues returns ALL venues (no status filter)
func GetAllVenues() ([]Venue, error) {
    query := `
        SELECT id, owner_id, status, name, sport_category, description, 
               address, price_per_hour, opening_time, closing_time,
               lunch_start_time, lunch_end_time, created_at
        FROM venues
    `
    rows, err := db.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    venues := make([]Venue, 0)

    for rows.Next() {
        v, err := scanVenue(rows)
        if err == nil {
            venues = append(venues, *v)
        }
    }

    return venues, nil
}


// ---------------------------------------------------------------
// FIND VENUE BY ID (only approved)
// ---------------------------------------------------------------
func FindApprovedVenueByID(id int64) (*Venue, error) {
	query := `
		SELECT id, owner_id, status, name, sport_category, description, address,
		       price_per_hour, opening_time, closing_time, lunch_start_time,
		       lunch_end_time, created_at
		FROM venues
		WHERE id = ? AND status = 'approved'
	`

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanVenue(rows)
	}

	return nil, sql.ErrNoRows
}

// ---------------------------------------------------------------
// UPDATE VENUE STATUS
// ---------------------------------------------------------------
func UpdateVenueStatusInDB(venueID int64, newStatus string) error {
	query := `UPDATE venues SET status = ? WHERE id = ?`

	_, err := db.DB.Exec(query, newStatus, venueID)
	if err != nil {
		log.Println("Error updating venue status:", err)
		return err
	}
	return nil
}

// ---------------------------------------------------------------
// GET PHOTOS BY VENUE ID
// ---------------------------------------------------------------
func GetPhotosByVenueID(venueID int64) ([]VenuePhoto, error) {
	query := `SELECT id, venue_id, image_url, created_at FROM venue_photos WHERE venue_id = ?`

	rows, err := db.DB.Query(query, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []VenuePhoto

	for rows.Next() {
		var p VenuePhoto
		rows.Scan(&p.ID, &p.VenueID, &p.ImageURL, &p.CreatedAt)
		photos = append(photos, p)
	}

	return photos, nil
}

// ---------------------------------------------------------------
// DELETE PHOTO (admin/owner validated at service layer)
// ---------------------------------------------------------------
func DeletePhoto(photoID int64) error {
	query := `DELETE FROM venue_photos WHERE id = ?`
	_, err := db.DB.Exec(query, photoID)
	return err
}

// ---------------------------------------------------------------
// GET VENUE ID FROM PHOTO
// ---------------------------------------------------------------
func GetVenueIDByPhotoID(photoID int64) (int64, error) {
	var venueID int64
	query := `SELECT venue_id FROM venue_photos WHERE id = ?`

	err := db.DB.QueryRow(query, photoID).Scan(&venueID)
	if err != nil {
		return 0, err
	}

	return venueID, nil
}

// ---------------------------------------------------------------
// SCAN VENUE ROW
// ---------------------------------------------------------------
func scanVenue(rows *sql.Rows) (*Venue, error) {
	var v Venue
	var desc, addr, lStart, lEnd sql.NullString
	var price sql.NullFloat64
	var created sql.NullTime

	err := rows.Scan(
		&v.ID, &v.OwnerID, &v.Status, &v.Name, &v.SportCategory,
		&desc, &addr, &price, &v.OpeningTime, &v.ClosingTime,
		&lStart, &lEnd, &created,
	)

	if err != nil {
		return nil, err
	}

	v.Description = desc.String
	v.Address = addr.String
	v.PricePerHour = price.Float64
	v.LunchStart = lStart.String
	v.LunchEnd = lEnd.String
	if created.Valid {
		v.CreatedAt = created.Time
	}

	return &v, nil
}
