package models

// Preference that user sets for who they want to see their globe
type PrivacyOption string

// Maps cords to location object
type LocationsMap map[string]Location

const (
	PrivacyOption_Public      PrivacyOption = "Public"
	PrivacyOption_FriendsOnly PrivacyOption = "Friends Only"
	PrivacyOption_Private     PrivacyOption = "Private"
)

var PrivacyMap = map[string]PrivacyOption{
	"Public":       PrivacyOption_Public,
	"Friends Only": PrivacyOption_FriendsOnly,
	"Private":      PrivacyOption_Private}

type User struct {
	Username      string        `json:"username" bson:"username"`
	Password      string        `json:"password" bson:"password"`
	Friends       []string      `json:"friends" bson:"friends"`
	Locations     LocationsMap  `json:"locations" bson:"locations"`
	PrivacyOption PrivacyOption `json:"privacy" bson:"privacy"`
}
