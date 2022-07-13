package models

// Preference that user sets for who they want to see their globe
type PrivacyOption string

const (
	PrivacyOption_Public      PrivacyOption = "Public"
	PrivacyOption_FriendsOnly PrivacyOption = "Friends Only"
	PrivacyOption_Private     PrivacyOption = "Private"
)

var PrivacyMap = map[string]PrivacyOption{
	"Public":       PrivacyOption_Public,
	"Friends Only": PrivacyOption_FriendsOnly,
	"Private":      PrivacyOption_Private,
}

type User struct {
	Username      string              `json:"username" bson:"username"`
	Password      string              `json:"password" bson:"password"`
	Friends       []string            `json:"friends" bson:"friends"`
	Locations     map[string]Location `json:"locations" bson:"locations"`
	Trips         map[string]Trip     `json:"trips" bson:"trips"`
	PrivacyOption PrivacyOption       `json:"privacy" bson:"privacy"`
}

type Location struct {
	Coordinates []float32 `json:"cords" bson:"cords"` //cords are [lat, long]
	// City        string    `json:"city" bson:"city"`
	// Country     string    `json:"country" bson:"country"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type Trip struct {
	StartCoordinates []float32 `json:"startcords" bson:"startcords"`
	EndCoordinates   []float32 `json:"endcords" bson:"endcords"`
	Title            string    `json:"title" bson:"title"`
	Description      string    `json:"description" bson:"description"`
}
