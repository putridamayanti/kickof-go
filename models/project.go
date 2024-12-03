package models

type Project struct {
	Id          string    `json:"id"`
	WorkspaceId string    `json:"workspaceId" bson:"workspaceId"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	UserIds     []string  `json:"userIds"`
	Members     []User    `json:"members" bson:"-"`
	Workspace   Workspace `json:"workspace" bson:"-"`
	BasicDate   `bson:",inline"`
}
