package event

import . "ue21/utils/color"

var GoesToSleep = 			Yellow + "%s is sleeping. (no patients)" + Reset
var WakesUp = 				Yellow + "%s woke up." + Reset
var RequestTreatment = 		Blue + "%s requested a treatment." + Reset
var DentistNotBusy = 		Green + "%s will be treated right away. (Dentist is not busy; I rang the clinic desk call bell to wake the dentist up!)" + Reset
var WaitingForTreatment = 	Red + "%s have to wait for treatment in the waiting room. (Dentist is busy treating another patient)" + Reset
var StartTreatingPatient = 	Green + "%s is treating the patient." + Reset
var IsGettingTreated = 		Yellow + "%s is getting treated. (They have been put to sleep until surgery is complete)" + Reset
var ChecksPatientTeeth = 	Purple + "%s finished the surgery! Dentist checks patient teeth <=" + Reset
var ShineTeeth = 			Purple + "=> %s has shiny teeth!" + Reset
var LeaveClinic = 			Gray + "%s is leaving the clinic." + Reset

var TreatmentMustBeInSync = Red + "Wait! Are you sure you're a dentist???" + Reset
var TreatmentIsComplete = 	Red + "Aren't we finished? Can I leave please?" + Reset
var GetOffTheChair = 		Red + "We're done here, can you get off the chair please?" + Reset
