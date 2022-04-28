package organization

import (
	"context"
	"pokergo/pkg/id"
)

type Organization struct {
	ID      id.ID   `bson:"_id" json:"id"`
	Admin   id.ID   `json:"admin"`
	Members []id.ID `bson:"members" json:"members"`
}

type Adapter interface {
	CreateOrganization(ctx context.Context, admin id.ID) (Organization, error)
	AddToOrganization(ctx context.Context, who []id.ID) error
}
