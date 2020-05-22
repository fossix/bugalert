package itracker

type User struct {
	*Bugzilla
	Email    string        `json:"email"`
	RealName string        `json:"real_name"`
	Name     string        `json:"name"`
	Groups   []interface{} `json:"groups"`
	CanLogin bool          `json:"can_login"`
	ID       int           `json:"id"`
}

func (u *User) Bugs() ([]Bug, error) {
	args, err := makeargs([]string{"assigned_to"}, []string{u.Email})
	if err != nil {
		return nil, err
	}
	bugs, err := u.GetBugs(args)
	if err != nil {
		return nil, err
	}

	return bugs, nil
}
