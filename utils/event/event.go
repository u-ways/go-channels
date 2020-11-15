package event

import "ue21/utils/color"

var Sleep =               color.Yellow + "%s is sleeping. (no patients)" + color.Reset
var WakeUp =              color.Yellow + "%s woke up. Someone rang the clinic desk call bell!" + color.Reset
var RequestTreatment =    color.Blue + "%s requested a treatment." + color.Reset
var DentistNotBusy =      color.Green + "%s will be treated right away. (Dentist is not busy)" + color.Reset
var WaitingForTreatment = color.Red + "%s have to wait for treatment in the waiting room. (Dentist is busy)" + color.Reset
var Treat =               color.Green + "%s is treating a patient." + color.Reset
var GettingTreated =      color.Green + "%s is getting treated." + color.Reset
var CheckTeeth =          color.Purple + "%s checks patient teeth..." + color.Reset
var ShineTeeth =          color.Purple + "%s has shiny teeth!" + color.Reset
var SaysThankYou =        color.Purple + "%s thanks the dentist!" + color.Reset
var LeaveClinic =         color.Gray + "%s is leaving the clinic."  + color.Reset
var OutOfSync =           color.Red + "Wait! Are you sure you're a doctor???" + color.Reset
