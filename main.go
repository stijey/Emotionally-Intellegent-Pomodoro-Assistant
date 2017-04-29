package main

import (
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
	"regexp"
	"crypto/sha256"
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
		fmt.Println(rows)
		var uid int
		var username string
		var password string
		var blob *map[string][]model.Goal
		err = rows.Scan(&uid, &username, &password, &blob)
		checkErr(err)
		TestUser = &model.User{Username: username, Password: password, WeeklyGoals: blob}
	}
	fmt.Println(TestUser.WeeklyGoals)

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
		loginSetGoals(w, r, username)
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
	err = db.QueryRow("INSERT INTO users(username,password) VALUES($1,$2) returning uid;", username, hashedPassword).Scan(&lastInsertId)

	if err != nil {
		return false;
	}
	return true;
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("data/"+filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile("data/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+p.Title, http.StatusFound)
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

var validPath = regexp.MustCompile("^/(index|edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
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

func loadGoalInformation(usr *model.User) *map[string][]model.Goal {
	days := [5]string{"1.) Monday", "2.) Tuesday",
		"3.) Wednesday", "4.) Thursday", "5.) Friday"}

	WeeklyGoals := make(map[string][]model.Goal)

	if usr.Goals[5].GoalName != "" {
		numOfGoals := len(usr.Goals)
		for day := range days {
			if numOfGoals >= 2 {
				WeeklyGoals[days[day]] = []model.Goal{usr.Goals[day], usr.Goals[day+1]}
				numOfGoals -= 2
			} else {
				if numOfGoals >= 0 {
					WeeklyGoals[days[day]] = []model.Goal{usr.Goals[day]}
					numOfGoals--
				} else {
					WeeklyGoals[days[day]] = []model.Goal{model.Goal{GoalName: ""}}
				}
			}
		}
	} else {
		for day := range days {
			WeeklyGoals[days[day]] = []model.Goal{usr.Goals[day]}
		}
	}
	return &WeeklyGoals
}

func saveGoals(TestUser *model.User) bool {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	var lastInsertId int
	err = db.QueryRow("UPDATE users SET weekly_goals=$1 WHERE username=$2 returning uid;", *TestUser.WeeklyGoals, TestUser.Username).Scan(&lastInsertId)
	fmt.Println(err)
	if err.Error() != "" {
		return true;
	} else {
		return false;
	}
}

func addGoalsToUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g1"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g2"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g3"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g4"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g5"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g6"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g7"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g8"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g9"][0]})
	TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: r.Form["g10"][0]})

	fmt.Println(TestUser.Goals)

	f := (*loadGoalInformation(TestUser))["1.) Monday"][0]

	p := &Page{WeeklyGoals: (*loadGoalInformation(TestUser)),
		Days:         days,
		FirstGoal:    f.GoalName,
		NumOfGoals:   []int{1, 2, 3},
		PomodoroTime: 10,
		Breaktime:    5}

	TestUser.WeeklyGoals = loadGoalInformation(TestUser)
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

func adjustNumberOfGoalsForTheWeek() {
	temp := getWeeklyGoalsFromSession()
	firstDayGoals := temp["1.) Monday"]
	copy(firstDayGoals[0:], firstDayGoals[1:])
	firstDayGoals[len(firstDayGoals)-1] = model.Goal{}
	firstDayGoals = firstDayGoals[:len(firstDayGoals)-1]
	temp["1.) Monday"] = firstDayGoals
	setWeeklyGoalsInSession(temp)
}

func pomodoroUpdate(w http.ResponseWriter, r *http.Request) {

	array := make([]int, affect.Deflection())
	for i := range array {
		array[i] = i + 1
	}

	r.ParseForm()

	if r.Form.Get("goal-complete") == "true" {
		adjustNumberOfGoalsForTheWeek()
	}

	f := getWeeklyGoalsFromSession()["1.) Monday"][0]

	p := &Page{WeeklyGoals: getWeeklyGoalsFromSession(),
		Days:         days,
		FirstGoal:    f.GoalName,
		NumOfGoals:   array,
		PomodoroTime: affect.PomodoroTime(),
		Breaktime:    affect.BreakTime()}

	renderTemplate(w, "pomodoro_action_view", p)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/createaccount", signupHandler)
	// http.HandleFunc("/signup", createUser)
	http.HandleFunc("/signup", prepareUserFormData)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/pomodoro", addGoalsToUser)
	http.HandleFunc("/pomodoro-update", pomodoroUpdate)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}

// package main
//
// import (
//     "database/sql"
//     "fmt"
//     _ "github.com/lib/pq"
//     "time"
// )
//
// const (
//     DB_USER     = "thefriyia"
//     DB_PASSWORD = ""
//     DB_NAME     = "temp"
// )
//
// func main() {
//     dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
//         DB_USER, DB_PASSWORD, DB_NAME)
//     db, err := sql.Open("postgres", dbinfo)
//     checkErr(err)
//     defer db.Close()
//
//     fmt.Println("# Inserting values")
//
//     var lastInsertId int
//     err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
//     checkErr(err)
//     fmt.Println("last inserted id =", lastInsertId)
//
//     fmt.Println("# Updating")
//     stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
//     checkErr(err)
//
//     res, err := stmt.Exec("astaxieupdate", lastInsertId)
//     checkErr(err)
//
//     affect, err := res.RowsAffected()
//     checkErr(err)
//
//     fmt.Println(affect, "rows changed")
//
//     fmt.Println("# Querying")
//     rows, err := db.Query("SELECT * FROM userinfo")
//     checkErr(err)
//
//     for rows.Next() {
//         var uid int
//         var username string
//         var department string
//         var created time.Time
//         err = rows.Scan(&uid, &username, &department, &created)
//         checkErr(err)
//         fmt.Println("uid | username | department | created ")
//         fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
//     }
//
//     fmt.Println("# Deleting")
//     stmt, err = db.Prepare("delete from userinfo where uid=$1")
//     checkErr(err)
//
//     res, err = stmt.Exec(lastInsertId)
//     checkErr(err)
//
//     affect, err = res.RowsAffected()
//     checkErr(err)
//
//     fmt.Println(affect, "rows changed")
// }
//
// func checkErr(err error) {
//     if err != nil {
//         panic(err)
//     }
// }
