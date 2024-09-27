package eventbus

import (
	"context"
	"root/internal/team/model"
)

type OrderRegisteredEvent struct {
	TeamName   string
	ResultChan chan Result // Канал для возврата результата
	OrderRole  string
	Context    context.Context
	Track      string
}

// Result содержит команду и ошибку для возврата
type Result struct {
	Team  *model.Team
	Error error
}
