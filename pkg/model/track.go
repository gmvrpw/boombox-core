package model

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Icon string `json:"icon"`
}

type Track struct {
	Url     string  `json:"url"`
	Name    string  `json:"name"`
	Author  Author  `json:"author"  gorm:"embedded;embeddedPrefix:author_"`
	Cover   string  `json:"cover"`
	Service Service `json:"service" gorm:"embedded;embeddedPrefix:service_"`
}

type UnplayableTrackError struct{}

func (e *UnplayableTrackError) Error() string {
	return "cannot find runner"
}

type UnspecifiedRequestError struct {
	Options []*Request
}

func (e *UnspecifiedRequestError) Error() string {
	return "too many tracks"
}
