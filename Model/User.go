package model

type User struct {
	Goals    []Goal
	Username string
	Password string
}

type Goal struct {
	GoalName string
	Priority int32
}

type PomodoroRound struct {
	Goals    [3]Goal
	Duration int32
}

type PomodoroSession struct {
	PomodoroRounds []PomodoroRound
	Days           [5]string
}
