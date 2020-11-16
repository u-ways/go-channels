package state

const START = 0
const QA = 1
const FINISH = 2

const READY = true

func Accept(operation interface{}, expected interface{}, msg string) {
	if operation != expected {
		panic(msg)
	}
}
