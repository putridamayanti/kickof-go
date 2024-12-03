package models

import "time"

type TaskLabel struct {
	Id          string `json:"id"`
	WorkspaceId string `json:"workspaceId" bson:"workspaceId"`
	ProjectId   string `json:"projectId" bson:"projectId"`
	Label       string `json:"label"`
	Color       string `json:"color"`
	BasicDate   `bson:",inline"`
}

type Task struct {
	Id          string      `json:"id"`
	WorkspaceId string      `json:"workspaceId" bson:"workspaceId"`
	ProjectId   string      `json:"projectId" bson:"projectId"`
	StateId     string      `json:"stateId" bson:"stateId"`
	Title       string      `json:"title"`
	Code        string      `json:"description"`
	StartDate   time.Time   `json:"startDate" bson:"startDate"`
	EndDate     time.Time   `json:"endDate" bson:"endDate"`
	LabelIds    []string    `json:"labelIds" bson:"labelIds"` // Task Label ids
	Labels      []TaskLabel `json:"labels" bson:"-"`
	AssigneeIds []string    `json:"assigneeIds"`
	Assignees   []User      `json:"assignees" bson:"-"` // User id
	BasicDate   `bson:",inline"`
}
