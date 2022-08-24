package input

import "errors"

type (
	UserInput struct {
		Login string `url:"login,omitempty"`
		Id    string `url:"id,omitempty"`
	}
)

func (u UserInput) Validate() error {
	if u.Login != "" && u.Id != "" {
		return errors.New("you should set id either login")
	}

	if u.Login == "" && u.Id == "" {
		return errors.New("null exception")
	}

	return nil
}
