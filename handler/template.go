package handler

import (
	"html/template"
)

var (
	newAccountTemplate   *template.Template
	newSuratTemplate     *template.Template
	newDisposisiTemplate *template.Template
	newStatusTemplate    *template.Template
)

func init() {
	newAccountTemplate, _ = template.New("emailTemplate").Delims("[[", "]]").Parse(emailNewAccount)
	newSuratTemplate, _ = template.New("emailTemplate").Delims("[[", "]]").Parse(emailNewSurat)
	newDisposisiTemplate, _ = template.New("emailTemplate").Delims("[[", "]]").Parse(emailNewDisposisi)
	newStatusTemplate, _ = template.New("emailTemplate").Delims("[[", "]]").Parse(emailNewStatus)
}

const emailNewAccount string = `
<!DOCTYPE html>
<html>

<head>
    <title>SIMAS - Sistem Manajemen Surat</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:400,600">
    <style type="text/css">
        @media screen and (max-width: 500px) {
            .main .header p {
                font-size: 12px
            }
            .main .content {
                font-size: 14px;
            }
        }

    </style>
</head>

<body style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
    <div class="main" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:100%;max-width:800px;margin:auto;border:1px solid #ddd;background-color:white;color:#424242;">
        <div class="header" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;color:white;background-color:#2d3e50;text-align:center;padding:16px;">
            <img src="https://[[.Domain]]/res/logo-email.jpg" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:80%;max-width:200px;margin-bottom:8px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;font-weight:600;font-size:14px;">Sistem Manajemen Surat</p>
        </div>
        <div class="content" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;padding:16px;font-size:16px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;margin-bottom:16px;">Anda telah didaftarkan ke Sistem Manajemen Surat Fakultas Teknik UPR dengan data sebagai berikut :</p>
            <table style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;border:none;outline:none;border-collapse:collapse;width:100%;margin-bottom:16px;">
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px;">Nama</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;width: 100%;">[[.Account.Nama]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Jabatan</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Account.Jabatan]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Telepon</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Account.Telepon]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Admin</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[if eq .Account.Admin 1]]Ya[[else]]Tidak[[end]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Keterangan</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[if eq .Account.Penginput 1]]Bisa menambah surat[[else]]-[[end]]</td>
                </tr>
            </table>
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;margin-bottom:16px;">Untuk mengakses Sistem Manajemen Surat, silakan login ke <a href="https://%5B%5B.Domain%5D%5D" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">[[.Domain]]</a> :</p>
            <table style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;border:none;outline:none;border-collapse:collapse;width:100%;margin-bottom:16px;">
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Email</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;width: 100%;">[[.Account.Email]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Password</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Password]]</td>
                </tr>
            </table>
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;">Setelah login, segera ganti password untuk mengamankan akun anda.</p>
        </div>
    </div>
</body>

</html>`

const emailNewSurat string = `
<!DOCTYPE html>
<html>

<head>
    <title>SIMAS - Sistem Manajemen Surat</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:400,600">
    <style type="text/css">
        @media screen and (max-width: 500px) {
            .main .header p {
                font-size: 12px
            }
            .main .content {
                font-size: 14px;
            }
        }

    </style>
</head>

<body style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
    <div class="main" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:100%;max-width:800px;margin:auto;border:1px solid #ddd;background-color:white;color:#424242;">
        <div class="header" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;color:white;background-color:#2d3e50;text-align:center;padding:16px;">
            <img src="https://[[.Domain]]/res/logo-email.jpg" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:80%;max-width:200px;margin-bottom:8px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;font-weight:600;font-size:14px;">Sistem Manajemen Surat</p>
        </div>
        <div class="content" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;padding:16px;font-size:16px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;margin-bottom:16px;">Anda mendapat surat baru dengan data sebagai berikut :</p>
            <table style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;border:none;outline:none;border-collapse:collapse;width:100%;margin-bottom:16px;">
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; min-width: 100px;">Nomor</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;width: 100%;">[[.Surat.Nomor]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Perihal</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Surat.Perihal]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Sumber surat</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Surat.Sumber]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Prioritas</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[if eq .Surat.Prioritas 0]]Biasa[[else if eq .Surat.Prioritas 1]]Segera[[else if eq .Surat.Prioritas 2]]Penting[[else]]Rahasia[[end]]</td>
                </tr>
            </table>
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;">Untuk mengakses Sistem Manajemen Surat, silakan login ke <a href="https://%5B%5B.Domain%5D%5D" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">[[.Domain]]</a></p>
        </div>
    </div>
</body>

</html>`

const emailNewDisposisi string = `
<!DOCTYPE html>
<html>

<head>
    <title>SIMAS - Sistem Manajemen Surat</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:400,600">
    <style type="text/css">
        @media screen and (max-width: 500px) {
            .main .header p {
                font-size: 12px
            }
            .main .content {
                font-size: 14px;
            }
        }

    </style>
</head>

<body style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
    <div class="main" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:100%;max-width:800px;margin:auto;border:1px solid #ddd;background-color:white;color:#424242;">
        <div class="header" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;color:white;background-color:#2d3e50;text-align:center;padding:16px;">
            <img src="https://[[.Domain]]/res/logo-email.jpg" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:80%;max-width:200px;margin-bottom:8px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;font-weight:600;font-size:14px;">Sistem Manajemen Surat</p>
        </div>
        <div class="content" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;padding:16px;font-size:16px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;margin-bottom:16px;">Anda mendapat disposisi dari <b style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">[[.Sumber]]</b> ([[.Jabatan]]) surat dengan data sebagai berikut :</p>
            <table style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;border:none;outline:none;border-collapse:collapse;width:100%;margin-bottom:16px;">
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; min-width: 100px;">Nomor</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;width: 100%;">[[.Surat.Nomor]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Perihal</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Surat.Perihal]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Sumber surat</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Surat.Sumber]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Prioritas</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[if eq .Surat.Prioritas 0]]Biasa[[else if eq .Surat.Prioritas 1]]Segera[[else if eq .Surat.Prioritas 2]]Penting[[else]]Rahasia[[end]]</td>
                </tr>
            </table>
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;">Untuk mengakses Sistem Manajemen Surat, silakan login ke <a href="https://%5B%5B.Domain%5D%5D" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">[[.Domain]]</a></p>
        </div>
    </div>
</body>

</html>`

const emailNewStatus string = `
<!DOCTYPE html>
<html>

<head>
    <title>SIMAS - Sistem Manajemen Surat</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Sans+Pro:400,600">
    <style type="text/css">
        @media screen and (max-width: 500px) {
            .main .header p {
                font-size: 12px
            }
            .main .content {
                font-size: 14px;
            }
        }

    </style>
</head>

<body style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
    <div class="main" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:100%;max-width:800px;margin:auto;border:1px solid #ddd;background-color:white;color:#424242;">
        <div class="header" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;color:white;background-color:#2d3e50;text-align:center;padding:16px;">
            <img src="https://[[.Domain]]/res/logo-email.jpg" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;width:80%;max-width:200px;margin-bottom:8px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;font-weight:600;font-size:14px;">Sistem Manajemen Surat</p>
        </div>
        <div class="content" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;padding:16px;font-size:16px;">
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;margin-bottom:16px;">Surat dengan data sebagai berikut :</p>
            <table style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;border:none;outline:none;border-collapse:collapse;width:100%;margin-bottom:16px;">
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; min-width: 100px;">Nomor</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;width: 100%;">[[.Nomor]]</td>
                </tr>
                <tr style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">Perihal</th>
                    <th style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;padding-right: 16px; padding-left: 16px;">:</th>
                    <td style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;text-align:left;vertical-align:top;padding-bottom:16px;">[[.Perihal]]</td>
                </tr>
            </table>
            <p style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;line-height:1.5em;word-break:break-all;word-wrap:break-word;">Telah selesai <b style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">[[if eq .Status 1]]diarsipkan[[else if eq .Status 2]]ditindaklanjuti[[else]]diarsip dan ditindaklanjuti[[end]]</b>. Untuk mengakses Sistem Manajemen Surat, silakan login ke <a href="https://%5B%5B.Domain%5D%5D" style="margin:0;padding:0;box-sizing:border-box;font-family:'Source Sans Pro', sans-serif;">[[.Domain]]</a></p>
        </div>
    </div>
</body>

</html>`
