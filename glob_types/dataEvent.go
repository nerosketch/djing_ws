package glob_types

type EventCallback = func(v []byte)

type DataEvent struct {
	Callback EventCallback
}
