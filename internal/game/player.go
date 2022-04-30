package game

import "pokergo/pkg/id"

type Player struct {
	// UserID may be nil for anonymous players
	UserID *id.ID `bson:"user_id,omitempty"`
	// UserName a player's name
	UserName string `bson:"user_name"`

	// BuyIn is total buy-in (including re-buy-ins) in subunits (for example 1zł20gr=120gr)
	BuyIn int64 `bson:"start_stack"`
	// BuyOut is a stack the player has at the end.
	BuyOut *int64 `bson:"finish_stack,omitempty"`

	// AdditionalIncomes is a list of inGameTransaction
	AdditionalIncomes []inGameTransaction `bson:"additional_incomes"`
}

// inGameTransaction special transaction between players
// works like re-buy-in, but they are from another player
// (increases buy out of that player)
type inGameTransaction struct {
	Amount   int64  `bson:"amount"`
	Reason   string `bson:"reason"`
	From     *id.ID `bson:"from_id,omitempty"`
	FromName string `bson:"from_name"`
}
