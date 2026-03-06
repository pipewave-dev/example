
/*
package main


source-subscriber

type Listener struct {
	c Connection
	unsubFn map[SourceName]func() // Key is source name, value is unsubscribe function
}

func Subscribe[DataT any](listener Listener, source SourceEvent[DataT], handler func(data DataT)) {
	unsub := source.Subscribe(handler)
	listener.unsubFn[source.Name()] = unsub
}

type SourceName string // Wrapper type for source name to easily lookup

type SourceEvent[DataT any] interface {
  Name() SourceName
	Emit(data DataT)
	Subscribe(func(data DataT)) (unsubscribe func())
}

*/
