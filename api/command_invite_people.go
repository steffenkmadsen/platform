// Copyright (c) 2016 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"strings"

	"github.com/mattermost/platform/app"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

type InvitePeopleProvider struct {
}

const (
	CMD_INVITE_PEOPLE = "invite_people"
)

func init() {
	RegisterCommandProvider(&InvitePeopleProvider{})
}

func (me *InvitePeopleProvider) GetTrigger() string {
	return CMD_INVITE_PEOPLE
}

func (me *InvitePeopleProvider) GetCommand(c *Context) *model.Command {
	return &model.Command{
		Trigger:          CMD_INVITE_PEOPLE,
		AutoComplete:     true,
		AutoCompleteDesc: c.T("api.command.invite_people.desc"),
		AutoCompleteHint: c.T("api.command.invite_people.hint"),
		DisplayName:      c.T("api.command.invite_people.name"),
	}
}

func (me *InvitePeopleProvider) DoCommand(c *Context, args *model.CommandArgs, message string) *model.CommandResponse {
	if !utils.Cfg.EmailSettings.SendEmailNotifications {
		return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command.invite_people.email_off")}
	}

	emailList := strings.Fields(message)

	for i := len(emailList) - 1; i >= 0; i-- {
		emailList[i] = strings.Trim(emailList[i], ",")
		if !strings.Contains(emailList[i], "@") {
			emailList = append(emailList[:i], emailList[i+1:]...)
		}
	}

	if len(emailList) == 0 {
		return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command.invite_people.no_email")}
	}

	if err := app.InviteNewUsersToTeam(emailList, c.TeamId, c.Session.UserId, c.GetSiteURL()); err != nil {
		c.Err = err
		return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command.invite_people.fail")}
	}

	return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command.invite_people.sent")}
}
