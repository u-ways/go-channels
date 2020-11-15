package state

const START  = 0
const QA     = 1
const FINISH = 2

func Accept(operation int, expected int, msg string) {
	if operation != expected {
		panic(msg)
	}
}