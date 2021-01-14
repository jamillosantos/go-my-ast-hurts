package models

/** Test multiline if 
	is multiline 
*/
type User struct {
	// ID User
	ID        int64     `json:"id,uuidTest"`
	/* Name user */
	Name      string    `json:"name" bson:""`
	// User email.
	// Test.
	Email     string    `json:"email"`
	Pass      string    `json:"pass"`
}