package AppCenterClient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type AppList struct {
	Count int   `json:"count"`
	List  []App `json:"value"`
}

type App struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Os          string    `json:"os"`
	Platform    string    `json:"platform"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *AppCenterClient) ListApps() (list AppList, error error) {
	defer c.concurrencyUnlock()
	c.concurrencyLock()
	response, err := c.rest().R().Get(fmt.Sprintf("orgs/%v/apps", url.QueryEscape(*c.organization)))
	if err := c.checkResponse(response, err); err != nil {
		error = err
		return
	}

	var slist []App
	err = json.Unmarshal(response.Body(), &slist)
	if err != nil {
		error = err
		return
	}

	list.Count = len(slist)
	list.List = slist
	return
}
