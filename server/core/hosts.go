package core

/*
	Sliver Implant Framework
	Copyright (C) 2021  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	consts "github.com/bishopfox/sliver/client/constants"
	"github.com/bishopfox/sliver/server/db"
	"github.com/bishopfox/sliver/server/db/models"
	"github.com/bishopfox/sliver/server/log"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

var (
	coreLog = log.NamedLogger("core", "hosts")
)

func StartEventAutomation() {
	go func() {
		for event := range EventBroker.Subscribe() {
			switch event.EventType {

			case consts.SessionOpenedEvent:
				if event.Session != nil {
					hostsSessionCallback(event.Session)
				}
			}

		}
	}()
}

// Triggered on new sessione events, checks to see if the host is in
// the database and adds it if not.
func hostsSessionCallback(session *Session) {
	dbSession := db.Session()
	_, err := db.HostByHostUUID(session.UUID)
	if err != nil && err != gorm.ErrRecordNotFound {
		coreLog.Error(err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err := dbSession.Create(&models.Host{
			HostUUID:      uuid.FromStringOrNil(session.UUID),
			Hostname:      session.Hostname,
			OSVersion:     session.Os,
			IOCs:          []models.IOC{},
			ExtensionData: []models.ExtensionData{},
		}).Error
		if err != nil {
			coreLog.Error(err)
			return
		}
	}
}