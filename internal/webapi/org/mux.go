package org

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"pokergo/internal/org"
	"pokergo/internal/users"
	"pokergo/internal/webapi"
	"pokergo/pkg/id"
	"time"
)

type mux struct {
	orgAdapter  org.Adapter
	userAdapter users.Adapter
}

func NewMux(orgAdapter org.Adapter, userAdapter users.Adapter) *mux {
	return &mux{orgAdapter: orgAdapter, userAdapter: userAdapter}
}

func (m *mux) Route(e *echo.Echo, prefix string) error {
	e.POST(prefix+"/newOrg", m.NewOrg)
	e.POST(prefix+"/addToOrg", m.AddToOrg)
	e.GET(prefix+"/listOrg", m.ListOrg)
	return nil
}

func (m *mux) NewOrg(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	jwtToken, err := webapi.GetJWTToken(c)
	if err != nil {
		return c.String(403, "jwt token invalid")
	}

	// no need to verify if the user exists - jwt already did the job

	var request newOrgRequest
	if err := c.Bind(&request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err))
	}

	uID := id.FromString(jwtToken.ID)
	o, err := m.orgAdapter.CreateOrg(reqCtx, uID, request.Name)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot create organization: %s", err))
	}

	return c.JSON(200, newOrgResponse{
		ID:   o.ID.Hex(),
		Name: o.Name,
	})
}

func (m *mux) AddToOrg(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	jwtToken, err := webapi.GetJWTToken(c)
	if err != nil {
		return c.String(403, "jwt token invalid")
	}
	requesterID := id.FromString(jwtToken.ID)

	var request addToOrgRequest
	if err := c.Bind(&request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err))
	}

	o, err := m.orgAdapter.GetOrgByName(reqCtx, request.OrgName)
	if err != nil {
		if errors.Is(err, org.ErrOrgNotExists) {
			return c.String(404, fmt.Sprintf("org not exists"))
		}
		return c.String(500, fmt.Sprintf("cannot find org: %s", err.Error()))
	}

	usr, err := m.userAdapter.GetUserByName(reqCtx, request.Who)
	if err != nil {
		if errors.Is(err, users.ErrUserNotExists) {
			return c.String(404, fmt.Sprintf("user not exists"))
		}
		return c.String(500, fmt.Sprintf("cannot perform the query: %s", err.Error()))
	}

	isMember := false
	for _, oo := range o.Members {
		if oo == requesterID {
			isMember = true
		}
		if oo == usr.ID {
			return c.String(400, "user already is a member of this org")
		}
	}
	if !isMember {
		return c.String(403, fmt.Sprintf("a user is NOT a member of the organization"))
	}

	if err := m.orgAdapter.AddToOrg(reqCtx, o.ID, usr.ID); err != nil {
		return c.String(500, fmt.Sprintf("cannot add user to org"))
	}

	return c.String(200, "ok")
}

func (m *mux) ListOrg(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	jwtToken, err := webapi.GetJWTToken(c)
	if err != nil {
		return c.String(403, "jwt token invalid")
	}
	requesterID := id.FromString(jwtToken.ID)

	var request listUserOrgRequest
	if err := c.Bind(&request); err != nil {
		return c.String(400, fmt.Sprintf("invalid request: %s", err.Error()))
	}

	orgs, err := m.orgAdapter.ListUserOrg(reqCtx, requesterID)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot fetch data: %s", err.Error()))
	}

	var response []orgResponse
	for _, o := range orgs {
		members, err := m.userAdapter.UserDetails(reqCtx, o.Members)
		if err != nil {
			return c.String(500, fmt.Sprintf("cannot get org-details: %s", err.Error()))
		}

		response = append(response, orgResponse{
			ID:        o.ID.Hex(),
			Name:      o.Name,
			Admin:     members[o.ID].Username,
			Members:   idAndNames(members),
			CreatedAt: o.CreatedAt,
		})
	}

	return c.JSON(200, listUserOrgResponse{response})
}

func idAndNames(input map[id.ID]users.User) []idWithName {
	var res []idWithName
	for _, v := range input {
		res = append(res, idWithName{
			ID:   v.ID.Hex(),
			Name: v.Username,
		})
	}
	return res
}
