package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Harits514/AdroadyTes/app/model"
	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "be_test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Println("connection error")
		panic(err.Error())
	}
	return db
}

func InsertData(w http.ResponseWriter, r *http.Request) {
	drv := model.Drivers{}
	db := dbConn()

	drv.Name = r.FormValue("Name")
	drv.Email = r.FormValue("Email")
	drv.PhoneNumber = r.FormValue("PhoneNumber")

	file, handler, err := r.FormFile("UploadFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	drv.ImageLink = "../AdroadyTes/image/" + handler.Filename

	f, err := os.OpenFile("../AdroadyTes/image/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	insForm, err := db.Prepare("INSERT INTO drivers(Name, Email, PhoneNumber, ImageLink) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(drv.Name, drv.Email, drv.PhoneNumber, drv.ImageLink)
	defer db.Close()
	w.WriteHeader(http.StatusCreated)
	return
}

func UpdateData(w http.ResponseWriter, r *http.Request) {
	drv := model.Drivers{}
	db := dbConn()
	drv.Name = r.FormValue("Name")
	drv.Email = r.FormValue("Email")
	drv.PhoneNumber = r.FormValue("PhoneNumber")
	temp := r.FormValue("id")
	drv.ID, _ = strconv.Atoi(temp)
	_, err := strconv.Atoi(temp)
	if err != nil {
		log.Println("ID is not a number")
		panic(err.Error())
	}

	insForm, err := db.Prepare("UPDATE drivers SET Name=?, Email=?, PhoneNumber=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(drv.Name, drv.Email, drv.PhoneNumber, drv.ID)
	defer db.Close()
	w.WriteHeader(http.StatusCreated)
	return
}

func ViewPhoto(w http.ResponseWriter, r *http.Request) {
	drv := model.Drivers{}
	db := dbConn()
	temp, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println("ID is not a number")
		panic(err.Error())
	}
	drv.ID = temp

	selDB, err := db.Query("SELECT ImageLink FROM drivers WHERE id=?", temp)
	var link string
	for selDB.Next() {
		err := selDB.Scan(&link)
		log.Println(link)
		if err != nil {
			log.Println("FileNotFound")
			panic(err.Error())
		}
	}

	img, err := os.Open(link)
	img2, _, _ := image.Decode(img)
	w.Header().Set("Content-Type", "image/jpeg")
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img2, nil); err != nil {
		log.Println("unable to encode image.")
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	drv := model.Drivers{}
	db := dbConn()
	temp, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println("ID is not a number")
		panic(err.Error())
	}
	drv.ID = temp

	selDB, err := db.Query("SELECT ImageLink FROM drivers WHERE id=?", temp)
	var link string
	for selDB.Next() {
		err := selDB.Scan(&link)
		log.Println(link)
		if err != nil {
			log.Println("FileNotFound")
			panic(err.Error())
		}
		os.Remove(link)
	}

	file, handler, err := r.FormFile("UploadFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", handler.Header)

	drv.ImageLink = "../AdroadyTes/image/" + handler.Filename

	f, err := os.OpenFile(drv.ImageLink, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	insForm, err := db.Prepare("UPDATE drivers SET ImageLink=? WHERE id=?")
	if err != nil {
		log.Println("data not inserted")
		panic(err.Error())
	}
	insForm.Exec(drv.ImageLink, drv.ID)
	defer db.Close()
	w.WriteHeader(http.StatusCreated)
	return
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	temp, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println("ID is not a number")
		panic(err.Error())
	}

	selDB, err := db.Query("SELECT ImageLink FROM drivers WHERE id=?", temp)
	var link string
	for selDB.Next() {
		err := selDB.Scan(&link)
		log.Println(link)
		if err != nil {
			log.Println("FileNotFound")
			panic(err.Error())
		}
		os.Remove(link)
	}

	insForm, err := db.Prepare("DELETE FROM drivers WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(temp)
	defer db.Close()
	w.WriteHeader(http.StatusOK)
	return
}

func LoadData(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	temp, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println("ID is not a number")
		panic(err.Error())
	}
	selDB, err := db.Query("SELECT * FROM drivers WHERE id = ?", temp)
	if err != nil {
		log.Println(err)
		http.Error(w, "data retrieval failed", http.StatusInternalServerError)
		return
	}
	emp := model.Drivers{}
	res := []model.Drivers{}
	for selDB.Next() {
		var id int
		var nama, email, linkphoto, nohp string
		err = selDB.Scan(&id, &nama, &email, &nohp, &linkphoto)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = nama
		emp.Email = email
		emp.PhoneNumber = nohp
		emp.ImageLink = linkphoto
		res = append(res, emp)
		log.Println(res)
	}
	usersJSON, err := json.Marshal(res)
	w.Write(usersJSON)
}
