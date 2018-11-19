package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"
    _ "github.com/go-sql-driver/mysql"
    "testing"
)

type Costumer struct {
    Id  int
    AcountNumber  string
    Name  string
    Email string
}
type RDeposit struct {
    AcountNumber  string
    Name  string
    Deposit int
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "bank"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM costumer ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    cos := Costumer{}
    res := []Costumer{}
    for selDB.Next() {
        var id int
        var acountNumber, name, email string
        err = selDB.Scan(&id, &acountNumber ,&name, &email)
        if err != nil {
            panic(err.Error())
        }
        cos.Id = id
        cos.AcountNumber = acountNumber
        cos.Name = name
        cos.Email = email
        res = append(res, cos)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}


func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        acountNumber := r.FormValue("acountNumber")
        email := r.FormValue("email")
        firstDeposit := r.FormValue("firstDeposit")
        insForm, err := db.Prepare("INSERT INTO costumer(name, acountNumber, email) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())
        }else {
        res, err := insForm.Exec(name, acountNumber, email)
        id_dept, err := res.LastInsertId()
        if err != nil {
            println("Error:", err.Error())
        } else {
            println("LastInsertId:", id_dept)
            insDeposit, err := db.Prepare("INSERT INTO transaction(idCostumer, deposit) VALUES(?,?)")
            if err != nil {
                panic(err.Error())
            }
            insDeposit.Exec(id_dept, firstDeposit)
        }
        }
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func NewDeposit(w http.ResponseWriter, r *http.Request) {
  db := dbConn()
  selDB, err := db.Query("SELECT * FROM costumer ORDER BY id DESC")
  if err != nil {
      panic(err.Error())
  }
  cos := Costumer{}
  res := []Costumer{}
  for selDB.Next() {
      var id int
      var acountNumber, name, email string
      err = selDB.Scan(&id, &acountNumber ,&name, &email)
      if err != nil {
          panic(err.Error())
      }
      cos.Id = id
      cos.AcountNumber = acountNumber
      cos.Name = name
      cos.Email = email
      res = append(res, cos)
  }
  tmpl.ExecuteTemplate(w, "InsDeposit", res)
  defer db.Close()
}

func InsDeposit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        idCostumer := r.FormValue("idCostumer")
        deposit := r.FormValue("deposit")
        insDep, err := db.Prepare("INSERT INTO transaction(idCostumer, deposit) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }else {
        insDep.Exec(idCostumer, deposit)
        }
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}


func ShowTransaction(w http.ResponseWriter, r *http.Request) {
  db := dbConn()
     IdCostumer := r.URL.Query().Get("id")
     selDB, err := db.Query("SELECT costumer.acountNumber, costumer.name, transaction.deposit FROM transaction, costumer WHERE costumer.id = transaction.idCostumer AND transaction.idCostumer=?", IdCostumer)
     if err != nil {
         panic(err.Error())
     }
     rdep := RDeposit{}
     res := []RDeposit{}
     for selDB.Next() {
         var deposit int
         var name, acountNumber string
         err = selDB.Scan(&name,&acountNumber,&deposit)
         if err != nil {
             panic(err.Error())
         }
         rdep.Name = name
         rdep.AcountNumber = acountNumber
         rdep.Deposit = deposit
         res = append(res, rdep)
     }
     tmpl.ExecuteTemplate(w, "HistoryTransaction", res)
     defer db.Close()
}
func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/newdeposit", NewDeposit)
    http.HandleFunc("/insDeposit", InsDeposit)
    http.HandleFunc("/new", New)
    http.HandleFunc("/ShowTransaction", ShowTransaction)
    http.HandleFunc("/insert", Insert)
    http.ListenAndServe(":8080", nil)
}
