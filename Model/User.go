package model

type User struct {
	Goals    []Goal
	Username string
	Password string
	WeeklyGoals *map[string][]Goal
}

type Goal struct {
	GoalName string
	Priority int
}

type Goals []Goal

func (slice Goals) Len() int {
    return len(slice)
}

func (slice Goals) Less(i, j int) bool {
    return slice[i].Priority < slice[j].Priority;
}

func (slice Goals) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

type PomodoroRound struct {
	Goals    [3]Goal
	Duration int32
}

type PomodoroSession struct {
	PomodoroRounds []PomodoroRound
	Days           [5]string
}
