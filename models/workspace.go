package models

type Workspace struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	Endpoint  string   `json:"endpoint"`
	Size      string   `json:"size"`
	UserIds   []string `json:"userIds"`
	Members   []User   `json:"members" bson:"-"`
	BasicDate `bson:",inline"`
}
