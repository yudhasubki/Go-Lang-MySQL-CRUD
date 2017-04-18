package main

import (
	"database/sql"
	"log"           // Display messages to console
	"net/http"      // Manage URL
	"html/template" // Manage HTML files
	_ "github.com/go-sql-driver/mysql" // MySQL Database driver
)

type Mahasiswa struct{
	Id int
	Nama string
	Kelas string
	Jurusan string
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func con() (db *sql.DB) {
	dbDriver := "mysql"   // Database driver
	dbUser := "root"      // Mysql username
	dbPass := "" // Mysql password
	dbName := "belajar_go"   // Mysql schema

	// Realize the connection with mysql driver
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	// If error stop the application
	checkErr(err)
	return db
}

func dir(){
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}

func main() {
	log.Println("You Started Server http://localhost:8080")
	dir()
	//Route Mahasiswa
	http.HandleFunc("/mahasiswa/", index)
	http.HandleFunc("/mahasiswa/tambah", TambahMahasiswa)
	http.HandleFunc("/mahasiswa/save", SimpanSiswa)
	http.HandleFunc("/mahasiswa/edit", EditSiswa)
	http.HandleFunc("/mahasiswa/update", UpdateSiswa)
	http.HandleFunc("/mahasiswa/hapus", DeleteSiswa)
	http.ListenAndServe(":8080",nil)
}

func index(w http.ResponseWriter , r *http.Request) {	
	koneksi := con()
	query, err := koneksi.Query("SELECT * FROM mahasiswa")
	checkErr(err)
	//Dijadiin Array
	n := Mahasiswa{}
	//di slice lalu di append
	data := []Mahasiswa{}

	for query.Next() {
		var id_mahasiswa int
		var nama , kelas , jurusan string

		//Scan harus sesuai Column pada table database
		err = query.Scan(&id_mahasiswa , &nama , &kelas , &jurusan)
		checkErr(err)
		n.Id = id_mahasiswa
		n.Nama = nama
		n.Kelas = kelas
		n.Jurusan = jurusan
		data = append(data,n)
	}
	log.Println(data)
	//Menghitung Query yang masuk
	//log.Println(len(data))
	tmpl.ExecuteTemplate(w,"index",data)
}

func TambahMahasiswa(w http.ResponseWriter , r *http.Request) {
	tmpl.ExecuteTemplate(w, "tambahMahasiswa", nil)
}

func SimpanSiswa(w http.ResponseWriter , r *http.Request){
	koneksi := con()
	if(r.Method == "POST"){
		r.ParseForm()
		query , err := koneksi.Prepare("INSERT mahasiswa SET nama=? , kelas=? , jurusan=?")
		checkErr(err)
		result , err := query.Exec(r.PostFormValue("nama"),r.PostFormValue("kelas"),r.PostFormValue("jurusan"))
		checkErr(err)
		result.LastInsertId()
		http.Redirect(w,r,"/mahasiswa",301)
	}
}

func EditSiswa(w http.ResponseWriter , r *http.Request){
	koneksi := con()
	id := r.URL.Query().Get("id")
	query , err := koneksi.Query("SELECT * FROM mahasiswa WHERE id_mahasiswa=?",id)
	checkErr(err)
	n := Mahasiswa{}

	for query.Next(){
		var id_mahasiswa int
		var nama , kelas , jurusan string

		err = query.Scan(&id_mahasiswa,&nama , &kelas , &jurusan)
		checkErr(err)
		n.Id = id_mahasiswa
		n.Nama = nama
		n.Kelas = kelas
		n.Jurusan = jurusan
	}

	tmpl.ExecuteTemplate(w, "EditSiswa", n)
}

func UpdateSiswa(w http.ResponseWriter , r *http.Request) {
	koneksi := con()
	if(r.Method == "POST"){
		r.ParseForm()
		query , err := koneksi.Prepare("Update mahasiswa SET nama=? , kelas=? , jurusan=? WHERE id_mahasiswa=?")
		checkErr(err)
		query.Exec(r.PostFormValue("nama"),r.PostFormValue("kelas"),r.PostFormValue("jurusan"),r.PostFormValue("id_mahasiswa"))
		defer koneksi.Close()
		http.Redirect(w,r, "/mahasiswa", 301)
	}
}

func DeleteSiswa(w http.ResponseWriter , r *http.Request){
	koneksi := con()
	id := r.URL.Query().Get("id")

	query , err := koneksi.Prepare("DELETE FROM mahasiswa WHERE id_mahasiswa=?")
	checkErr(err)
	query.Exec(id)
	defer koneksi.Close()

	http.Redirect(w, r , "/mahasiswa", 301)
}

func checkErr(err error){
	if(err != nil){
		log.Println(err)
	}
}

