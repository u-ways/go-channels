package main

import (
	. "go-channels/utils/event"
	. "go-channels/utils/log"
	. "go-channels/utils/state"
	"math/rand"
	"runtime"
	. "time"
)

// 2.b.
//   "Can you identify one possibility of starvation in the problem statement of part 1?"
//
//   If we assume go does **not** have fair semantics, i.e. say the channels are not FIFO:
//
//   1. The dentist picks a high priority patient X from the waiting queue to treat.
//   2. And a soon as the dentist is done patient X requests another treatment.
//   3. Repeat 2.
//   4. Other patients will have a lower chance to get treated and thus starvation.
//
//   > Starvation can happen by uncontrolled resource management.
//
//   Looking at one of the problem statement points of part 1:
//
//   > - If the dentist is busy with another patient, the arriving patient
//       goes in the waiting room and waits (i.e., sleeps).
//
//   This means as long as dentist is busy with another patient, other patients
//   will starve until dentist is done with said patient.
//
//   ---
//
//   "If you identify one, say what you could change your solution to make the system starvation-free?"
//
//   If I cannot enforce higher-level synchronization through channels and communication, I would use
//   basic synchronization primitives such as mutual exclusion locks to enforce fairness.

/** dentist **********************************************************/

func dentist(hwait chan chan int, lwait <-chan chan int, dent <-chan chan int) {
	limit := 3000 * Millisecond
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

	for {
		select {
		case hPatient := <-hwait:
			Dentist(FoundAHighPriorityPatient)
			treat(hPatient)
		default:
			select {
			case lPatient := <-lwait:
				timer.Reset(limit)
				Dentist(FoundALowPriorityPatient)
				treat(lPatient)
			default:
				// Sleep until a patient shows up and requests a treatment
				Dentist(WentToSleep)
				newlyArrivedPatient := <-dent
				Dentist(WakesUp)
				treat(newlyArrivedPatient)
			}
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
	dent := make(chan chan int)

	const lwaitChannelSize = 5
	const hwaitChannelSize = 50

	// creates an asynchronous channels for hwait and lwait
	lwait := make(chan chan int, lwaitChannelSize)
	hwait := make(chan chan int, hwaitChannelSize)

	go dentist(hwait, lwait, dent)

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
