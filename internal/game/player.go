package game

import "pokergo/pkg/id"

type Player struct {
	UserID id.ID `bson:"userID"`

	// BuyIn in subunits (for example 1zł20gr=120gr)
	BuyIn int64 `bson:"startStack" json:"startStack"`
	// BuyOut is a stack the player has at the end.
	BuyOut *int64 `bson:"finishStack,omitempty" json:"finishStack,omitempty"`

	// AdditionalIncomes is a list of inGameTransaction
	AdditionalIncomes []inGameTransaction `bson:"additionalIncomes,omitempty" json:"additionalIncomes,omitempty"`
}
