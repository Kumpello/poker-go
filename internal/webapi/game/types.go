package game

type createGameRequest struct {
	Org string `json:"org" validate:"required"`
}

type createGameResponse struct {
	ID string `json:"id"`
}

type appendPlayerRequest struct {
	GameID     string  `json:"game_id" validate:"required"`
	UserID     *string `json:"user_id"`
	UserName   string  `json:"user_name" validate:"required"`
	StartStack *int64  `json:"start_stack" validate:"required"` // ptr to allow 0 value
}

type setFinishStack struct {
	GameID      string `json:"game_id" validate:"required"`
	UserName    string `json:"user_name" validate:"required"`
	FinishStack *int64 `json:"finish_stack" validate:"required"` // ptr to allow 0 value
}

type reBuyIn struct {
	GameID   string `json:"game_id" validate:"required"`
	UserName string `json:"user_name" validate:"required"`
	BuyIn    int64  `json:"buy_in" validate:"required"`
}

type reBuyInFromPlayer struct {
	GameID   string `json:"game_id" validate:"required"`
	UserName string `json:"user_name" validate:"required"`
	FromName string `json:"from_name" validate:"required"`
	BuyIn    int64  `json:"buy_in" validate:"required"`
}
