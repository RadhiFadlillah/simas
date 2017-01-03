package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"simas/model"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

const (
	simasDomain    string = "simas.tokoyulia.com"
	zenzivaUserKey string = "jcyhzp"
	zenzivaPassKey string = "teYa52CLLLR5FHdtQDdY"
	emailAddress   string = "m.radhi.f@gmail.com"
	emailPassword  string = "Z9t7U6NsPhQYTWGSMppbkFSsp"
	emailServer    string = "smtp.gmail.com"
	emailPort      int    = 587
	fileSuratDir   string = "."
)

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var request model.LoginRequest
	checkError(json.NewDecoder(r.Body).Decode(&request))

	// Validate input value
	if request.Email == "" {
		panic(errors.New("Email harus diisi"))
	}

	if request.Password == "" {
		panic(errors.New("Password harus diisi"))
	}

	// Get account data from database
	account := model.Account{}
	err := handler.DB.Get(&account, "SELECT * FROM account WHERE email = ?", request.Email)
	if err != nil {
		panic(errors.New("Email tidak terdaftar"))
	}

	// Compare password with database
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password))
	if err != nil {
		panic(errors.New("Email dan password tidak cocok"))
	}

	// Calculate expiration time
	nbf := time.Now()
	exp := time.Now().Add(2 * time.Hour)

	if request.Remember {
		exp = time.Date(nbf.Year(), nbf.Month(), nbf.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	}

	// Generate token
	isAdmin := false
	if account.Admin == 1 {
		isAdmin = true
	}

	isPenginput := false
	if account.Penginput == 1 {
		isPenginput = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf":   nbf.Unix(),
		"exp":   exp.Unix(),
		"sub":   account.ID,
		"admin": isAdmin,
		"input": isPenginput,
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	checkError(err)

	// Return login result
	account.Password = ""
	result := model.LoginResult{
		Account: account,
		Token:   tokenString,
	}

	delay()
	w.Header().Add("Content-Type", "application/json")
	checkError(json.NewEncoder(w).Encode(&result))
}

func (handler *Handler) SelectAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	tokenMustExist(r)

	// Parse URL query
	queries := r.URL.Query()
	page := queries.Get("page")
	keyword := queries.Get("keyword")

	// Parse page number
	pageNumber, _ := strconv.Atoi(page)
	if pageNumber == 0 {
		pageNumber = 1
	}
	offset := (pageNumber - 1) * 20

	// Create base query
	sqlCount := "SELECT FLOOR(COUNT(*) / 20) FROM account"
	sqlSelect := "SELECT id, email, nama, jabatan, telepon, admin, penginput FROM account"
	whereClause := "WHERE 1"
	args := []interface{}{}

	// Add keyword to where clause
	if keyword != "" {
		keyword += "%"
		whereClause += " AND (nama LIKE ? OR email LIKE ? OR jabatan LIKE ?)"
		args = append(args, "%"+keyword, keyword, "%"+keyword)
	}

	// Finalize query
	sqlCount += " " + whereClause
	sqlSelect += " " + whereClause

	// Get max page from database
	var maxPage int
	err := handler.DB.Get(&maxPage, sqlCount, args...)
	checkError(err)

	// Add order and limit clause to select query
	sqlSelect += " ORDER BY nama LIMIT 20 OFFSET ?"
	args = append(args, offset)

	// Select all account from database
	listAccount := []model.Account{}
	err = handler.DB.Select(&listAccount, sqlSelect, args...)
	checkError(err)

	// Encode result
	pageListAccount := model.PageListAccount{
		Page:    pageNumber,
		MaxPage: maxPage,
		Item:    listAccount,
	}

	delay()
	w.Header().Add("Content-Type", "application/json")
	checkError(json.NewEncoder(w).Encode(&pageListAccount))
}

func (handler *Handler) InsertAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	tokenMustExist(r)

	// Decode request
	var account model.Account
	checkError(json.NewDecoder(r.Body).Decode(&account))

	// Validate input
	if account.Nama == "" {
		panic(errors.New("Nama harus diisi"))
	}

	if account.Email == "" {
		panic(errors.New("Email harus diisi"))
	}

	if account.Jabatan == "" {
		panic(errors.New("Jabatan harus diisi"))
	}

	// Generate password
	randomPassword := randomString(10)

	// Hash password with bcrypt
	password := []byte(randomPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	checkError(err)

	// Insert account to database
	res := handler.DB.MustExec(`INSERT INTO account 
		(nama, email, jabatan, telepon, password, admin, penginput) VALUES 
		(?, ?, ?, ?, ?, ?, ?)`,
		account.Nama,
		account.Email,
		account.Jabatan,
		account.Telepon,
		hashedPassword,
		account.Admin,
		account.Penginput)

	// Prepare SMS
	smsMessage := fmt.Sprintf(
		`Anda telah didaftarkan ke Sistem Manajemen Surat Fakultas Teknik UPR. 
		Silakan login ke %s dengan menggunakan email %s dan password %s`,
		simasDomain, account.Email, password)

	// Prepare email
	buffer := bytes.Buffer{}
	newAccountTemplate.Execute(&buffer, model.EmailNewAccount{
		Account:  account,
		Domain:   simasDomain,
		Password: randomPassword})

	// Send SMS and Email
	go handler.sendSMS(account.Telepon, smsMessage)
	go handler.sendEmail(account.Email, "Selamat Datang di SIMAS FT UPR", buffer.String())

	// Return inserted ID
	delay()
	id, _ := res.LastInsertId()
	fmt.Fprint(w, id)
}

func (handler *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	tokenMustExist(r)

	// Decode request
	var account model.Account
	checkError(json.NewDecoder(r.Body).Decode(&account))

	// Validate input
	if account.Nama == "" {
		panic(errors.New("Nama harus diisi"))
	}

	if account.Email == "" {
		panic(errors.New("Email harus diisi"))
	}

	if account.Jabatan == "" {
		panic(errors.New("Jabatan harus diisi"))
	}

	// Get requested account data
	var oldData model.Account
	checkError(handler.DB.Get(&oldData, "SELECT * FROM account WHERE id = ?", account.ID))

	// If it is admin, check number of existing admin in database
	if oldData.Admin == 1 {
		var nAdmin int
		checkError(handler.DB.Get(&nAdmin, "SELECT COUNT(*) FROM account WHERE admin = 1"))

		if nAdmin == 1 && account.Admin == 0 {
			panic(errors.New("Setidaknya harus ada satu admin"))
		}
	}

	// If it is penginput, check number of penginput admin in database
	if oldData.Penginput == 1 {
		var nPenginput int
		checkError(handler.DB.Get(&nPenginput, "SELECT COUNT(*) FROM account WHERE penginput = 1"))

		if nPenginput == 1 && account.Penginput == 0 {
			panic(errors.New("Setidaknya harus ada satu penginput"))
		}
	}

	// Update account in database
	handler.DB.MustExec(`UPDATE account SET nama = ?, email = ?, jabatan = ?, 
		telepon = ?, admin = ?, penginput = ? WHERE id = ?`,
		account.Nama,
		account.Email,
		account.Jabatan,
		account.Telepon,
		account.Admin,
		account.Penginput,
		account.ID)

	// Return updated ID
	delay()
	fmt.Fprint(w, account.ID)
}

func (handler *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	tokenMustExist(r)

	// Get id account from URL address
	idAccount := ps.ByName("id")

	// Get requested account data
	var account model.Account
	checkError(handler.DB.Get(&account, "SELECT * FROM account WHERE id = ?", idAccount))

	// If it is admin, check number of existing admin in database
	if account.Admin == 1 {
		var nAdmin int
		checkError(handler.DB.Get(&nAdmin, "SELECT COUNT(*) FROM account WHERE admin = 1"))

		if nAdmin == 1 {
			panic(errors.New("Setidaknya harus ada satu admin"))
		}
	}

	// If it is penginput, check number of penginput admin in database
	if account.Penginput == 1 {
		var nPenginput int
		checkError(handler.DB.Get(&nPenginput, "SELECT COUNT(*) FROM account WHERE penginput = 1"))

		if nPenginput == 1 {
			panic(errors.New("Setidaknya harus ada satu penginput"))
		}
	}

	// Delete account in database
	handler.DB.MustExec("DELETE FROM account WHERE id = ?", idAccount)

	// Return ID
	delay()
	fmt.Fprint(w, idAccount)
}

func (handler *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	claims := tokenMustExist(r)

	// Decode request
	var request model.UpdatePasswordRequest
	checkError(json.NewDecoder(r.Body).Decode(&request))

	// Validate input
	if request.Password == "" {
		panic(errors.New("Password baru harus diisi"))
	}

	// Get account data from database
	id := claims["sub"]
	account := model.Account{}
	err := handler.DB.Get(&account, "SELECT * FROM account WHERE id = ?", id)
	if err != nil {
		panic(errors.New("User tidak terdaftar"))
	}

	// Compare password with database
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.PasswordLama))
	if err != nil {
		panic(errors.New("Password lama salah"))
	}

	// Hash password with bcrypt
	password := []byte(request.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 10)
	checkError(err)

	// Update password in database
	handler.DB.MustExec("UPDATE account SET password = ? WHERE id = ?", hashedPassword, id)

	// Return account ID
	delay()
	fmt.Fprint(w, id)
}

func (handler *Handler) SelectSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	claims := tokenMustExist(r)

	// Parse URL query
	queries := r.URL.Query()
	page := queries.Get("page")
	keyword := queries.Get("keyword")
	_, diproses := queries["diproses"]
	_, diarsip := queries["diarsip"]
	_, ditindak := queries["ditindak"]
	_, kelola := queries["kelola"]

	// Parse page number
	pageNumber, _ := strconv.Atoi(page)
	if pageNumber == 0 {
		pageNumber = 1
	}
	offset := (pageNumber - 1) * 20

	// Create base query
	sqlCount := "SELECT FLOOR(COUNT(*) / 20) FROM disposisi d LEFT JOIN surat s ON d.surat_id = s.id"
	sqlSelect := "SELECT s.*, d.status, d.read FROM disposisi d LEFT JOIN surat s ON d.surat_id = s.id"
	whereClause := "WHERE d.tujuan_id = ?"
	args := []interface{}{claims["sub"]}

	if penginput, ok := claims["input"].(bool); kelola && ok && penginput {
		sqlCount = "SELECT FLOOR(COUNT(*) / 20) + 1 FROM surat"
		sqlSelect = "SELECT *, 1 `read` FROM surat"
		whereClause = "WHERE 1"
		args = []interface{}{}
	}

	// Add keyword to where clause
	if keyword != "" {
		keyword += "%"
		whereClause += " AND (nomor LIKE ? OR perihal LIKE ? OR sumber LIKE ?)"
		args = append(args, keyword, "%"+keyword, "%"+keyword)
	}

	// Add diproses to where clause
	if diproses {
		whereClause += " AND status = 0"
	}

	// Add diarsip to where clause
	if diarsip {
		whereClause += " AND (status = 1 OR status = 3)"
	}

	// Add ditindak to where clause
	if ditindak {
		whereClause += " AND (status = 2 OR status = 3)"
	}

	// Finalize query
	sqlCount += " " + whereClause
	sqlSelect += " " + whereClause

	// Get max page from database
	var maxPage int
	err := handler.DB.Get(&maxPage, sqlCount, args...)
	checkError(err)

	// Add order and limit clause to select query
	sqlSelect += " ORDER BY prioritas DESC, tanggal DESC LIMIT 20 OFFSET ?"
	args = append(args, offset)

	// Select all surat from database
	listSurat := []model.Surat{}
	err = handler.DB.Select(&listSurat, sqlSelect, args...)
	checkError(err)

	// Get count of new surat
	newCount := []model.NewSuratCount{}
	err = handler.DB.Select(&newCount,
		`SELECT status, COUNT(*) count 
		FROM disposisi WHERE disposisi.read = 0 AND tujuan_id = ? GROUP BY status`,
		claims["sub"])
	checkError(err)

	// Encode result
	pageListSurat := model.PageListSurat{
		Page:     pageNumber,
		MaxPage:  maxPage,
		NewCount: newCount,
		Item:     listSurat,
	}

	delay()
	w.Header().Add("Content-Type", "application/json")
	checkError(json.NewEncoder(w).Encode(&pageListSurat))
}

func (handler *Handler) GetSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check token
	claims := tokenMustExist(r)

	// Get id surat from URL address
	idSurat := ps.ByName("id")
	intIdSurat, _ := strconv.Atoi(idSurat)

	// Parse URL query
	queries := r.URL.Query()
	_, kelola := queries["kelola"]

	// Retrieve data surat from database
	surat := model.Surat{}
	err := handler.DB.Get(&surat, `SELECT s.id, s.nomor, s.perihal, s.sumber, 
		a.nama tujuan, a.jabatan, s.tanggal, s.waktu_terima, s.prioritas, d.status
		FROM surat s LEFT JOIN disposisi d ON s.id = d.surat_id
		LEFT JOIN account a ON d.tujuan_id = a.id 
		WHERE s.id = ? AND d.parent_id IS NULL`, idSurat)
	checkError(err)

	// Retrieve data disposisi and timeline from database
	disposisi := model.Disposisi{}
	timeline := []model.TimelineItem{}
	if !kelola {
		err = handler.DB.Get(&disposisi, `SELECT d.id, d.parent_id, 
			IFNULL(a.id, 0) sumber_id, IFNULL(a.nama, '') sumber, 
			IFNULL(a.jabatan, '') jabatan, d.waktu, d.status,
			IFNULL(de.deskripsi, '') deskripsi
			FROM disposisi d
			LEFT JOIN disposisi_detail de ON d.id = de.id
			LEFT JOIN account a ON a.id = (
				SELECT tujuan_id FROM disposisi 
				WHERE surat_id = d.surat_id AND id = d.parent_id
			) WHERE d.surat_id = ? AND d.tujuan_id = ?`,
			idSurat, claims["sub"])
		checkError(err)

		if err != sql.ErrNoRows {
			timeline = handler.getSuratTimeline(intIdSurat, disposisi.ID)
		}
	}

	// Mark as read
	handler.DB.MustExec("UPDATE disposisi SET `read` = 1 WHERE id = ?", disposisi.ID)

	// Retrieve image files
	filePattern := fmt.Sprintf("%s/%d-*", fileSuratDir, intIdSurat)

	imageFiles, _ := filepath.Glob(filePattern)
	if len(imageFiles) == 0 {
		imageFiles = make([]string, 0)
	}

	// Encode result
	pageSurat := model.PageSurat{
		Surat:     surat,
		Disposisi: disposisi,
		Files:     imageFiles,
		Timeline:  timeline,
	}

	delay()
	w.Header().Add("Content-Type", "application/json")
	checkError(json.NewEncoder(w).Encode(&pageSurat))
}

func (handler *Handler) InsertSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode form request
	reader, err := r.MultipartReader()
	checkError(err)

	form, err := reader.ReadForm(1000000)
	checkError(err)

	// Validate input value
	tujuanID := 0
	if len(form.Value["tujuanId"]) != 0 {
		tujuanID, _ = strconv.Atoi(form.Value["tujuanId"][0])
	}

	prioritas := 0
	if len(form.Value["prioritas"]) != 0 {
		prioritas, _ = strconv.Atoi(form.Value["prioritas"][0])
	}

	if len(form.Value["nomor"]) == 0 || form.Value["nomor"][0] == "" {
		panic(errors.New("Nomor surat harus diisi"))
	}

	if len(form.Value["perihal"]) == 0 || form.Value["perihal"][0] == "" {
		panic(errors.New("Perihal surat harus diisi"))
	}

	if len(form.Value["perihal"]) == 0 || form.Value["sumber"][0] == "" {
		panic(errors.New("Sumber surat harus diisi"))
	}

	if tujuanID == 0 {
		panic(errors.New("Tujuan surat harus diisi"))
	}

	if len(form.Value["tanggal"]) == 0 || form.Value["tanggal"][0] == "" {
		panic(errors.New("Tanggal surat harus diisi"))
	}

	if len(form.Value["waktuTerima"]) == 0 || form.Value["waktuTerima"][0] == "" {
		panic(errors.New("Waktu terima surat harus diisi"))
	}

	// Validate files
	if len(form.File["files"]) > 5 {
		panic(errors.New("Jumlah file maksimal 5 gambar JPG atau PNG"))
	}

	for _, uploadedFile := range form.File["files"] {
		contentType := uploadedFile.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" {
			panic(errors.New("File harus berupa gambar JPG atau PNG"))
		}
	}

	// Convert form to struct
	surat := model.Surat{
		Nomor:       form.Value["nomor"][0],
		Perihal:     form.Value["perihal"][0],
		Sumber:      form.Value["sumber"][0],
		Tanggal:     form.Value["tanggal"][0],
		WaktuTerima: form.Value["waktuTerima"][0],
		Prioritas:   prioritas,
	}

	// Begin transaction
	tx := handler.DB.MustBegin()

	// Insert surat to database
	res := tx.MustExec(`INSERT INTO surat
		(nomor, perihal, sumber, tanggal, waktu_terima, prioritas)
		VALUES (?, ?, ?, ?, ?, ?)`,
		surat.Nomor,
		surat.Perihal,
		surat.Sumber,
		surat.Tanggal,
		surat.WaktuTerima,
		surat.Prioritas)

	// Get inserted ID
	id, err := res.LastInsertId()
	checkError(err)

	surat.ID = int(id)

	// Insert disposisi to database
	tx.MustExec(`INSERT INTO disposisi
		(surat_id, tujuan_id, waktu) VALUES (?, ?, ?)`,
		surat.ID, tujuanID, surat.WaktuTerima)

	// Save image
	jpgFile := []string{}
	pngFile := []string{}
	for i, uploadedFile := range form.File["files"] {
		fileName := fmt.Sprintf("%s/%d-%d", fileSuratDir, surat.ID, i+1)

		contentType := uploadedFile.Header.Get("Content-Type")
		if contentType == "image/jpeg" {
			fileName += ".jpg"
			jpgFile = append(jpgFile, fileName)
		} else if contentType == "image/png" {
			fileName += ".png"
			pngFile = append(pngFile, fileName)
		}

		src, err := uploadedFile.Open()
		defer src.Close()
		checkError(err)

		dst, err := os.Create(fileName)
		defer dst.Close()
		checkError(err)

		_, err = io.Copy(dst, src)
		checkError(err)
	}

	// Compress image
	pngquantArgs := append([]string{"--speed", "11", "--force", "-ext", ".png"}, pngFile...)
	jpegOptimArgs := append([]string{"-q", "-s", "--all-progressive", "-f", "-m", "80"}, jpgFile...)

	cmdPngquant := exec.Command("pngquant", pngquantArgs...)
	cmdJpegOptim := exec.Command("jpegoptim", jpegOptimArgs...)

	go cmdPngquant.Run()
	go cmdJpegOptim.Run()

	// Commit database
	err = tx.Commit()
	checkError(err)

	// Prepare SMS
	smsMessage := fmt.Sprintf(
		`Anda mendapat surat baru dari %s dengan nomor %s, perihal "%s". SIMAS FT UPR`,
		surat.Sumber, surat.Nomor, surat.Perihal)

	// Prepare email
	buffer := bytes.Buffer{}
	newSuratTemplate.Execute(&buffer, model.EmailNewSurat{
		Surat:  surat,
		Domain: simasDomain})

	emailSubject := fmt.Sprintf("Surat Baru No. %s - SIMAS FT UPR", surat.Nomor)

	// Get accound tujuan surat
	account := model.Account{}
	handler.DB.Get(&account, "SELECT * FROM account WHERE id = ?", tujuanID)

	// Send SMS and Email
	go handler.sendSMS(account.Telepon, smsMessage)
	go handler.sendEmail(account.Email, emailSubject, buffer.String())

	delay()
	fmt.Fprint(w, surat.ID)
}

func (handler *Handler) UpdateSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	reader, err := r.MultipartReader()
	checkError(err)

	form, err := reader.ReadForm(1000000)
	checkError(err)

	// Validate input value
	suratID, _ := strconv.Atoi(form.Value["id"][0])

	prioritas := 0
	if len(form.Value["prioritas"]) != 0 {
		prioritas, _ = strconv.Atoi(form.Value["prioritas"][0])
	}

	if len(form.Value["nomor"]) == 0 || form.Value["nomor"][0] == "" {
		panic(errors.New("Nomor surat harus diisi"))
	}

	if len(form.Value["perihal"]) == 0 || form.Value["perihal"][0] == "" {
		panic(errors.New("Perihal surat harus diisi"))
	}

	if len(form.Value["perihal"]) == 0 || form.Value["sumber"][0] == "" {
		panic(errors.New("Sumber surat harus diisi"))
	}

	if len(form.Value["tanggal"]) == 0 || form.Value["tanggal"][0] == "" {
		panic(errors.New("Tanggal surat harus diisi"))
	}

	if len(form.Value["waktuTerima"]) == 0 || form.Value["waktuTerima"][0] == "" {
		panic(errors.New("Waktu terima surat harus diisi"))
	}

	// Retrieve existing files
	filePattern := fmt.Sprintf("%s/%d-*", fileSuratDir, suratID)
	existingFiles, _ := filepath.Glob(filePattern)

	// Fetch file to be deleted
	deletedFiles := form.Value["deleted"]

	// Check total of file
	newFiles := form.File["files"]
	if len(existingFiles)-len(deletedFiles)+len(newFiles) > 5 {
		panic(errors.New("Jumlah file maksimal 5 gambar JPG atau PNG"))
	}

	// Make sure file is JPG or PNG
	for _, uploadedFile := range newFiles {
		contentType := uploadedFile.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" {
			panic(errors.New("File harus berupa gambar JPG atau PNG"))
		}
	}

	// Update surat in database
	handler.DB.MustExec(`UPDATE surat SET nomor = ?, perihal = ?,
		sumber = ?, tanggal = ?, waktu_terima = ?, prioritas = ? WHERE id = ?`,
		form.Value["nomor"][0],
		form.Value["perihal"][0],
		form.Value["sumber"][0],
		form.Value["tanggal"][0],
		form.Value["waktuTerima"][0],
		prioritas,
		suratID)

	// Delete file
	for _, deletedFile := range deletedFiles {
		path := filepath.Join(fileSuratDir, deletedFile)
		os.Remove(path)
	}

	// Get max number
	maxNumber := 0
	for _, existedFile := range existingFiles {
		name := existedFile[:len(existedFile)-4]
		name = strings.Split(name, "-")[1]
		number, _ := strconv.Atoi(name)
		if number > maxNumber {
			maxNumber = number
		}
	}

	// Save new image
	jpgFile := []string{}
	pngFile := []string{}
	for i, uploadedFile := range newFiles {
		fileName := fmt.Sprintf("%s/%d-%d", fileSuratDir, suratID, maxNumber+i+1)

		contentType := uploadedFile.Header.Get("Content-Type")
		if contentType == "image/jpeg" {
			fileName += ".jpg"
			jpgFile = append(jpgFile, fileName)
		} else if contentType == "image/png" {
			fileName += ".png"
			pngFile = append(pngFile, fileName)
		}

		src, err := uploadedFile.Open()
		defer src.Close()
		checkError(err)

		dst, err := os.Create(fileName)
		defer dst.Close()
		checkError(err)

		_, err = io.Copy(dst, src)
		checkError(err)
	}

	// Compress new image
	pngquantArgs := append([]string{"--speed", "11", "--force", "-ext", ".png"}, pngFile...)
	jpegOptimArgs := append([]string{"-q", "-s", "--all-progressive", "-f", "-m", "80"}, jpgFile...)

	cmdPngquant := exec.Command("pngquant", pngquantArgs...)
	cmdJpegOptim := exec.Command("jpegoptim", jpegOptimArgs...)

	go cmdPngquant.Run()
	go cmdJpegOptim.Run()

	delay()
	fmt.Fprint(w, suratID)
}

func (handler *Handler) DeleteSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	suratID := ps.ByName("id")

	// Delete surat in database
	handler.DB.MustExec("DELETE FROM surat WHERE id = ?", suratID)

	// Delete existing files
	filePattern := fmt.Sprintf("%s/%s-*", fileSuratDir, suratID)
	existingFiles, _ := filepath.Glob(filePattern)

	for _, fileName := range existingFiles {
		path := filepath.Join(fileSuratDir, fileName)
		os.Remove(path)
	}

	// Return id surat
	delay()
	fmt.Fprint(w, suratID)
}

func (handler *Handler) GetFileSurat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileName := ps.ByName("name")
	path := filepath.Join(fileSuratDir, fileName)

	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		delay()
		http.ServeFile(w, r, path)
		return
	}

	http.NotFound(w, r)
}

func (handler *Handler) InsertDisposisi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var disposisi model.Disposisi
	checkError(json.NewDecoder(r.Body).Decode(&disposisi))

	// Convert parent id to null if needed
	parentID := interface{}(nil)
	if disposisi.ParentID.Int64 != 0 {
		parentID = interface{}(disposisi.ParentID.Int64)
	}

	// Begin transaction
	tx := handler.DB.MustBegin()

	// Insert disposisi to database
	res := tx.MustExec(`INSERT INTO disposisi 
		(surat_id, tujuan_id, parent_id) VALUES (?, ?, ?)`,
		disposisi.SuratID, disposisi.TujuanID, parentID)

	// Get inserted id
	id, _ := res.LastInsertId()

	// Insert detail disposisi
	tx.MustExec(`INSERT INTO disposisi_detail (id, deskripsi) VALUES (?, ?)`,
		id, disposisi.Deskripsi)

	// Commit
	err := tx.Commit()
	checkError(err)

	// Get sumber disposisi
	sumber := model.Account{}
	handler.DB.Get(&sumber, `SELECT * FROM account 
		WHERE id = (SELECT tujuan_id FROM disposisi WHERE id = ?)`, parentID)

	// Get surat data
	surat := model.Surat{}
	handler.DB.Get(&surat, "SELECT * FROM surat WHERE id = ?", disposisi.SuratID)

	// Prepare SMS
	smsMessage := fmt.Sprintf(
		`Anda mendapat disposisi dari %s tentang surat dengan nomor %s, perihal "%s". SIMAS FT UPR`,
		sumber.Jabatan, surat.Nomor, surat.Perihal)

	// Prepare email
	buffer := bytes.Buffer{}
	newDisposisiTemplate.Execute(&buffer, model.EmailNewDisposisi{
		Sumber:  sumber.Nama,
		Jabatan: sumber.Jabatan,
		Surat:   surat,
		Domain:  simasDomain})

	emailSubject := fmt.Sprintf("Disposisi Surat %s Dari %s - SIMAS FT UPR", surat.Nomor, sumber.Jabatan)

	// Get account tujuan surat
	tujuan := model.Account{}
	handler.DB.Get(&tujuan, "SELECT * FROM account WHERE id = ?", disposisi.TujuanID)

	// Send SMS and Email
	go handler.sendSMS(tujuan.Telepon, smsMessage)
	go handler.sendEmail(tujuan.Email, emailSubject, buffer.String())

	// Return ID
	fmt.Fprint(w, id)
}

func (handler *Handler) InsertDiarsip(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var disposisi model.Disposisi
	checkError(json.NewDecoder(r.Body).Decode(&disposisi))

	// Begin transaction
	tx := handler.DB.MustBegin()

	// Insert diarsip to database
	tx.MustExec("UPDATE disposisi SET status = 1 WHERE id = ?", disposisi.ID)

	// Update parent status
	handler.updateParentStatus(tx, disposisi.SuratID, 0)

	// Commit
	err := tx.Commit()
	checkError(err)

	// Return ID
	fmt.Fprint(w, disposisi.ID)
}

func (handler *Handler) InsertDitindak(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Decode request
	var disposisi model.Disposisi
	checkError(json.NewDecoder(r.Body).Decode(&disposisi))

	// Begin transaction
	tx := handler.DB.MustBegin()

	// Insert ditindak to database
	tx.MustExec("UPDATE disposisi SET status = 2 WHERE id = ?", disposisi.ID)

	// Update parent status
	handler.updateParentStatus(tx, disposisi.SuratID, 0)

	// Commit
	err := tx.Commit()
	checkError(err)

	// Return ID
	fmt.Fprint(w, disposisi.ID)
}

func (handler *Handler) getSuratTimeline(suratID int, parentID int) []model.TimelineItem {
	// Init variable
	items := []model.TimelineItem{}

	// Get timeline from database
	err := handler.DB.Select(&items,
		`SELECT di.id, ac.nama tujuan, ac.jabatan, 
		CONVERT_TZ(di.modified, @@session.time_zone, '+00:00') waktu, di.status
		FROM disposisi di LEFT JOIN account ac ON di.tujuan_id = ac.id
		WHERE di.surat_id = ? AND IFNULL(di.parent_id, 0) = ? ORDER BY di.modified`,
		suratID, parentID)
	checkError(err)

	// Get children for each timeline item
	for i, item := range items {
		items[i].Children = handler.getSuratTimeline(suratID, item.ID)
	}

	return items
}

func (handler *Handler) updateParentStatus(tx *sqlx.Tx, suratID int, parentID int) ([]model.DisposisiUpdate, int) {
	// Init variable
	items := []model.DisposisiUpdate{}
	listStatus := make(map[int]int)

	// Get timeline from database
	err := tx.Select(&items,
		`SELECT d.id, d.status, s.nomor, s.perihal, a.email, a.telepon 
		FROM disposisi d 
		LEFT JOIN surat s ON d.surat_id = s.id
		LEFT JOIN account a ON d.tujuan_id = a.id
		WHERE surat_id = ? AND IFNULL(parent_id, 0) = ?`,
		suratID, parentID)
	checkError(err)

	// Get children for each timeline item
	for i, item := range items {
		var childStatus = 0
		items[i].Children, childStatus = handler.updateParentStatus(tx, suratID, item.ID)

		if len(items[i].Children) == 0 {
			listStatus[item.Status]++
		} else {
			listStatus[childStatus]++
			if item.Status != childStatus {
				tx.MustExec("UPDATE disposisi SET status = ?, `read` = 0 WHERE id = ?",
					childStatus, item.ID)

				if item.Status == 0 {
					// Set status value
					statusValue := "diarsipkan dan ditindaklanjuti"
					if childStatus == 1 {
						statusValue = "diarsipkan"
					} else if childStatus == 2 {
						statusValue = "ditindaklanjuti"
					}

					// Prepare SMS
					smsMessage := fmt.Sprintf(
						`Surat dengan nomor %s perihal "%s" telah %s. SIMAS FT UPR`,
						item.Nomor, item.Perihal, statusValue)

					// Prepare email
					buffer := bytes.Buffer{}
					newStatusTemplate.Execute(&buffer, model.EmailNewStatus{
						Nomor:   item.Nomor,
						Perihal: item.Perihal,
						Status:  childStatus,
						Domain:  simasDomain})

					emailSubject := fmt.Sprintf("Surat No. %s %s - SIMAS FT UPR", item.Nomor, statusValue)

					// Send SMS and Email
					go handler.sendSMS(item.Telepon, smsMessage)
					go handler.sendEmail(item.Email, emailSubject, buffer.String())
				}
			}
		}
	}

	// Calculate final status
	finalStatus := 0
	if listStatus[0] == 0 {
		if listStatus[3] != 0 {
			finalStatus = 3
		} else if listStatus[1] != 0 && listStatus[2] != 0 {
			finalStatus = 3
		} else if listStatus[2] != 0 {
			finalStatus = 2
		} else if listStatus[1] != 0 {
			finalStatus = 1
		}
	}

	return items, finalStatus
}

func (handler *Handler) sendSMS(number string, message string) {
	if number == "" {
		return
	}

	zenzivaQuery := url.Values{}
	zenzivaQuery.Set("userkey", zenzivaUserKey)
	zenzivaQuery.Set("passkey", zenzivaPassKey)
	zenzivaQuery.Set("nohp", number)
	zenzivaQuery.Set("pesan", message)

	zenzivaURL, _ := url.Parse("https://reguler.zenziva.net/apps/smsapi.php")

	zenzivaURL.RawQuery = zenzivaQuery.Encode()

	http.Get(zenzivaURL.String())
}

func (handler *Handler) sendEmail(target string, subject string, body string) {
	header := "MIME-Version: 1.0" + "\r\n" +
		"Content-type: text/html" + "\r\n" +
		"Subject: " + subject + "\r\n\r\n"

	auth := smtp.PlainAuth("", emailAddress, emailPassword, emailServer)
	server := fmt.Sprintf("%s:%d", emailServer, emailPort)
	smtp.SendMail(server, auth, emailAddress, []string{target}, []byte(header+body))
}
