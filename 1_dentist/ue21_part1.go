package main

import (
	. "go-channels/utils/event"
	. "go-channels/utils/log"
	. "go-channels/utils/state"
	"math/rand"
	"runtime"
	. "time"
)

/** dentist **********************************************************/

func dentist(wait <-chan chan int, dent <-chan chan int) {
	for {
		select {
		case nextPatient := <-wait:
			treat(nextPatient)
		default:
			// Sleep when no patients found in the waiting room
			Dentist(WentToSleep)
			// But wake up when a patient shows up and requests a treatment
			newlyArrivedPatient := <-dent
			Dentist(WakesUp)
			treat(newlyArrivedPatient)
		}
	}
}

func treat(patient chan int) {
	Dentist(StartTreatingPatient)

	patient <- START
	// Emulate dentist treatment activity
	dentistTreatmentActivity()

	// Dentist making sure patient has shinny teeth
	Dentist(ChecksPatientTeeth)
	patient <- QA

	// Handshake to acknowledge treatment is complete
	Accept(<-patient, FINISH, GetOffTheChair)
	patient <- FINISH
}

// dentistTreatmentActivity is implemented as a time-consuming action (i.e. pausing
// the current goroutine based on maximum and minimum "treatment" time.)
func dentistTreatmentActivity() {
	const minDuration = 1
	const maxDuration = 3

	random := rand.New(rand.NewSource(Now().UnixNano()))
	treatmentTime := Duration(minDuration+random.Intn(maxDuration)) * Second

	Sleep(treatmentTime)
}

/** patient **********************************************************/

func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	Patient(id, RequestTreatment)

	// Creates an appointed treatment channel
	treatment := make(chan int)

	select {
	// Request treatment (wakes up the dentist if asleep)
	case dent <- treatment:
		Patient(id, DentistNotBusy)
		receiveTreatment(id, treatment)
	default:
		// Dentist is busy, go to the waiting room and wait (i.e. sleep)
		wait <- treatment
		Patient(id, WaitingForTreatment)
		receiveTreatment(id, treatment)
	}

	close(treatment)
}

func receiveTreatment(id int, treatment chan int) {
	// Wait until you start receiving the treatment
	Accept(<-treatment, START, TreatmentMustBeInSync)

	// When START is received, dentist start the treatment
	Patient(id, IsGettingTreated)

	// Patient "sleeps" until operation is complete (i.e. gets blocked)
	Accept(<-treatment, QA, TreatmentMustBeInSync)

	// When QA is received, dentist asks the Patient to smile.
	Patient(id, ShineTeeth)

	treatment <- FINISH
	Patient(id, LeaveClinic)
	Accept(<-treatment, FINISH, TreatmentIsComplete)
}

func main() {
	const maxThreads = 5
	const numberOfPatients = 10
	const channelSize = 5

	runtime.GOMAXPROCS(maxThreads)

	// creates a synchronous channel
	dent := make(chan chan int)
	// creates an asynchronous channel of size `channelSize`
	wait := make(chan chan int, channelSize)

	go dentist(wait, dent)

	Sleep(2 * Second)

	for i := 1; i <= numberOfPatients; i++ {
		go patient(wait, dent, i)
		Sleep(Second)
	}

	Sleep(3 * numberOfPatients * Second)
}
