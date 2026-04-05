# 🚀 Proxy Checker Go

Proxy checker super ringan dan kencang yang dibuat dengan Bahasa Pemrograman Go. Menggunakan Goroutines dan Worker Pool untuk pengecekan paralel yang efisien.

## ✨ Fitur
- **Parallel Checking**: Menggunakan Worker Pool (Goroutines) agar proses scan super cepat.
- **Multi-Protocol**: Mendukung pengecekan protokol **HTTP** dan **SOCKS5**.
- **Output Bersih**: Menampilkan status Alive/Dead lengkap dengan *Response Time*.
- **Auto-Save**: Menyimpan hasil proxy yang aktif ke dalam file `active.txt`.
- **Lightweight**: Hemat memori dan CPU.

## 📋 Prasyarat
- [Go (Golang)](https://golang.org/dl/) versi 1.16 atau yang lebih baru.

## 🛠️ Cara Instalasi & Penggunaan

1. **Clone Repository**
   ```bash
   git clone https://github.com/enylvia/proxy_checker.git
   cd proxy_checker
   ```

2. **Siapkan List Proxy**
   Buat file bernama `proxy.txt` di direktori yang sama. Isi dengan daftar proxy format `IP:Port`. Contoh:
   ```text
   127.0.0.1:8080
   1.2.3.4:1080
   ```

3. **Jalankan Skrip**
   ```bash
   go run main.go
   ```

## ⚙️ Konfigurasi
Lo bisa ubah variabel di dalam `main.go` sesuai kebutuhan:
- `workerCount`: Jumlah worker paralel (default: 50).
- `timeout`: Waktu tunggu respon (default: 5 detik).
- `proxyFile`: Nama file sumber (default: `proxy.txt`).
- `activeFile`: Nama file hasil (default: `active.txt`).

## 📁 Struktur Folder
- `main.go`: Logika utama program.
- `proxy.txt`: File input (tidak masuk git).
- `active.txt`: File output hasil scan (tidak masuk git).
- `.gitignore`: Mengatur file yang tidak perlu di-upload ke GitHub.

---
Dibuat dengan ❤️ pake Go.
