package ws

type ClientEventCallback = func(v []byte)

type ClientDataEvent struct {
	Callback ClientEventCallback
}
