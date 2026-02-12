package utils

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/google/uuid"
)

func SortListsByPosition(lists []models.List, order []uuid.UUID) []models.List {
	sortedLists := make([]models.List, 0, len(lists))
	listMap := make(map[uuid.UUID]models.List)

	//untuk mapping list by publicID
	for _, list := range lists {
		// mapping list by publicID
		listMap[list.PublicID] = list
	}

	//untuk sorting berdasarkan order
	for _, pos := range order {
		if list, exists := listMap[pos]; exists {
			sortedLists = append(sortedLists, list)
		}
	}

	//contoh
	// order: [id3, id1, id2]
	// lists: [{id1}, {id2}, {id3}]
	// result: [{id3}, {id1}, {id2}]

	return sortedLists
}

/*
Oke, aku jelasin fungsi ini step-by-step ya:

Tujuan Fungsi
Fungsi ini mengubah urutan list sesuai dengan order yang diberikan.

Penjelasan Detail:
func SortListsByPosition(lists []models.List, order []uuid.UUID) []models.List

Input 1: lists = daftar list yang belum terurut 
Input 2: order = urutan ID yang diinginkan 
Output: list yang sudah terurut sesuai order


Langkah 1: Buat Map
sortedLists := make([]models.List, 0, len(lists))
listMap := make(map[uuid.UUID]models.List)
sortedLists = array kosong untuk hasil akhir listMap = dictionary untuk akses cepat

Langkah 2: Isi Map
for _, list := range lists {
    listMap[list.PublicID] = list
}
Ubah list menjadi key-value:
{id1 → {name: "Todo"}},
{id2 → {name: "In Progress"}},
{id3 → {name: "Done"}}


Langkah 3: Susun Urutan
for _, pos := range order {
    if list, exists := listMap[pos]; exists {
        sortedLists = append(sortedLists, list)
    }
}

Ambil dari map sesuai urutan di order:

Contoh:

order = [id3, id1, id2] 
Ambil id3 → tambah ke hasil 
Ambil id1 → tambah ke hasil 
Ambil id2 → tambah ke hasil 
Hasil: urutan berubah sesuai keinginan ✅
Gampangnya: Fungsi ini seperti playlist - list awal acak, order adalah urutan putar yang diinginkan.

Ada yang masih bingung?
*/