package tracker

type User struct {
	*Bugzilla
	Email    string        `json:"email"`
	RealName string        `json:"real_name"`
	Name     string        `json:"name"`
	Groups   []interface{} `json:"groups"`
	CanLogin bool          `json:"can_login"`
	ID       int           `json:"id"`
}

func (u *User) Bugs(filter map[string]string) ([]*Bug, error) {
	var _filter map[string]string

	if filter == nil {
		_filter = make(map[string]string)
	} else {
		_filter = filter
	}

	_filter["assigned_to"] = u.Email
	bugs, err := u.Search(_filter)
	if err != nil {
		return nil, err
	}

	return bugs, nil
}
