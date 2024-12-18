package runtime

type RunnerSession struct {
	ID       string `json:"id"`
	Url      string `json:"url"`
	Playback uint64 `json:"playback"`
	Port     int    `json:"port"`

	Stop chan bool `json:"-"`
}
