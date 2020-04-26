package structs

type UserData struct {
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone" bson:"phone"`
	Password string `json:"password" bson:"-"`
	Username string `json:"username" bson:"username"`
	Age		 int    `json:"age" bson:"age"`
	Gender   string `json:"gender" bson:"gender"`
	Country  string `json:"country" bson:"country"`
	City     string `json:"city" bson:"city"`
	MaxDist  int    `json:"max_dist" bson:"max_dist"`
	LookFor  string `json:"look_for" bson:"look_for"`
	MinAge   int    `json:"min_age" bson:"min_age"`
	MaxAge   int    `json:"max_age" bson:"max_age"`
}

type LoginData struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password"`
}
