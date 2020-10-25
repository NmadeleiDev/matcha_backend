package types

type UserData struct {
	Id        string   `json:"id,omitempty" bson:"id,omitempty"`
	Email     string   `json:"email" bson:"email"`
	Phone     string   `json:"phone" bson:"phone"`
	Password  string   `json:"password,omitempty" bson:"-"`
	Username  string   `json:"username" bson:"username"`
	Name  string   `json:"name" bson:"name"`
	Surname  string   `json:"surname" bson:"surname"`
	BirthDate int64    `json:"birthDate" bson:"birth_date"`
	Gender    string   `json:"gender" bson:"gender"`
	Country   string   `json:"country" bson:"country"`
	City      string   `json:"city" bson:"city"`
	MaxDist   int      `json:"maxDist" bson:"max_dist"`
	LookFor   string   `json:"lookFor" bson:"look_for"`
	MinAge    int      `json:"minAge" bson:"min_age"`
	MaxAge    int      `json:"maxAge" bson:"max_age"`
	Images    []string `json:"images" bson:"images"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	LikedBy   []string `json:"likedBy" bson:"liked_by"`
	LookedBy []string	`json:"lookedBy" bson:"looked_by"`
	Matches	[]string	`json:"matches" bson:"matches"`
	GeoPosition Coordinates	`json:"position,omitempty"`
}

type Coordinates struct {
	Lat			float64		`json:"lat"`
	Lon			float64		`json:"lon"`
}

type LoginData struct {
	Id		string	`json:"id,omitempty" bson:"id,omitempty"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password"`
}

type Message struct {
	Id					string		`json:"id" bson:"id"`
	Sender				string		`json:"sender" bson:"sender"`
	Recipient				string		`json:"recipient" bson:"recipient"`
	Date				int			`json:"date" bson:"date"`
	State				int			`json:"state" bson:"state"`
	Text				string		`json:"text" bson:"text"`
}

type VerifyRequest struct {
	AuthKey		string		`json:"authKey"`
}

