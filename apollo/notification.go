package apollo

import (
	"encoding/json"
	"sync"
)

const defaultNotificationID int = -1

type notificationsMgr struct {
	notifications sync.Map
}

type notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationID int    `json:"notificationId"`
}

func newNotificationManager(namespaceNames []string) *notificationsMgr {
	n := &notificationsMgr{
		notifications: sync.Map{},
	}
	for _, namespaceName := range namespaceNames {
		n.notifications.Store(namespaceName, defaultNotificationID)
	}
	return n
}

func (n *notificationsMgr) String() string {
	var notifications []*notification
	n.notifications.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		v, _ := value.(int)
		notifications = append(notifications, &notification{
			NamespaceName:  k,
			NotificationID: v,
		})
		return true
	})
	res, err := json.Marshal(&notifications)
	if err != nil {
		return ""
	}
	return string(res)
}

func (n *notificationsMgr) Store(namespaceName string, notificationID int) {
	n.notifications.Store(namespaceName, notificationID)
}
