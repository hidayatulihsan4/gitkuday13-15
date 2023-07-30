package main

import (
	"b48s1/connection"
	"context"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// untuk menyimpan
type Project struct {
	Id          int
	ProjectName string
	StartDate   time.Time
	EndDate     time.Time
	Duration    string
	Description string
	NodeJs      bool
	ReactJs     bool
	Golang      bool
	Javascript  bool
	Image       string
	Technology  []string
}

type User struct {
	id       int
	name     string
	email    string
	password string
}

type SessionData struct {
	isLogin bool
	name    string
}

var userData = SessionData{}

func main() {

	e := echo.New()
	connection.DatabaseConnect()

	// Mengatur penanganan file static(jss,css,gambar)
	e.Static("/public", "public")

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// Daftar Routes GET(digunakan untuk permintaan get)
	e.GET("/hello", helloWorld)
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/testimonial", testimonial)
	e.GET("/add-project", addProject)
	e.GET("/edit-project/:id", editProject)
	e.GET("/project/:id", projectDetail)
	e.GET("/register", register)
	e.GET("/login", login)
	e.GET("/logout", logout)

	//Daftar Routes POST(digunakan untuk permintaan post)
	e.POST("/submitregister", submitRegister)
	e.POST("/submitlogin", submitLogin)
	// e.POST("/", submitLogin)
	e.POST("/", submitProject)
	e.POST("/edit-project/:id", submitEditedProject)
	e.POST("/delete-project/:id", deleteProject)
	// e.POST("/login", logout)
	// e.POST("/submitlogout", submitlogout)

	// Server(akan mengirimkan pesan fatal dan menghentikan eksekusi program)
	e.Logger.Fatal(e.Start("localhost:8000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Worldl!")
}

func login(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func register(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, project_name, description, image, start_date, end_date, technology FROM tb_project")

	dataProjects := []Project{}
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.ProjectName, &each.Description, &each.Image, &each.StartDate, &each.EndDate, &each.Technology)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		each.Duration = countDuration(each.StartDate, each.EndDate)

		if checkValue(each.Technology, "ReactJs") {
			each.ReactJs = true
		}
		if checkValue(each.Technology, "Javascript") {
			each.Javascript = true
		}
		if checkValue(each.Technology, "Golang") {
			each.Golang = true
		}
		if checkValue(each.Technology, "NodeJs") {
			each.NodeJs = true
		}

		dataProjects = append(dataProjects, each)
	}

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.isLogin = false
	} else {
		userData.isLogin = session.Values["isLogin"].(bool)
		userData.name = session.Values["name"].(string)
	}
	projects := map[string]interface{}{
		"Projects":     dataProjects,
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}
	return tmpl.Execute(c.Response(), projects)

	// return tmpl.Execute(c.Response(), dataSession)
}

func addProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.isLogin = false
	} else {
		userData.isLogin = session.Values["isLogin"].(bool)
		userData.name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}

	return tmpl.Execute(c.Response(), dataSession)
}

func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.isLogin = false
	} else {
		userData.isLogin = session.Values["isLogin"].(bool)
		userData.name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}

	return tmpl.Execute(c.Response(), dataSession)
}

func testimonial(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimonial.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.isLogin = false
	} else {
		userData.isLogin = session.Values["isLogin"].(bool)
		userData.name = session.Values["name"].(string)
	}

	dataSession := map[string]interface{}{
		"dataSession":  userData,
		"FlashStatus":  session.Values["status"],
		"FlashMessage": session.Values["message"],
		"FlashName":    session.Values["name"],
	}

	return tmpl.Execute(c.Response(), dataSession)
}

func projectDetail(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	idToInt, _ := strconv.Atoi(id)

	ProjectDetail := Project{}

	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_project WHERE id=$1", idToInt).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.Description, &ProjectDetail.Image, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Technology)

	if errQuery != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ProjectDetail.Duration = countDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	if checkValue(ProjectDetail.Technology, "ReactJs") {
		ProjectDetail.ReactJs = true
	}
	if checkValue(ProjectDetail.Technology, "Javascript") {
		ProjectDetail.Javascript = true
	}
	if checkValue(ProjectDetail.Technology, "Golang") {
		ProjectDetail.Golang = true
	}
	if checkValue(ProjectDetail.Technology, "NodeJs") {
		ProjectDetail.NodeJs = true
	}

	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.isLogin = false
	} else {
		userData.isLogin = session.Values["isLogin"].(bool)
		userData.name = session.Values["name"].(string)
	}

	data := map[string]interface{}{
		"Id":              id,
		"Project":         ProjectDetail,
		"startDateString": ProjectDetail.StartDate.Format("2006-01-02"),
		"endDateString":   ProjectDetail.EndDate.Format("2006-01-02"),
		"dataSession":     userData,
		"FlashStatus":     session.Values["status"],
		"FlashMessage":    session.Values["message"],
		"FlashName":       session.Values["name"],

		// "startDateString": ProjectDetail.StartDate.Format("12-31-2002"),
		// "endDateString":   ProjectDetail.EndDate.Format("12-31-2002"),
	}

	return tmpl.Execute(c.Response(), data)
}

func editProject(c echo.Context) error {
	id := c.Param("id")

	tmpl, err := template.ParseFiles("views/edit-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Project Not Found"})
	}

	idToInt, _ := strconv.Atoi(id)

	ProjectDetail := Project{}

	errQuery := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_project WHERE id=$1", idToInt).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.Description, &ProjectDetail.Image, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Technology)

	if errQuery != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ProjectDetail.Duration = countDuration(ProjectDetail.StartDate, ProjectDetail.EndDate)

	if checkValue(ProjectDetail.Technology, "ReactJs") {
		ProjectDetail.ReactJs = true
	}
	if checkValue(ProjectDetail.Technology, "Javascript") {
		ProjectDetail.Javascript = true
	}
	if checkValue(ProjectDetail.Technology, "Golang") {
		ProjectDetail.Golang = true
	}
	if checkValue(ProjectDetail.Technology, "NodeJs") {
		ProjectDetail.NodeJs = true
	}
	session, _ := session.Get("session", c)

	if session.Values["isLogin"] != true {
		userData.isLogin = false
	} else {
		userData.isLogin = session.Values["isLogin"].(bool)
		userData.name = session.Values["name"].(string)
	}

	data := map[string]interface{}{
		"Id":              id,
		"Project":         ProjectDetail,
		"startDateString": ProjectDetail.StartDate.Format("2006-01-02"),
		"endDateString":   ProjectDetail.EndDate.Format("2006-01-02"),
		"dataSession":     userData,
		"FlashStatus":     session.Values["status"],
		"FlashMessage":    session.Values["message"],
		"FlashName":       session.Values["name"],
	}

	return tmpl.Execute(c.Response(), data)
}

func logout(c echo.Context) error {
	session, _ := session.Get("session", c)
	session.Options.MaxAge = -1
	session.Values["isLogin"] = false
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func submitProject(c echo.Context) error {

	title := c.FormValue("input-name")
	image := c.FormValue("input-image")
	startdate := c.FormValue("startDate")
	enddate := c.FormValue("endDate")
	content := c.FormValue("input-description")
	technoReactJs := c.FormValue("ReactJs")
	technoJavascript := c.FormValue("Javascript")
	technoGolang := c.FormValue("Golang")
	technoNodeJs := c.FormValue("NodeJs")

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (project_name, description, image, start_date, end_date, technology[1], technology[2], technology[3], technology[4]) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", title, content, image, startdate, enddate, technoReactJs, technoJavascript, technoGolang, technoNodeJs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	// fmt.Println(startdate)
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func submitEditedProject(c echo.Context) error {

	// Menangkap Id dari Query Params
	id := c.FormValue("id")
	title := c.FormValue("input-name")
	image := c.FormValue("input-image")
	startdate := c.FormValue("startDate")
	enddate := c.FormValue("endDate")
	content := c.FormValue("input-description")
	technoReactJs := c.FormValue("ReactJs")
	technoJavascript := c.FormValue("Javascript")
	technoGolang := c.FormValue("Golang")
	technoNodeJs := c.FormValue("NodeJs")

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET project_name=$1, description=$2, image=$7, start_date=$8, end_date=$9, technology[1]=$3, technology[2]=$4, technology[3]=$5, technology[4]=$6, WHERE id=$10", title, content, image, startdate, enddate, id, technoReactJs, technoJavascript, technoGolang, technoNodeJs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id := c.Param("id")

	idToInt, _ := strconv.Atoi(id)

	connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", idToInt)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func countDuration(d1 time.Time, d2 time.Time) string {

	diff := d2.Sub(d1)
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := days / 30

	if months >= 12 {
		return strconv.Itoa(months/12) + " tahun"
	}
	if months > 0 {
		return strconv.Itoa(months) + " bulan"
	}
	if weeks > 0 {
		return strconv.Itoa(weeks) + " minggu"
	}
	return strconv.Itoa(days) + " hari"
}

func checkValue(slice []string, object string) bool {
	for _, data := range slice {
		if data == object {
			return true
		}
	}
	return false
}
func submitRegister(c echo.Context) error {

	// Menangkap Id dari Query Params
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_akun (name, email, password) VALUES ($1, $2, $3)", name, email, password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusMovedPermanently, "/login")
}

func redirectMessage(c echo.Context, message string, status bool, path string) error {
	session, _ := session.Get("session", c)
	session.Values["message"] = message
	session.Values["status"] = status
	session.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, path)
}

func submitLogin(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	var user = User{}

	errEmail := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_akun WHERE email=$1", email).Scan(&user.id, &user.name, &user.email, &user.password)
	errPass := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_akun WHERE password=$1", password).Scan(&user.id, &user.name, &user.email, &user.password)

	if errEmail != nil {
		c.JSON(http.StatusInternalServerError, "Email or Password wrong!")
	}

	if errPass != nil {
		c.JSON(http.StatusInternalServerError, "Email or Password wrong!")
	}

	session, _ := session.Get("session", c)
	session.Options.MaxAge = 36000 //3satuan detik
	session.Values["message"] = "login Success"
	session.Values["status"] = true // show alert
	session.Values["name"] = user.name
	session.Values["id"] = user.id
	session.Values["isLogin"] = true // access login
	session.Save(c.Request(), c.Response())

	return redirectMessage(c, "Login Succes", true, "/")
}
