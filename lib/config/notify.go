package config

type Notify struct {
	Signal chan struct{}
}

func NewNotify() Notify {
	return Notify{Signal: make(chan struct{})}
}
