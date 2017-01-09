<p align="center"><a href="https://github.com/RadhiFadlillah/simas" _target="blank"><img width="400" src="https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/logo.png"></img></a></p>

SIMAS adalah aplikasi Sistem Manajemen Surat yang dikembangkan untuk Fakultas Teknik Universitas Palangkaraya. Aplikasi ini berbasis web dan dikembangkan menggunakan bahasa pemrograman Go untuk back-end dan framework Vue.js untuk front-end.

## Latar Belakang

Sebelumnya, sistem persuratan di Fakultas Teknik UPR secara sederhana bekerja dengan cara sebagai berikut :

1. Pihak luar datang membawa surat dan menyerahkannya ke staff Dekanat.

   ![Proses lama #1](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-1.png)
   
2. Staff Dekanat menyerahkan surat yang diterima ke Dekan.

   ![Proses lama #2](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-2.png)
   
3. Dekan membaca surat yang diterima, lalu membuat disposisi surat yang ditujukan kepada bawahannya (sebagai contoh di sini adalah Wakil Dekan).

   ![Proses lama #3](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-3.png)
   
4. Dekan menyerahkan surat dan disposisinya ke staff Dekanat.

   ![Proses lama #4](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-4.png)
   
5. Staff Dekanat memfotokopi surat, lalu menyerahkan surat dan disposisinya ke tujuan, yaitu Wakil Dekan.

   ![Proses lama #5](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-5.png)
   
6. Wakil Dekan melaksanakan isi surat, lalu menuliskan laporannya ke lembar disposisi, apakah diarsipkan atau ditindaklanjuti.

   ![Proses lama #6](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-6.png)
   
7. Wakil Dekan mengembalikan surat dan lembar disposisinya ke staff Dekanat.

   ![Proses lama #7](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-7.png)
   
8. Staff Dekanat menyerahkan surat dan lembar disposisinya ke Dekan.

   ![Proses lama #8](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-8.png)
   
9. Dekan mengetahui bahwa surat tersebut sudah dilaksanakan atau diarsipkan.

   ![Proses lama #9](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-lama-9.png)

Proses lama ini memiliki beberapa masalah, yaitu :

1. Terlalu banyak kertas yang digunakan. Seandainya Dekan mendisposisikan surat ke tiga tujuan, maka surat tersebut akan difotokopi tiga kali, dan surat disposisi yang dibuat oleh Dekan juga harus sebanyak tiga buah.
2. Proses penyerahan disposisi dan laporan hasil disposisi seluruhnya melalui staff Dekanat, sehingga semakin banyak surat yang masuk, semakin berat tugas untuk staff Dekanat.
3. Dekan tidak dapat mengetahui apa yang terjadi terhadap surat yang telah dia disposisikan, sampai dia mendapat laporan akhir dari staff Dekanat. Jadi dia tidak tahu apakah suratnya didisposisikan lagi ke orang lain, atau sudah dilaksanakan.

## Cara Kerja Aplikasi

Dari latar belakang di atas, aplikasi SIMAS dirancang untuk bekerja sebagai berikut :

1. Pihak luar datang membawa surat dan menyerahkannya ke staff Dekanat.

   ![Proses baru #1](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-1.png)
   
2. Staff Dekanat menscan surat yang masuk. Setelah itu mengupload hasil scan dan data surat ke dalam aplikasi.

   ![Proses baru #2](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-2.png)
   
3. Aplikasi mengirimkan notifikasi melalui SMS dan email kepada Dekan bahwa ada surat baru yang masuk.

   ![Proses baru #3](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-3.png)
   
4. Dekan membuka aplikasi dan membaca surat, lalu mengklik tombol disposisi surat. Sebagai contoh, di sini Dekan mendisposisikan surat ke Wakil Dekan.

   ![Proses baru #4](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-4.png)
   
5. Aplikasi mengirimkan notifikasi melalui SMS dan email kepada Wakil Dekan bahwa dia mendapat disposisi surat dari Dekan.

   ![Proses baru #5](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-5.png)
   
6. Wakil Dekan membuka aplikasi, lalu membaca surat yang didisposisikan kepadanya. Setelah selesai dia laksanakan, dia mengklik tombol tindaklanjuti surat di aplikasi.

   ![Proses baru #6](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-6.png)
   
7. Aplikasi mengirimkan notifikasi bahwa surat sudah ditindaklanjuti oleh Wakil Dekan.

   ![Proses baru #7](https://raw.githubusercontent.com/RadhiFadlillah/simas/master/readme/diagram/proses-baru-7.png)

Dengan menggunakan aplikasi ini, ketiga masalah lama dapat teratasi :

1. Tidak lagi perlu memfotokopi surat, sebab data surat telah di-scan dan diupload ke aplikasi. Begitu juga surat disposisi tidak perlu lagi dibuat sebab sudah melalui aplikasi.
2. Tugas staff Dekanat sekarang hanya menscan dan mengupload data surat ke aplikasi. Sedangkan pemberitahuan bahwa ada surat baru, ada disposisi yang masuk dan lain sebagainya telah ditangani oleh aplikasi.
3. Dekan dapat kapan saja melihat apakah suratnya didisposisikan lagi, sudah diarsip atau sudah ditindaklanjuti.

## Screenshot

#### Desktop

|![Halaman login](https://github.com/RadhiFadlillah/simas/blob/master/readme/screenshot/desktop-1.png)|![Halaman utama](https://github.com/RadhiFadlillah/simas/blob/master/readme/screenshot/desktop-2.png)|
|:-:|:-:|

#### Mobile

|![Halaman login](https://github.com/RadhiFadlillah/simas/blob/master/readme/screenshot/mobile-1.png)|![Halaman daftar surat](https://github.com/RadhiFadlillah/simas/blob/master/readme/screenshot/mobile-2.png)|![Halaman detail surat](https://github.com/RadhiFadlillah/simas/blob/master/readme/screenshot/mobile-3.png)|
|:-:|:-:|:-:|

## Lisensi

SIMAS disebarkan dengan lisensi [MIT](http://choosealicense.com/licenses/mit/) sehingga setiap orang dapat dengan bebas mendownload, menggunakan dan memodifikasi aplikasi ini, dengan syarat tetap memberikan atribusi kepada [Radhi Fadlillah](https://github.com/RadhiFadlillah).

## Kredit

Ikon dan yang digunakan untuk illustrasi proses kerja di atas adalah karya [Freepik](http://www.freepik.com/) dan didownload di [Flaticon](www.flaticon.com)
