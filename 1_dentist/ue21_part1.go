package main

import (
	"math/rand"
	"runtime"
	"time"
	"ue21/utils/event"
	"ue21/utils/log"
	"ue21/utils/state"
)

/** dentist **********************************************************/

func dentist(wait <-chan chan int, dent <-chan chan int) {
	for {
		select {
		case nextPatient := <-wait:
			// Treat the FIFO patient in the waiting room
			treat(nextPatient)
		default:
			// Sleep when no patients found in the waiting room
			log.Dentist(event.Sleep)
			// But wake up when a patient shows up and requests a treatment
			patientArrived := <-dent
			log.Dentist(event.WakeUp)
			// Treat the newly arrived patient
			treat(patientArrived)
		}
	}
}

func treat(patient chan<- int) {
	const minDuration = 1
	const maxDuration = 3

	log.Dentist(event.Treat)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	treatmentTime := time.Duration(minDuration+random.Intn(maxDuration)) * time.Second

	patient <- state.START
	// Dentist activity (treatment time)
	time.Sleep(treatmentTime)

	patient <- state.QA
	// Dentist making sure patient has shinny teeth
	log.Dentist(event.CheckTeeth)

	patient <- state.FINISH
	// Let the dentist rest for a bit
	time.Sleep(time.Second)
}

/** patient **********************************************************/

func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	log.Patient(id, event.RequestTreatment)

	// Create an appointed treatment channel
	treatment := make(chan int)

	select {
	// Request treatment (wakes up the dentist if asleep)
	case dent <- treatment:
		log.Patient(id, event.DentistNotBusy)
		receiveTreatment(id, treatment)
	default:
		// Dentist is busy, go to the waiting room and wait (i.e. sleep)
		wait <- treatment
		log.Patient(id, event.WaitingForTreatment)
		// Wait until you start receiving the treatment
		receiveTreatment(id, treatment)
	}

	close(treatment)
	log.Patient(id, event.LeaveClinic)
}

func receiveTreatment(id int, treatment <-chan int) {
	accept(<-treatment, state.START)
	log.Patient(id, event.GettingTreated)

	accept(<-treatment, state.QA)
	log.Patient(id, event.ShineTeeth)

	accept(<-treatment, state.FINISH)
	log.Patient(id, event.SaysThankYou)
}

func accept(operation int, expected int) {
	if operation != expected {
		panic(event.OutOfSync)
	}
}

func main() {
	const maxThreads = 10
	const numberOfPatients = 20
	const channelSize = 5

	runtime.GOMAXPROCS(maxThreads)

	// creates a synchronous channel
	dent := make(chan chan int)
	// creates an asynchronous channel of size `channelSize`
	wait := make(chan chan int, channelSize)

	go dentist(wait, dent)

	time.Sleep(2 * time.Second)

	for i := 1; i <= numberOfPatients; i++ {
		go patient(wait, dent, i)
		time.Sleep(time.Second)
	}

	time.Sleep(numberOfPatients * 3 * time.Second)
}
