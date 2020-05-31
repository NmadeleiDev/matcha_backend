package structs

type UserData struct {
	Id		string	`json:"id,omitempty" bson:"id,omitempty"`
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password,omitempty" bson:"-"`
	Username string `json:"username" bson:"username"`
	Age		 int    `json:"age" bson:"age"`
	Gender   string `json:"gender" bson:"gender"`
	Country  string `json:"country" bson:"country"`
	City     string `json:"city" bson:"city"`
	MaxDist  int    `json:"max_dist" bson:"max_dist"`
	LookFor  string `json:"look_for" bson:"look_for"`
	MinAge   int    `json:"min_age" bson:"min_age"`
	MaxAge   int    `json:"max_age" bson:"max_age"`
	Images	[]string `json:"images" bson:"images,omitempty"`
	Avatar	string		`json:"avatar,omitempty" bson:"avatar,omitempty"`
	LikedBy	[]string	`json:"liked_by" bson:"liked_by"`
	LookedBy []string	`json:"looked_by" bson:"looked_by"`
	Matches	[]string	`json:"matches" bson:"matches"`
}

type LoginData struct {
	Id		string	`json:"id,omitempty" bson:"id,omitempty"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password"`
}
