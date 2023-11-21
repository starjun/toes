package jobs

import "log"

type Job01 struct {
	// filtered
	Test string
}

// ReminderEmails.Run() will get triggered automatically.
func (e Job01) Run() {
	// Queries the DB
	// Sends some email
	log.Println("Every 10 sec send reminder emails \n", e.Test)
}
