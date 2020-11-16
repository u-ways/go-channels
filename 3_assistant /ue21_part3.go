package main

import (
	. "go-channels/utils/event"
	. "go-channels/utils/log"
	. "go-channels/utils/state"
	"math/rand"
	"runtime"
	. "time"
)

// 3.b.
//   "What would be a deadlock possibility that affects part 2 but not part 3?"
//
//   > Deadlock is where no process proceeds, and get blocked.
//
//   1. Dentist have one low priority patient
//   2. Dentist checks the high priority queue first => found no one
//   3. Low priority patient has aged => move to high priority queue
//   4. Dentist checks the low priority queue => found no one
//   5. Dentist goes to sleep... (deadlock)
//
//   This is fixed in part 3 because the dentist has one "waiting" queue only.

func assistant(hwait chan chan int, lwait <-chan chan int, wait chan<- chan int) {
	limit := 500 * Millisecond
	timer := NewTimer(limit)

	// Aging algorithm:
	// Move a patient from lwait to hwait whenever limit has passed
	go func() {
		for {
			select {
			case <-timer.C:
				select {
				case lPatient := <-lwait:
					Dentist(MovingLPatientToHwait)
					hwait <- lPatient
					timer.Reset(limit)
				}
			}
		}
	}()

	// Serve patients in the high priority queue first.
	// And age low priority patients by limit everytime hwait is read.
	for {
		select {
		case hPatient := <-hwait:
			Assistant(PlacingAHighPriorityPatient)
			wait <- hPatient
		default:
			select {
			case lPatient := <-lwait:
				timer.Reset(limit)
				Assistant(PlacingALowPriorityPatient)
				wait <- lPatient
			default:
				break
			}
		}
	}
}

/** dentist **********************************************************/

func dentist(wait chan chan int, dent <-chan chan int, ready chan<- bool) {
	for {
		select {
		case nextPatient := <-wait:
			treat(nextPatient)
		default:
			// Sleep until a patient shows up and requests a treatment
			Dentist(WentToSleep)
			ready <- READY
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
	runtime.GOMAXPROCS(maxThreads)

	// creates a synchronous channel
	ready := make(chan bool)
	dent := make(chan chan int)

	const lwaitChannelSize = 10
	const hwaitChannelSize = 20
	const waitChannelSize = 15

	// creates an asynchronous channels for hwait and lwait
	wait := make(chan chan int, waitChannelSize)
	lwait := make(chan chan int, lwaitChannelSize)
	hwait := make(chan chan int, hwaitChannelSize)

	go dentist(wait, dent, ready)
	go assistant(hwait, lwait, wait)

	Accept(<-ready, READY, DentistIsNotReady)

	const lPatients = 10
	const hPatients = 20

	for i := lPatients; i <= hPatients; i++ {
		go patient(hwait, dent, i)
	}

	for i := 1; i <= lPatients; i++ {
		go patient(lwait, dent, i)
	}

	Sleep(5 * (hPatients + lPatients) * Second)
}
