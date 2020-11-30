package main

import (
	"fmt"
	"log"
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

/*
 * assistant that communicates with the patients using the queues hwait and lwait,
 * and communicates with the dentist using one single queue wait. The dentist will
 * not see or act on the queues hwait and hwait but only receive patients on wait.
 */
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
					assistantLog(movingLPatientToHwait)
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
			assistantLog(placingAHighPriorityPatient)
			wait <- hPatient
		default:
			select {
			case lPatient := <-lwait:
				timer.Reset(limit)
				assistantLog(placingALowPriorityPatient)
				wait <- lPatient
			default:
				break
			}
		}
	}
}

/** dentist **********************************************************/

/**
 * The dentist. The dentist checks for patients in the waiting room.
 *   • If there are no patients, the dentist falls asleep.
 *   • If there are is at least one patient, the dentist calls the first one in.
 *     The remaining patients keep waiting. During the treatment, the dentist is
 *     active while the patient is sleeping1. When the dentist finishes the treatment,
 *     the patient is woken up, and the dentist checks for patients in the waiting room.
 *     And so on...
 */
func dentist(wait chan chan int, dent <-chan chan int, ready chan<- bool) {
	for {
		select {
		case nextPatient := <-wait:
			treat(nextPatient)
		default:
			// Sleep until a patient shows up and requests a treatment
			dentistLog(wentToSleep)
			ready <- signal
			newlyArrivedPatient := <-dent
			dentistLog(wakesUp)
			treat(newlyArrivedPatient)
		}
	}
}

/**
 * Emulates a treatment operation activity
 */
func treat(patient chan int) {
	dentistLog(startTreatingPatient)

	patient <- start
	// Emulate dentist treatment activity
	dentistTreatmentActivity()

	// Dentist making sure patient has shinny teeth
	dentistLog(checksPatientTeeth)
	patient <- qa

	// Handshake to acknowledge treatment is complete
	accept(<-patient, finish, getOffTheChair)
	patient <- finish
}

/**
 * The dentistTreatmentActivity is a time-consuming action (i.e. pausing
 * the current goroutine based on maximum and minimum "treatment" time.)
 */func dentistTreatmentActivity() {
	const minDuration = 1
	const maxDuration = 3

	random := rand.New(rand.NewSource(Now().UnixNano()))
	treatmentTime := Duration(minDuration+random.Intn(maxDuration)) * Second

	Sleep(treatmentTime)
}

/** patient **********************************************************/

/**
 * The patient. The patient, upon arrival, checks if the dentist
 * is busy with other patients or sleeping.
 *   • If the dentist is sleeping, the patient wakes the dentist up and falls
 *     asleep while being treated. The patient is woken up when the treatment
 *     is completed, and leaves (i.e., terminates).
 *   • If the dentist is busy with another patient, the arriving patient goes
 *     in the waiting room and waits (i.e., sleeps). When the patient is woken
 *     up, the treatment starts: the patient falls asleep until being woken up
 *     at the end of the treatment.
 */
func patient(wait chan<- chan int, dent chan<- chan int, id int) {
	patientLog(id, requestTreatment)

	// Creates an appointed treatment channel
	treatment := make(chan int)

	select {
	// Request treatment (wakes up the dentist if asleep)
	case dent <- treatment:
		patientLog(id, dentistNotBusy)
		receiveTreatment(id, treatment)
	default:
		// Dentist is busy, go to the waiting room and wait (i.e. sleep)
		wait <- treatment
		patientLog(id, waitingForTreatment)
		receiveTreatment(id, treatment)
	}

	close(treatment)
}

/**
 * Emulates receiving a treatment operation
 */
func receiveTreatment(id int, treatment chan int) {
	// Wait until you start receiving the treatment
	accept(<-treatment, start, treatmentMustBeInSync)

	// When start is received, dentist start the treatment
	patientLog(id, isGettingTreated)

	// Patient "sleeps" until operation is complete (i.e. gets blocked)
	accept(<-treatment, qa, treatmentMustBeInSync)

	// When qa is received, dentist asks the Patient to smile.
	patientLog(id, shineTeeth)

	treatment <- finish
	patientLog(id, leaveClinic)
	accept(<-treatment, finish, treatmentIsComplete)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Main Method                                                                                             //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

	accept(<-ready, signal, dentistIsNotReady)

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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Utils (code duplicated across files (instead of having a utils package) to adhere to submission format) //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

/** state **********************************************************/

/**
 * Available state codes the dentist and patient use for channel communication
 */
const start = 0
const qa = 1
const finish = 2

/**
 * Signal is used to indicate the dentist is "ready" to treat patients
 */
const signal = true

/**
 * A function to enforce consuming expected channel values
 */
func accept(operation interface{}, expected interface{}, msg string) {
	if operation != expected {
		panic(msg)
	}
}

/** loggers **********************************************************/

/**
 * A log function identifying assistant
 */
func assistantLog(action string) {
	log.SetFlags(log.Ltime)
	log.Printf(action, "Assistant")
}

/**
 * A log function identifying dentist
 */
func dentistLog(action string) {
	log.SetFlags(log.Ltime)
	log.Printf(action, "Dentist")
}

/**
 * A log function identifying patient
 */
func patientLog(id int, action string) {
	log.SetFlags(log.Ltime)
	var patient = fmt.Sprintf("%s (%d)", "Patient", id)
	log.Printf(action, patient)
}

/** colors **********************************************************/

/**
 * Colors used to make logs more readable
 */
var clear = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var cyan = "\033[36m"
var purple = "\033[35m"
var gray = "\033[37m"

func init() {
	if runtime.GOOS == "windows" {
		clear = ""
		red = ""
		green = ""
		yellow = ""
		blue = ""
		cyan = ""
		purple = ""
		gray = ""
	}
}

/** events **********************************************************/

// Dentist log events
var wentToSleep = yellow + "%s is sleeping. (no patients)" + clear
var wakesUp = yellow + "%s woke up." + clear
var dentistNotBusy = green + "%s will be treated right away. (Dentist is not busy)" + clear
var startTreatingPatient = green + "%s is treating the patient." + clear
var checksPatientTeeth = purple + "%s finished the surgery! Dentist checks patient teeth <=" + clear

// Patient log events
var requestTreatment = blue + "%s requested a treatment." + clear
var waitingForTreatment = red + "%s have to wait for treatment. (Dentist is not ready yet)" + clear
var isGettingTreated = yellow + "%s is getting treated. (They have been put to sleep until surgery is complete)" + clear
var shineTeeth = purple + "=> %s has shiny teeth!" + clear
var leaveClinic = gray + "%s is leaving the clinic." + clear

// Panic log events
var dentistIsNotReady = red + "Sorry, I am not ready yet..." + clear
var treatmentMustBeInSync = red + "Wait! Are you sure you're a dentist???" + clear
var treatmentIsComplete = red + "Aren't we finished? Can I leave please?" + clear
var getOffTheChair = red + "We're done here, can you get off the chair please?" + clear

// Assistant log events
var movingLPatientToHwait = cyan + "%s is moving one low priority patient to high priority." + clear
var placingAHighPriorityPatient = cyan + "%s placed a HIGH priority patient in the waiting area" + clear
var placingALowPriorityPatient = cyan + "%s placed a LOW priority patient in the waiting area" + clear
