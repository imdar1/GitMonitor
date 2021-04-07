package observer

type Observer interface {
	Update()
}

func NotifyObserver(o Observer) {
	o.Update()
}
