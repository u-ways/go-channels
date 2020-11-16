package event

import . "go-channels/utils/color"

// Dentist
var WentToSleep = Yellow + "%s is sleeping. (no patients)" + Reset
var WakesUp = Yellow + "%s woke up." + Reset
var DentistNotBusy = Green + "%s will be treated right away. (Dentist is not busy)" + Reset
var StartTreatingPatient = Green + "%s is treating the patient." + Reset
var ChecksPatientTeeth = Purple + "%s finished the surgery! Dentist checks patient teeth <=" + Reset

// Patient
var RequestTreatment = Blue + "%s requested a treatment." + Reset
var WaitingForTreatment = Red + "%s have to wait for treatment. (Dentist is not ready yet)" + Reset
var IsGettingTreated = Yellow + "%s is getting treated. (They have been put to sleep until surgery is complete)" + Reset
var ShineTeeth = Purple + "=> %s has shiny teeth!" + Reset
var LeaveClinic = Gray + "%s is leaving the clinic." + Reset

// Panic
var DentistIsNotReady = Red + "Sorry, I am not ready yet..." + Reset
var TreatmentMustBeInSync = Red + "Wait! Are you sure you're a dentist???" + Reset
var TreatmentIsComplete = Red + "Aren't we finished? Can I leave please?" + Reset
var GetOffTheChair = Red + "We're done here, can you get off the chair please?" + Reset

// Assistant
var FoundAHighPriorityPatient = Cyan + "%s found a high priority patient." + Reset
var FoundALowPriorityPatient = Cyan + "%s found a low priority patient while no high priority patients were available." + Reset
var MovingLPatientToHwait = Cyan + "%s is moving one low priority patient to high priority." + Reset

var PlacingAHighPriorityPatient = Cyan + "%s placed a HIGH priority patient in the waiting area" + Reset
var PlacingALowPriorityPatient = Cyan + "%s placed a LOW priority patient in the waiting area" + Reset
