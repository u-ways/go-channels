package log

import (
	"fmt"
	"log"
)

func Dentist(action string) {
	log.SetFlags(log.Ltime)
	log.Printf(action, "Dentist")
}

func Patient(id int, action string) {
	log.SetFlags(log.Ltime)
	var patient = fmt.Sprintf("%s (%d)", "Patient", id)
	log.Printf(action, patient)
}