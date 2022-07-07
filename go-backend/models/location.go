package models

type Location struct {
	Coordinates []float32 `json:"cords" bson:"cords"` //cords are [lat, long]
	City        string    `json:"city" bson:"city"`
	Country     string    `json:"country" bson:"country"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
}
