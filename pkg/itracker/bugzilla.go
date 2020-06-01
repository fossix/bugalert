package itracker

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Bug struct {
	*Bugzilla
	sync.Mutex
	Blocks              []int         `json:"blocks"`
	IsCcAccessible      bool          `json:"is_cc_accessible"`
	Keywords            []string      `json:"keywords"`
	URL                 string        `json:"url"`
	QaContactID         string        `json:"qa_contact"`
	UpdateToken         string        `json:"update_token"`
	CcDetail            []User        `json:"cc_detail"`
	Summary             string        `json:"summary"`
	Platform            string        `json:"platform"`
	Version             string        `json:"version"`
	Deadline            interface{}   `json:"deadline"`
	IsCreatorAccessible bool          `json:"is_creator_accessible"`
	IsConfirmed         bool          `json:"is_confirmed"`
	Priority            string        `json:"priority"`
	AssignedTo          User          `json:"assigned_to_detail"`
	CreatorID           string        `json:"creator"`
	LastChangeTime      time.Time     `json:"last_change_time"`
	Creator             User          `json:"creator_detail"`
	Cc                  []string      `json:"cc"`
	SeeAlso             []interface{} `json:"see_also"`
	Groups              []interface{} `json:"groups"`
	AssignedToID        string        `json:"assigned_to"`
	CreationTime        time.Time     `json:"creation_time"`
	Whiteboard          string        `json:"whiteboard"`
	ID                  int           `json:"id"`
	DependsOn           []int         `json:"depends_on"`
	DupeOf              int           `json:"dupe_of"`
	QaContact           User          `json:"qa_contact_detail"`
	Resolution          string        `json:"resolution"`
	Classification      string        `json:"classification"`
	Alias               []interface{} `json:"alias"`
	OpSys               string        `json:"op_sys"`
	Status              string        `json:"status"`
	IsOpen              bool          `json:"is_open"`
	Severity            string        `json:"severity"`
	Flags               []Flag        `json:"flags"`
	Component           string        `json:"component"`
	TargetMilestone     string        `json:"target_milestone"`
	Product             string        `json:"product"`
	History             []*History    `json:"history"`
	// custom fields
	CustomFields interface{}
}

type Flag struct {
	TypeID           int       `json:"type_id"`
	ModificationDate time.Time `json:"modification_date"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	ID               int       `json:"id"`
	Setter           string    `json:"setter"`
	Requestee        string    `json:"requestee"`
	CreationDate     time.Time `json:"creation_date"`
}

type History struct {
	When    time.Time `json:"when"`
	Who     string    `json:"who"`
	Changes []struct {
		Added     string `json:"added"`
		FieldName string `json:"field_name"`
		Removed   string `json:"removed"`
	} `json:"changes"`
}

type HistoryList struct {
	History []History `json:"history"`
}

type Bugzilla struct {
	url      string
	endpoint string
	apikey   string
}

func NewBugzilla(url, endpoint string) *Bugzilla {
	bz := &Bugzilla{
		url:      url,
		endpoint: endpoint,
	}
	return bz
}

func (b *Bugzilla) SetAPIKey(key string) {
	b.apikey = key
}

func (b *Bugzilla) SetRestEndPoint(endpoint string) {
	b.endpoint = endpoint
}

func (b *Bugzilla) get(api string, args map[string]string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", b.url, b.endpoint, api)

	if args == nil {
		args = make(map[string]string)
	}
	args["api_key"] = b.apikey

	return get(endpoint, args)
}

func (b *Bugzilla) GetBugs(args map[string]string) ([]*Bug, error) {
	body, err := b.get("bug", args)
	if err != nil {
		return nil, err
	}

	type _bugs struct {
		Bugs   []*Bug        `json:"bugs"`
		Faults []interface{} `json:"faults"`
	}

	var _b _bugs
	if err := json.Unmarshal(body, &_b); err != nil {
		return nil, err
	}

	if len(_b.Bugs) == 0 {
		return nil, fmt.Errorf("No bugs found")
	}

	for _, bug := range _b.Bugs {
		bug.Bugzilla = b
		go bug.GetHistory()
	}

	return _b.Bugs, nil
}

func (b *Bugzilla) GetBug(id int) (*Bug, error) {
	args, err := makeargs([]string{"id"}, []string{fmt.Sprintf("%d", id)})
	if err != nil {
		return nil, err
	}

	bugs, err := b.GetBugs(args)
	if err != nil {
		return nil, err
	}

	if len(bugs) > 1 {
		return nil, fmt.Errorf("Unexpected output, expected 1, got %d", len(bugs))
	}

	return bugs[0], nil
}

func (bug *Bug) GetAssignee() (*User, error) {
	u := &bug.AssignedTo
	if u == nil {
		return nil, fmt.Errorf("Bug not assigned yet")
	}

	u.Bugzilla = bug.Bugzilla
	return u, nil
}

func (b *Bugzilla) GetUser(id string) (*User, error) {
	endpoint := fmt.Sprintf("user/%s", id)
	body, err := b.get(endpoint, nil)
	if err != nil {
		return nil, err
	}

	var u struct {
		Users []User `json:"users"`
	}

	if err := json.Unmarshal(body, &u); err != nil {
		return nil, err
	}

	if len(u.Users) == 0 {
		return nil, fmt.Errorf("No such user")
	}
	if len(u.Users) > 1 {
		return nil, fmt.Errorf("Unexpected output, expected 1, got %d", len(u.Users))
	}

	u.Users[0].Bugzilla = b
	return &u.Users[0], nil
}

func (bug *Bug) GetHistory() error {
	bug.Lock()
	defer bug.Unlock()

	// don't fetch if already available
	if len(bug.History) > 0 {
		return nil
	}

	endpoint := fmt.Sprintf("bug/%d/history", bug.ID)
	body, err := bug.get(endpoint, nil)
	if err != nil {
		return err
	}

	type _bugs struct {
		Bugs []struct {
			History []*History `json:"history"`
		} `json:"bugs"`
		Alias string `json:"alias"`
		ID    int    `json:"id"`
	}

	var bugs _bugs

	if err := json.Unmarshal(body, &bugs); err != nil {
		return err
	}

	if len(bugs.Bugs) == 0 {
		return fmt.Errorf("Cannot find history for bug")
	}
	if len(bugs.Bugs) > 1 {
		return fmt.Errorf("Unexpected output, expected 1, got %d", len(bugs.Bugs))
	}

	t := bugs.Bugs[0]
	for _, h := range t.History {
		bug.History = append(bug.History, h)
	}

	return nil
}
