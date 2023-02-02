package org

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"pokergo/internal/org"
	"pokergo/internal/users"
	"pokergo/internal/webapi/binder"
	"pokergo/pkg/id"
)

type mux struct {
	orgAdapter  org.Adapter
	userAdapter users.Adapter
}

func NewMux(orgAdapter org.Adapter, userAdapter users.Adapter) *mux {
	return &mux{orgAdapter, userAdapter}
}

func (m *mux) Route(g *echo.Group) {
	g.POST("/newOrg", m.NewOrg)
	g.POST("/addToOrg", m.AddToOrg)
	g.GET("/listOrg", m.ListOrg)
}

func (m *mux) NewOrg(c echo.Context) error {
	data, bindErr := binder.BindRequest[newOrgRequest](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	o, err := m.orgAdapter.CreateOrg(data.Context(), data.UserID(), data.Request.Name)
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot create organization: %s", err))
	}

	return c.JSON(200, newOrgResponse{
		ID:   o.ID.Hex(),
		Name: o.Name,
	})
}

func (m *mux) AddToOrg(c echo.Context) error {
	data, bindErr := binder.BindRequest[addToOrgRequest](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	o, err := m.orgAdapter.GetOrgByName(data.Context(), data.Request.OrgName)
	if err != nil {
		if errors.Is(err, org.ErrOrgNotExists) {
			return c.String(404, "org not exists")
		}
		return c.String(500, fmt.Sprintf("cannot find org: %s", err.Error()))
	}

	usr, err := m.userAdapter.GetUserByName(data.Context(), data.Request.Who)
	if err != nil {
		if errors.Is(err, users.ErrUserNotExists) {
			return c.String(404, "user not exists")
		}
		return c.String(500, fmt.Sprintf("cannot perform the query: %s", err.Error()))
	}

	alreadyPresent := o.IsMember(usr.ID)
	if alreadyPresent {
		return c.String(400, "user already is a member of this org")
	}

	canAddMember := o.IsMember(data.UserID())
	if !canAddMember {
		return c.String(403, "a user is NOT a member of the organization")
	}

	if err := m.orgAdapter.AddToOrg(data.Context(), o.ID, usr.ID); err != nil {
		return c.String(500, "cannot add user to org")
	}

	return c.String(200, "ok")
}

func (m *mux) ListOrg(c echo.Context) error {
	data, bindErr := binder.BindRequest[listUserOrgRequest](c, true)
	if bindErr != nil {
		return c.String(bindErr.Code, bindErr.Message)
	}
	defer data.Cancel()

	orgs, err := m.orgAdapter.ListUserOrg(data.Context(), data.UserID())
	if err != nil {
		return c.String(500, fmt.Sprintf("cannot fetch data: %s", err.Error()))
	}

	var response []orgResponse
	for _, o := range orgs {
		members, err := m.userAdapter.UserDetails(data.Context(), o.Members)
		if err != nil {
			return c.String(500, fmt.Sprintf("cannot get org-details: %s", err.Error()))
		}

		response = append(response, orgResponse{
			ID:        o.ID.Hex(),
			Name:      o.Name,
			Admin:     members[o.ID].Username,
			Members:   idsAndNames(members),
			CreatedAt: o.CreatedAt,
		})
	}

	return c.JSON(200, listUserOrgResponse{response})
}

func idsAndNames(input map[id.ID]users.User) []idWithName {
	var res []idWithName
	for _, v := range input {
		res = append(res, idWithName{
			ID:   v.ID.Hex(),
			Name: v.Username,
		})
	}
	return res
}
