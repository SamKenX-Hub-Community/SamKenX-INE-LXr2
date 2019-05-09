// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/snap"
)

// The maximum number of Usb bInterfaceNumber.
const UsbMaxInterfaces = 32

// AppLabelExpr returns the specification of the apparmor label describing
// all the apps bound to a given slot. The result has one of three forms,
// depending on how apps are bound to the slot:
//
// - "snap.$snap_instance.$app" if there is exactly one app bound
// - "snap.$snap_instance.{$app1,...$appN}" if there are some, but not all, apps bound
// - "snap.$snap_instance.*" if all apps are bound to the slot
func appLabelExpr(apps map[string]*snap.AppInfo, snap *snap.Info) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `"snap.%s.`, snap.InstanceName())
	if len(apps) == 1 {
		for appName := range apps {
			buf.WriteString(appName)
		}
	} else if len(apps) == len(snap.Apps) {
		buf.WriteByte('*')
	} else {
		appNames := make([]string, 0, len(apps))
		for appName := range apps {
			appNames = append(appNames, appName)
		}
		sort.Strings(appNames)
		buf.WriteByte('{')
		for _, appName := range appNames {
			buf.WriteString(appName)
			buf.WriteByte(',')
		}
		buf.Truncate(buf.Len() - 1)
		buf.WriteByte('}')
	}
	buf.WriteByte('"')
	return buf.String()
}

func slotAppLabelExpr(slot *interfaces.ConnectedSlot) string {
	return appLabelExpr(slot.Apps(), slot.Snap())
}

func plugAppLabelExpr(plug *interfaces.ConnectedPlug) string {
	return appLabelExpr(plug.Apps(), plug.Snap())
}
