package dto


type RegisterReq struct {
	FIO         string     `json:"fio" validate:"required"`
	Age         int        `json:"age" validate:"required"`
	//Role        model.Role `json:"role" validate:"required"`
	//PhoneNumber string     `json:"phone_number" validate:"required,e164"`
	Email       string     `json:"email" validate:"required,email"`
	// TeamName    string     `json:"team_name" validate:"required,alpha"`
}
