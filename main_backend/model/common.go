package model

import "time"

type FullUserData struct {
	Id        string   `json:"id,omitempty" bson:"id,omitempty"`
	Email     string   `json:"email,omitempty" bson:"email"`
	Phone     string   `json:"phone,omitempty" bson:"phone"`
	Password  string   `json:"password,omitempty" bson:"-"`
	Username  string   `json:"username" bson:"username"`
	Name  string   `json:"name" bson:"name"`
	Surname  string   `json:"surname" bson:"surname"`
	BirthDate int64    `json:"birthDate" bson:"birth_date"`
	Gender    string   `json:"gender" bson:"gender"`
	Country   string   `json:"country" bson:"country"`
	City      string   `json:"city" bson:"city"`
	Bio      string   `json:"bio" bson:"bio"`
	MaxDist   int      `json:"maxDist,omitempty" bson:"max_dist"`
	LookFor   string   `json:"lookFor,omitempty" bson:"look_for"`
	MinAge    int      `json:"minAge,omitempty" bson:"min_age"`
	MaxAge    int      `json:"maxAge,omitempty" bson:"max_age"`
	Images    []string `json:"images" bson:"images"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	LikedBy   []string `json:"likedBy,omitempty" bson:"liked_by"`
	LookedBy []string	`json:"lookedBy,omitempty" bson:"looked_by"`
	Matches	[]string	`json:"matches,omitempty" bson:"matches"`
	TagIds	[]int64		`json:"-" bson:"tag_ids"`
	Tags	[]string		`json:"tags" bson:"-"`
	BannedUserIds	[]string	`json:"bannedUserIds,omitempty" bson:"banned_user_ids"`
	UseLocation	bool	`json:"useLocation,omitempty" bson:"use_location"`
	GeoPosition Coordinates	`json:"position,omitempty" bson:"-"`
	MongoLocation MongoCoors	`json:"-" bson:"position,omitempty"`

	IsOnline	bool	`json:"isOnline" bson:"is_online"`
	Rating		float64	`json:"rating" bson:"rank"`
}

func (d *FullUserData) ConvertFromDbCoords() {
	d.GeoPosition.Lon = d.MongoLocation.Coordinates[0]
	d.GeoPosition.Lat = d.MongoLocation.Coordinates[1]
}

func (d *FullUserData) ConvertToDbCoords() {
	d.MongoLocation.Type = "Point"
	d.MongoLocation.Coordinates = []float64{d.GeoPosition.Lon, d.GeoPosition.Lat}
}

type ShortUserData struct {
	Id        string   `json:"id,omitempty" bson:"id,omitempty"`
	Username  string   `json:"username" bson:"username"`
	BirthDate int64    `json:"birthDate" bson:"birth_date"`
	Gender    string   `json:"gender" bson:"gender"`
	Avatar    string   `json:"avatar" bson:"avatar"`
	Images    []string `json:"-" bson:"images"` // для того, чтобы при выгрузке взять аватар из картинок, если сам аватар не установлен. На фронт никогда не передается
	City      string   `json:"city" bson:"city"`
	Country   string   `json:"country" bson:"country"`
}

type Coordinates struct {
	Lat			float64		`json:"lat"`
	Lon			float64		`json:"lon"`
}

type MongoCoors struct {
	Type string			`bson:"type"`
	Coordinates	[]float64		`bson:"coordinates"`
}

type LoginData struct {
	Id		string	`json:"id,omitempty" bson:"id,omitempty"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password"`
}

type Chat struct {
	Id			string		`json:"id"`
	UserIds		[]string	`json:"userIds"`
	Messages	[]Message	`json:"messages"`
}

func (c *Chat) LastMessage() *Message  {
	if len(c.Messages) > 0 {
		return &c.Messages[len(c.Messages) - 1]
	}
	return nil
}

type Message struct {
	Id					string		`json:"id" bson:"id"`
	Sender				string		`json:"sender" bson:"sender"`
	Recipient				string		`json:"recipient" bson:"recipient"`
	ChatId				string		`json:"chatId" bson:"chat_id"`
	Date				int			`json:"date" bson:"date"`
	State				int			`json:"state" bson:"state"`
	Text				string		`json:"text" bson:"text"`
}

type Report struct {
	Date	time.Time		`json:"date" bson:"date"`
	Category	string		`json:"category" bson:"category"`
	Complaint string		`json:"complaint" bson:"complaint"`
	AuthorId	string	`json:"authorId" bson:"author_id"`
}

type VerifyRequest struct {
	AuthKey		string		`json:"authKey"`
}

type Tags struct {
	Tags []string		`json:"tags"`
}
