package main

import (
	"fmt"
	"github.com/the-friyia/go-affect/AuthenticationSystem"
	_ "github.com/the-friyia/go-affect/Memory"
	"github.com/the-friyia/go-affect/Model"
	"github.com/the-friyia/go-affect/AffectControlLib"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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
	"tmpl/fragments/pomodoro_activity_view.html"))

var globalSessions *session.Manager
var TestUser = &model.User{Goals: []model.Goal{}, Username: "Daniel", Password: "pass"}
var days = [5]string{"1 Monday", "2 Tuesday", "3 Wednesday", "4 Thursday", "5 Friday"}

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

type Page struct {
	Title    string
	Username string
	Body     []byte
	Goal     []string
	User	 *model.User
	NumOfGoals []int
	WeeklyGoals map[string][]model.Goal
	Days	[5]string
	FirstGoal string
	PomodoroTime int
	Breaktime int
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

func createUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key: ", k)
		fmt.Println("Value: ", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello World")
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

func loadGoalInformation(usr *model.User) map[string][]model.Goal {
	days := [5]string{"1.) Monday", "2.) Tuesday", "3.) Wednesday", "4.) Thursday", "5.) Friday"}
	WeeklyGoals := make(map[string][]model.Goal)

	if len(usr.Goals) > 5 {
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
	return WeeklyGoals
}

func addGoalsToUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for _, v := range r.Form {
		if v[0] != "" {
			TestUser.Goals = append(TestUser.Goals, model.Goal{GoalName: v[0], Priority: 0})
		}
	}


	f := loadGoalInformation(TestUser)["1.) Monday"][0]
	p := &Page{WeeklyGoals: loadGoalInformation(TestUser), Days: days, FirstGoal: f.GoalName, NumOfGoals: []int{1, 2, 3}, PomodoroTime: 10, Breaktime: 5}

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
		if r.Form["username"][0] == TestUser.Username && r.Form["password"][0] == TestUser.Password {
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

func pomodoroUpdate(w http.ResponseWriter, r *http.Request) {

	f := loadGoalInformation(TestUser)["1.) Monday"][0]
	array := make([]int, affect.Deflection())
	for i := range array {
		array[i] = i+1
	}

	p := &Page{WeeklyGoals: loadGoalInformation(TestUser),
		 	   Days: days,
			   FirstGoal: f.GoalName,
			   NumOfGoals: array,
			   PomodoroTime: affect.PomodoroTime(),
		   	   Breaktime: affect.BreakTime()}

	renderTemplate(w, "pomodoro_action_view", p)
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", createUser)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/pomodoro", addGoalsToUser)
	http.HandleFunc("/pomodoro-update", pomodoroUpdate)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe(":8080", nil)
}
