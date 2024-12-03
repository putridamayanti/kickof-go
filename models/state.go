package models

type State struct {
	Id          string    `json:"id"`
	WorkspaceId string    `json:"workspaceId" bson:"workspaceId"`
	ProjectId   string    `json:"projectId" bson:"projectId"`
	Name        string    `json:"name"`
	Workspace   Workspace `json:"workspace" bson:"-"`
	Project     Project   `json:"project" bson:"-"`
	Tasks       []Task    `json:"tasks" bson:"-"`
	BasicDate   `bson:",inline"`
}
