package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/the-friyia/go-affect/AffectControlLib"
	"github.com/the-friyia/go-affect/AuthenticationSystem"
	_ "github.com/the-friyia/go-affect/Memory"
	"github.com/the-friyia/go-affect/Model"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
)

var templates = template.Must(template.ParseFiles(
	"tmpl/edit.html",
	"tmpl/view.html",
	"tmpl/index.html",
	"tmpl/initial_goals.html",
	"tmpl/pomodoro_action_view.html",
	"tmpl/fragments/header.html",
	"tmpl/fragments/footer.html",
	"tmpl/fragments/login.html",
	"tmpl/fragments/signup.html",
	"tmpl/fragments/login_failure.html",
	"tmpl/fragments/weekly_goals.html",
	"tmpl/create_user.html",
	"tmpl/fragments/pomodoro_activity_view.html"))

var globalSessions *session.Manager
var TestUser *model.User
var days = [5]string{"1 Monday", "2 Tuesday", "3 Wednesday",
	"4 Thursday", "5 Friday"}

var affectiveState affect.AffectiveState

const (
	DB_USER     = "thefriyia"
	DB_PASSWORD = ""
	DB_NAME     = "test"
)

type Page struct {
	Title        string
	Username     string
	Body         []byte
	Goal         []string
	User         *model.User
	NumOfGoals   []int
	WeeklyGoals  map[string][]model.Goal
	Days         [5]string
	FirstGoal    string
	PomodoroTime int
	Breaktime    int
}

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	affectiveState = affect.MakeAffectiveState()
	go globalSessions.GC()
}

func authenticateUser(username string, password string) bool {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE username=($1)", username)

	for rows.Next() {
		// fmt.Println(rows)
		var uid int
		var username string
		var password string
		var blob *map[string][]model.Goal
		err = rows.Scan(&uid, &username, &password, &blob)
		checkErr(err)
		TestUser = &model.User{Username: username,
			Password:    password,
			WeeklyGoals: blob}
	}

	if TestUser.Password == hashPassword(password) {
		return true
	}
	return false
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Welcome", Goal: nil}
	renderTemplate(w, "create_user", p)
}

func prepareUserFormData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	passwordConfirm := r.Form["password-confirm"][0]
	err := createUser(username, password, passwordConfirm)

	if !err {
		p := &Page{Title: "Welcome", Goal: nil}
		renderTemplate(w, "create_user", p)
	} else {
		indexHandler(w, r)
	}
}

func hashPassword(password string) string {
	h := sha256.New()
	passwordBytes := []byte(password)
	passwordHashed := h.Sum(passwordBytes)
	return string(passwordHashed)
}

func createUser(username string, password string, passwordConfirm string) bool {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	hashedPassword := hashPassword(password)
	hashedPasswordConfirm := hashPassword(passwordConfirm)

	if hashedPassword != hashedPasswordConfirm {
		return false
	}

	var lastInsertId int

	err = db.QueryRow("INSERT INTO users(username,password) VALUES($1,$2) re"+
		"turning uid;", username, hashedPassword).Scan(&lastInsertId)

	if err != nil {
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("data/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Welcome", Goal: nil}
	renderTemplate(w, "index", p)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginOrdinary(w http.ResponseWriter, r *http.Request, username string) {
	p := &Page{}

	if username == "" {
		p = &Page{Title: "Welcome", Username: ""}
	} else {
		p = &Page{Title: "Welcome", Username: username}
	}
	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveGoals(TestUser *model.User) bool {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	var lastInsertId int
	err = db.QueryRow("UPDATE users SET weekly_goals=$1 WHERE username=$2 re"+
		"turning uid;", *TestUser.WeeklyGoals, TestUser.Username).
		Scan(&lastInsertId)

	fmt.Println(err)
	if err.Error() != "" {
		return true
	} else {
		return false
	}
}

func delegateGoals(goals []model.Goal, w http.ResponseWriter,
	r *http.Request) *map[string][]model.Goal {

	days := [5]string{"1.) Monday", "2.) Tuesday",
		"3.) Wednesday", "4.) Thursday", "5.) Friday"}

	WeeklyGoals := make(map[string][]model.Goal)

	if goals[0].GoalName == "" {
		renderTemplate(w, "pomodoro_action_view", &Page{})
	} else if goals[1].GoalName == "" || goals[2].GoalName == "" ||
		goals[3].GoalName == "" || goals[4].GoalName == "" ||
		goals[5].GoalName == "" {

		for i := 0; i < 5; i++ {
			WeeklyGoals[days[i]] = []model.Goal{goals[i]}
		}
	} else if goals[6].GoalName == "" {
		WeeklyGoals[days[0]] = []model.Goal{goals[0], goals[1]}
		for i := 1; i < 5; i++ {
			WeeklyGoals[days[i]] = []model.Goal{goals[i+1]}
		}
	} else if goals[7].GoalName == "" {
		WeeklyGoals[days[0]] = []model.Goal{goals[0], goals[1]}
		WeeklyGoals[days[1]] = []model.Goal{goals[2], goals[3]}
		for i := 2; i < 5; i++ {
			WeeklyGoals[days[i]] = []model.Goal{goals[i+2]}
		}
	} else if goals[8].GoalName == "" {
		WeeklyGoals[days[0]] = []model.Goal{goals[0], goals[1]}
		WeeklyGoals[days[1]] = []model.Goal{goals[2], goals[3]}
		WeeklyGoals[days[2]] = []model.Goal{goals[4], goals[5]}
		for i := 3; i < 5; i++ {
			WeeklyGoals[days[i]] = []model.Goal{goals[i+3]}
		}
	} else if goals[9].GoalName == "" {
		WeeklyGoals[days[0]] = []model.Goal{goals[0], goals[1]}
		WeeklyGoals[days[1]] = []model.Goal{goals[2], goals[3]}
		WeeklyGoals[days[2]] = []model.Goal{goals[4], goals[5]}
		WeeklyGoals[days[3]] = []model.Goal{goals[6], goals[7]}
		for i := 4; i < 5; i++ {
			WeeklyGoals[days[i]] = []model.Goal{goals[i+4]}
		}
	} else if goals[9].GoalName != "" {
		j := 0
		for i := 0; i < 5; i++ {
			WeeklyGoals[days[i]] = []model.Goal{goals[j], goals[j+1]}
			j += 2
		}
	}

	return &WeeklyGoals
}

func addGoalsToUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if len(TestUser.Goals) != 0 {
		TestUser.Goals = []model.Goal{}
	}

	for i := 1; i <= 10; i++ {
		get := fmt.Sprintf("g%d", i)
		TestUser.Goals = append(TestUser.Goals,
			model.Goal{GoalName: r.Form[get][0], Priority: i})
	}

	f := *delegateGoals(TestUser.Goals, w, r)

	p := &Page{WeeklyGoals: (f),
		Days:         days,
		FirstGoal:    f["1.) Monday"][0].GoalName,
		NumOfGoals:   []int{1, 2, 3},
		PomodoroTime: 30,
		Breaktime:    5}

	TestUser.WeeklyGoals = &f
	saveGoals(TestUser)
	setWeeklyGoalsInSession(*TestUser.WeeklyGoals)
	renderTemplate(w, "pomodoro_action_view", p)
}

func loginSetGoals(w http.ResponseWriter, r *http.Request, username string) {
	p := &Page{Title: username}
	err := templates.ExecuteTemplate(w, "initial_goals.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		if authenticateUser(r.Form["username"][0], r.Form["password"][0]) {
			sess.Set(r.Form["username"][0], r.Form["username"])
			if len(TestUser.Goals) == 0 {
				loginSetGoals(w, r, r.Form["username"][0])
				return
			}
			loginOrdinary(w, r, r.Form["username"][0])
		} else {
			loginOrdinary(w, r, "FAILURE")
		}
	}
}

func setWeeklyGoalsInSession(value map[string][]model.Goal) {
	val := globalSessions.Provider
	aval, _ := val.SessionRead("username")
	aval.Set("weekly-goals", value)
}

func getWeeklyGoalsFromSession() map[string][]model.Goal {
	val := globalSessions.Provider
	aval, _ := val.SessionRead("username")
	return aval.Get("weekly-goals").(map[string][]model.Goal)
}

func adjustNumberOfGoalsForTheWeek(w http.ResponseWriter, r *http.Request) {
	temp := getWeeklyGoalsFromSession()
	var newArr model.Goals

	for _, v := range temp {
		for val := range v {
			newArr = append(newArr, v[val])
		}
	}

	fmt.Println()

	for i := len(newArr); i < 10; i++ {
		newArr = append(newArr, model.Goal{GoalName: "",
			Priority: 999})
	}

	sort.Sort(newArr)

	for i := 0; i < len(newArr)-1; i++ {
		newArr[i] = newArr[i+1]
	}

	newArr[len(newArr)-1] = model.Goal{GoalName: "", Priority: 999}
	sort.Sort(newArr)
	setWeeklyGoalsInSession(*delegateGoals(newArr, w, r))
}

func calculateRoundPerameters(userInput string) (goals int, breaks int,
	workTime int) {

	affectiveState.PropegateForward(affectiveState.UserInputToEPA(userInput))

	if affectiveState.Deflection > 21 {
		goals = 1
		breaks = 20
		workTime = 40
	} else if affectiveState.Deflection > 10 {
		goals = 2
		breaks = 15
		workTime = 30
	} else {
		goals = 3
		breaks = 10
		workTime = 20
	}
	affectiveState.Respond()

	return goals, breaks, workTime
}

func pomodoroUpdate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	goals, breakTime, workTime :=
		calculateRoundPerameters(r.Form["feedback-text"][0])

	array := make([]int, goals)
	for i := range array {
		array[i] = i + 1
	}

	if r.Form.Get("goal-complete") == "true" {
		adjustNumberOfGoalsForTheWeek(w, r)
	}

	var f model.Goal

	if len(getWeeklyGoalsFromSession()["1.) Monday"]) == 0 {
		f.GoalName = ""
		http.Redirect(w, r, "/setnewgoals", 302)
		return
	} else {
		f = getWeeklyGoalsFromSession()["1.) Monday"][0]
	}

	p := &Page{WeeklyGoals: getWeeklyGoalsFromSession(),
		Days:         days,
		FirstGoal:    f.GoalName,
		NumOfGoals:   array,
		PomodoroTime: workTime,
		Breaktime:    breakTime}

	renderTemplate(w, "pomodoro_action_view", p)
}

func setNewGoalsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "initial_goals", &Page{})
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/createaccount", signupHandler)
	http.HandleFunc("/signup", prepareUserFormData)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/setnewgoals", setNewGoalsHandler)
	http.HandleFunc("/pomodoro", addGoalsToUser)
	http.HandleFunc("/pomodoro-update", pomodoroUpdate)
	http.Handle("/assets/", http.StripPrefix("/assets/",
		http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}
