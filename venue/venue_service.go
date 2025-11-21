package venue

import (
	"github.com/JkD004/playarena-backend/db"
	"github.com/JkD004/playarena-backend/notification"
	"github.com/JkD004/playarena-backend/user"
	"errors"
	"log"
)

// ---------------------------------------------------------------
// CREATE VENUE
// ---------------------------------------------------------------
func CreateNewVenue(v *Venue, ownerID int64) error {
	v.OwnerID = ownerID
	return CreateVenue(v) // repository function
}

// ---------------------------------------------------------------
// GET ALL APPROVED VENUES
// ---------------------------------------------------------------
func GetAllVenues() ([]Venue, error) {
	return FindApprovedVenues() // repository
}

// ---------------------------------------------------------------
// GET BY STATUS
// ---------------------------------------------------------------
func GetVenuesByStatus(status string) ([]Venue, error) {
	return FindVenuesByStatus(status)
}

// ---------------------------------------------------------------
// UPDATE STATUS (admin)
// ---------------------------------------------------------------
func UpdateVenueStatus(venueID int64, newStatus string) error {

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ownerID, err := GetVenueOwner(venueID) // repository new signature
	if err != nil {
		return err
	}

	// update DB
	err = UpdateVenueStatusInDB(venueID, newStatus)
	if err != nil {
		return err
	}

	// role upgrade
	if newStatus == "approved" {
		_ = user.UpdateUserRole(tx, ownerID, "owner")
		_ = notification.CreateNotification(ownerID,
			"Your venue has been approved!",
			"success",
		)
	} else if newStatus == "rejected" {
		_ = notification.CreateNotification(ownerID,
			"Your venue was rejected.",
			"error",
		)
	}

	return tx.Commit()
}

// ---------------------------------------------------------------
// GET BY ID
// ---------------------------------------------------------------
func GetVenueByID(id int64) (*Venue, error) {
	return FindApprovedVenueByID(id)
}

// ---------------------------------------------------------------
// PHOTOS
// ---------------------------------------------------------------
func GetVenuePhotos(venueID int64) ([]VenuePhoto, error) {
	return GetPhotosByVenueID(venueID)
}

func DeleteVenuePhoto(photoID int64) error {
	return DeletePhoto(photoID)
}

// ---------------------------------------------------------------
// OWNER VENUES
// ---------------------------------------------------------------
func GetVenuesForOwner(ownerID int64) ([]Venue, error) {
	return FindVenuesByOwnerID(ownerID)
}

// ---------------------------------------------------------------
// REVIEWS
// ---------------------------------------------------------------
func AddReview(venueID, userID int64, rating int, comment string) error {
	review := &Review{
		VenueID: venueID,
		UserID:  userID,
		Rating:  rating,
		Comment: comment,
	}
	return CreateReview(review) // repository correct name
}

func GetVenueReviews(venueID int64) ([]Review, error) {
	return GetReviewsByVenueID(venueID) // repository
}

// ---------------------------------------------------------------
// UPDATE VENUE INFO
// ---------------------------------------------------------------
func ModifyVenue(venueID int64, v *Venue) error {
	v.ID = venueID
	return UpdateVenueDetails(v) // repository
}

// ---------------------------------------------------------------
// OWNER CHECK
// ---------------------------------------------------------------
func VerifyVenueOwnership(venueID int64, userID int64) error {
	isOwner, err := IsVenueOwner(venueID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return errors.New("you do not own this venue")
	}
	return nil
}

// ---------------------------------------------------------------
// PHOTO → VENUE ID
// ---------------------------------------------------------------
func GetVenueIdFromPhoto(photoID int64) (int64, error) {
	return GetVenueIDByPhotoID(photoID)
}
