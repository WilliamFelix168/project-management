package types

import (
	"database/sql/driver"
	"errors"  // untuk membuat error kustom
	"fmt"     // untuk formatting error
	"strings" // untuk operasi string: Trim, Split

	"github.com/google/uuid" // paket UUID pihak ketiga untuk parsing UUID
)

// UUIDArray adalah tipe custom yang merepresentasikan slice dari uuid.UUID
type UUIDArray []uuid.UUID

// Scan mengimplementasikan interface sql.Scanner sehingga tipe ini bisa
// dipakai langsung untuk memindai (scan) nilai yang dikembalikan oleh
// database (mis. kolom array UUID di PostgreSQL).
//
// Parameter 'value' bertipe interface{} karena contract sql.Scanner
// memberikan nilai dalam berbagai bentuk (mis. []byte, string, atau nil).
// Fungsi ini bertugas mengubah representasi database menjadi UUIDArray.
func (a *UUIDArray) Scan(value interface{}) error {
	// Contoh bentuk value yang diharapkan: "{uuid1,uuid2,uuid3}"
	var str string

	// Tentukan tipe konkret dari value
	switch v := value.(type) {
	case []byte:
		// Jika database mengembalikan []byte, ubah menjadi string
		str = string(v)
	case string:
		// Jika sudah string, gunakan langsung
		str = v
	default:
		// Tipe lain tidak didukung — kembalikan error
		return errors.New("failed to parse UUIDArray : unsupport data type")
	}

	// Hapus kurung kurawal pembuka '{' jika ada di awal string
	str = strings.TrimPrefix(str, "{")
	// Hapus kurung kurawal penutup '}' jika ada di akhir string
	str = strings.TrimSuffix(str, "}")
	// Pisahkan string berdasarkan koma menjadi potongan-potongan UUID
	parts := strings.Split(str, ",")

	// Inisialisasi slice tujuan dengan kapasitas awal sama dengan jumlah part
	// Ini efisien karena menghindari alokasi ulang yang berlebihan
	*a = make(UUIDArray, 0, len(parts))

	// Iterasi setiap potongan hasil split
	for _, s := range parts {
		// Hapus spasi di kedua ujung dan juga hapus tanda kutip jika ada
		s = strings.TrimSpace(strings.Trim(s, `"`)) // trim spasi dan kutip

		// Jika setelah trimming string kosong, lewati (mungkin nilai NULL)
		if s == "" {
			continue
		}

		// Parse string menjadi uuid.UUID menggunakan library google/uuid
		u, err := uuid.Parse(s)
		if err != nil {
			// Jika parse gagal, kembalikan error berisi detail
			return fmt.Errorf("invalid UUID in Array : %v", err)
		}

		// Tambahkan UUID yang berhasil di-parse ke slice hasil
		*a = append(*a, u)
	}

	// Berhasil, kembalikan nil error
	return nil
}

// Kenapa ga pakai pointer di receiver?
// kalau kita pakai pointer, kita bisa mengubah isi dari UUIDArray
// karena kita tidak mengubah isi dari UUIDArray
// hanya mengembalikan nilai string di PostreSQL

// Value mengimplementasikan interface driver.Valuer agar UUIDArray bisa
// otomatis dikonversi ke format string array PostgreSQL saat disimpan ke database.
// Fungsi ini dipakai saat ingin menyimpan data ke database.
func (a UUIDArray) Value() (driver.Value, error) {
	// Jika slice kosong, kembalikan string array kosong PostgreSQL
	if len(a) == 0 {
		// Format array kosong di PostgreSQL adalah "{}"
		return "{}", nil
	}

	// postgreFormat akan menampung string UUID hasil konversi
	postgreFormat := make([]string, 0, len(a))
	// Loop setiap elemen UUIDArray
	for _, u := range a {
		// Konversi UUID ke string dan masukkan ke slice postgreFormat
		postgreFormat = append(postgreFormat, fmt.Sprintf(`"%s"`, u.String()))
	}
	// Gabungkan semua string UUID dengan koma, lalu bungkus dengan kurung kurawal
	// Contoh hasil: "{uuid1,uuid2,uuid3}"
	return "{" + strings.Join(postgreFormat, ",") + "}", nil
}


//dibuat gorm tipe custom sebagai uuid[]
//dipanggil otomtatis oleh gorm
//tujuannya agar gorm tau tipe data apa yang digunakan
func (UUIDArray) GormDataType() string {
	return "uuid[]"
}