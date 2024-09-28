package dto

import "root/internal/team/model"

type TeamFilterParam struct {
	Name string `json:"name"`
}

type UpdateTeamReq struct {
	TeamName string      `json:"team_name,omitempty"`
	Link     string      `json:"link,omitempty"`
	Track    model.Track `json:"track,omitempty"`
}
