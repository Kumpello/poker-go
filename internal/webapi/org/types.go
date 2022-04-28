package org

import "pokergo/internal/org"

type newOrgRequest struct {
	Name string `json:"name"`
}

type newOrgResponse struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

type addToOrgRequest struct {
	OrgName string `json:"name"`
	Who     string `json:"who"`
}

type listUserOrgRequest struct {
	// empty
}

type listUserOrgResponse struct {
	Orgs []org.Org `json:"orgs"`
}
