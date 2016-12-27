package model

import "gopkg.in/guregu/null.v3"

type Account struct {
	ID        int    `db:"id"        json:"id"`
	Email     string `db:"email"     json:"email"`
	Nama      string `db:"nama"      json:"nama"`
	Jabatan   string `db:"jabatan"   json:"jabatan"`
	Telepon   string `db:"telepon"   json:"telepon"`
	Password  string `db:"password"  json:"password"`
	Admin     int    `db:"admin"     json:"admin"`
	Penginput int    `db:"penginput" json:"penginput"`
}

type Surat struct {
	ID          int    `db:"id"           json:"id"`
	Nomor       string `db:"nomor"        json:"nomor"`
	Perihal     string `db:"perihal"      json:"perihal"`
	Sumber      string `db:"sumber"       json:"sumber"`
	Tanggal     string `db:"tanggal"      json:"tanggal"`
	WaktuTerima string `db:"waktu_terima" json:"waktuTerima"`
	Prioritas   int    `db:"prioritas"    json:"prioritas"`
}

type Disposisi struct {
	ID        int      `db:"id"        json:"id"`
	SuratID   int      `db:"surat_id"  json:"suratId"`
	ParentID  null.Int `db:"parent_id" json:"parentId"`
	SumberID  null.Int `db:"sumber_id" json:"sumberId"`
	TargetID  int      `db:"target_id" json:"targetId"`
	Waktu     string   `db:"waktu"     json:"waktu"`
	Status    int      `db:"status"    json:"status"`
	Modified  string   `db:"modified"  json:"modified"`
	Deskripsi string   `db:"deskripsi" json:"deskripsi"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginResult struct {
	Account Account `json:"account"`
	Token   string  `json:"token"`
}

type UpdatePasswordRequest struct {
	PasswordLama string `json:"passwordLama"`
	Password     string `json:"password"`
}

type PageListSurat struct {
	Page    int     `json:"page"`
	MaxPage int     `json:"maxPage"`
	Item    []Surat `json:"item"`
}

type PageListAccount struct {
	Page    int       `json:"page"`
	MaxPage int       `json:"maxPage"`
	Item    []Account `json:"item"`
}

type PageSurat struct {
	Surat    Surat          `json:"surat"`
	Timeline []TimelineItem `json:"timeline"`
}

type TimelineItem struct {
	ID       int            `db:"id"      json:"id"`
	Target   string         `db:"target"  json:"target"`
	Jabatan  string         `db:"jabatan" json:"jabatan"`
	Waktu    string         `db:"waktu"   json:"waktu"`
	Status   int            `db:"status"  json:"status"`
	Children []TimelineItem `json:"children"`
}

type TimelineStatus struct {
	ID       int              `db:"id"     json:"id"`
	Status   int              `db:"status" json:"status"`
	Children []TimelineStatus `json:"children"`
}
