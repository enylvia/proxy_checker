# 🚀 Proxy Checker Go

Proxy checker super ringan dan kencang yang dibuat dengan Bahasa Pemrograman Go. Menggunakan Goroutines dan Worker Pool untuk pengecekan paralel yang efisien.

## ✨ Fitur
- **Parallel Checking**: Menggunakan Worker Pool (Goroutines) agar proses scan super cepat.
- **Multi-Protocol**: Mendukung pengecekan protokol **HTTP** dan **SOCKS5**.
- **Real-time Log**: Menampilkan status Alive/Dead dengan warna yang menarik di terminal.
- **Verification**: Menampilkan IP asli yang terdeteksi melalui proxy (bukti proxy bekerja).
- **Auto-Save**: Menyimpan hasil proxy yang aktif ke dalam file `active.txt`.
- **Command Line Flags**: Konfigurasi mudah tanpa ubah kode.

## 📋 Prasyarat
- [Go (Golang)](https://golang.org/dl/) versi 1.16 atau yang lebih baru.

## 🛠️ Cara Instalasi & Penggunaan

1. **Clone Repository**
   ```bash
   git clone https://github.com/enylvia/proxy_checker.git
   cd proxy_checker
   ```

2. **Siapkan List Proxy**
   Buat file bernama `proxy.txt` di direktori yang sama. Isi dengan daftar proxy format `IP:Port`.

3. **Jalankan Skrip**
   ```bash
   go run main.go
   ```

## ⚙️ Opsi Command Line
Lo bisa kustomisasi pengecekan lewat flags:
```bash
go run main.go -w 100 -t 15 -f list_proxy.txt
```
| Flag | Deskripsi | Default |
| :--- | :--- | :--- |
| `-w` | Jumlah worker paralel | `50` |
| `-t` | Timeout (detik) | `10` |
| `-f` | Nama file sumber proxy | `proxy.txt` |

## 📁 Struktur Folder
- `main.go`: Logika utama program.
- `proxy.txt`: File input (tidak masuk git).
- `active.txt`: File output hasil scan (tidak masuk git).
- `LICENSE`: Lisensi MIT.

---
Dibuat dengan ❤️ pake Go.
