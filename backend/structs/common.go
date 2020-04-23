package structs

type UserData struct {
	Email		string		`json:"email"`
	Phone		string		`json:"phone"`
	Password	string		`json:"password"`
	Username	string		`json:"username"`
	BornDate	int			`json:"born_date"`
	Gender		string		`json:"gender"`
	Country		string		`json:"country"`
	City		string		`json:"city"`
	MaxDist		int			`json:"max_dist"`
	LookFor		string		`json:"look_for"`
	MinAge		int			`json:"min_age"`
	MaxAge		int			`json:"max_age"`
}

type LoginData struct {
	Email		string		`json:"email"`
	Password	string		`json:"password"`
}

