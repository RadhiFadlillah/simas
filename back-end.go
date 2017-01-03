package main

import (
	"fmt"
	"log"
	"net/http"
	"simas/handler"
	"simas/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type BackEnd struct {
	DB         *sqlx.DB
	Config     model.Configuration
	PortNumber int
}

func NewBackEnd(port int, config model.Configuration) BackEnd {
	dbSource := fmt.Sprintf("%s:%s@/%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName)

	backEnd := BackEnd{
		DB:         sqlx.MustConnect("mysql", dbSource),
		Config:     config,
		PortNumber: port,
	}

	backEnd.generateTable()

	return backEnd
}

func (backend *BackEnd) ServeApp() {
	// Create handler
	hdl := handler.Handler{
		DB:     backend.DB,
		Config: backend.Config,
	}

	// Create router
	router := httprouter.New()

	// Handle path to UI
	router.GET("/res/*filepath", hdl.ServeFile)
	router.GET("/style/*filepath", hdl.ServeFile)
	router.GET("/", hdl.ServeIndexPage)
	router.GET("/login", hdl.ServeLoginPage)

	// Handle path to API
	router.POST("/api/login", hdl.Login)

	router.GET("/api/account", hdl.SelectAccount)
	router.PUT("/api/account", hdl.UpdateAccount)
	router.POST("/api/account", hdl.InsertAccount)
	router.POST("/api/account/password", hdl.UpdatePassword)
	router.DELETE("/api/account/:id", hdl.DeleteAccount)

	router.GET("/api/surat", hdl.SelectSurat)
	router.GET("/api/surat/id/:id", hdl.GetSurat)
	router.GET("/api/surat/image/:name", hdl.GetFileSurat)
	router.PUT("/api/surat", hdl.UpdateSurat)
	router.POST("/api/surat", hdl.InsertSurat)
	router.DELETE("/api/surat/:id", hdl.DeleteSurat)

	router.POST("/api/disposisi", hdl.InsertDisposisi)
	router.POST("/api/diarsip", hdl.InsertDiarsip)
	router.POST("/api/ditindak", hdl.InsertDitindak)

	// Set panic handler
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, arg interface{}) {
		http.Error(w, fmt.Sprint(arg), 500)
	}

	// Serve app
	log.Printf("Serve app in port %d\n", backend.PortNumber)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", backend.PortNumber), router))
}

func (backend *BackEnd) Close() {
	backend.DB.Close()
}

func (backend *BackEnd) generateTable() {
	// Begin transaction
	tx := backend.DB.MustBegin()

	// Create table
	tx.MustExec(`CREATE TABLE IF NOT EXISTS account (
		id int(11) NOT NULL AUTO_INCREMENT,
		email varchar(100) NOT NULL,
		nama varchar(100) NOT NULL,
		jabatan varchar(100) NOT NULL,
		telepon varchar(20) NOT NULL,
		password binary(80) NOT NULL,
		admin tinyint(4) NOT NULL DEFAULT '0',
		penginput tinyint(4) NOT NULL DEFAULT '0',
		PRIMARY KEY (id),
		UNIQUE KEY account_email_UN (email))`)

	tx.MustExec(`CREATE TABLE IF NOT EXISTS surat (
		id int(11) NOT NULL AUTO_INCREMENT,
		nomor varchar(50) NOT NULL,
		perihal varchar(300) NOT NULL,
		sumber varchar(100) NOT NULL,
		tanggal date NOT NULL,
		waktu_terima datetime NOT NULL,
		prioritas tinyint(4) NOT NULL DEFAULT '0',
		PRIMARY KEY (id))`)

	tx.MustExec(`CREATE TABLE IF NOT EXISTS disposisi (
		id int(11) NOT NULL AUTO_INCREMENT,
		surat_id int(11) NOT NULL,
		parent_id int(11) DEFAULT NULL,
		tujuan_id int(11) NOT NULL,
		waktu timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		status tinyint(4) NOT NULL DEFAULT '0',
		modified timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,` +
		"`read` tinyint(4) NOT NULL DEFAULT '0'," +
		`PRIMARY KEY (id),
		KEY disposisi_surat_FK (surat_id),
		KEY disposisi_account_target_FK (tujuan_id),
		KEY disposisi_parent_FK (parent_id),
		CONSTRAINT disposisi_account_target_FK FOREIGN KEY (tujuan_id) REFERENCES account (id) ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT disposisi_parent_FK FOREIGN KEY (parent_id) REFERENCES disposisi (id) ON DELETE CASCADE ON UPDATE CASCADE,
		CONSTRAINT disposisi_surat_FK FOREIGN KEY (surat_id) REFERENCES surat (id) ON DELETE CASCADE ON UPDATE CASCADE)`)

	tx.MustExec(`CREATE TABLE IF NOT EXISTS disposisi_detail (
		id int(11) NOT NULL,
		deskripsi text NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT disposisi_detail_FK FOREIGN KEY (id) REFERENCES disposisi (id) ON DELETE CASCADE ON UPDATE CASCADE)`)

	// If there are no account, create new account
	// with username admin and password admin
	var nAccount int
	err := tx.Get(&nAccount, "SELECT COUNT(*) FROM account")
	checkError(err)

	if nAccount == 0 {
		password := []byte("admin")
		hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
		checkError(err)

		tx.MustExec(`INSERT INTO account 
			(email, nama, password, jabatan, admin, penginput) VALUES (?, ?, ?, ?, ?, ?)`,
			"admin@simas", "Administrator", hashedPassword, "Administrator", 1, 1)
	}

	// Commit transaction
	err = tx.Commit()
	checkError(err)
}
