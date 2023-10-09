// package main

// import (
// 	"fmt"
// 	"go-dbms/config"
// 	"go-dbms/server"
// 	"go-dbms/services"
// )

// func main() {
// 	configs := config.New()
// 	svcs := services.New()

// 	err := server.Start(configs.ServerConfig, svcs)
// 	fmt.Printf(err)
// }

// package main

// import (
// 	"fmt"
// 	r "math/rand"
// 	"os"
// 	"path"
// 	"time"

// 	"go-dbms/pkg/column"
// 	"go-dbms/pkg/table"
// 	"go-dbms/pkg/types"

// 	"github.com/sirupsen/logrus"
// )

// var rand = r.New(r.NewSource(time.Now().Unix()))

// func main() {
// 	logrus.SetLevel(logrus.DebugLevel)

// 	dir, _ := os.Getwd()
// 	tablePath := path.Join(dir, "testtable")
// 	var options *table.Options = nil

// 	options = &table.Options{
// 		Columns: []*column.Column{
// 			column.New("id",        types.Meta(types.TYPE_INTEGER, true, 4)),
// 			column.New("firstname", types.Meta(types.TYPE_VARCHAR, 32)),
// 			column.New("lastname",  types.Meta(types.TYPE_VARCHAR, 32)),
// 		},
// 	}

// 	table, err := table.Open(tablePath, options)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}

// 	start := time.Now()
// 	exitFunc := func() {
// 		fmt.Println("DURATION =>", time.Since(start))
// 		_ = table.Close()
// 		// os.Remove(path.Join(tablePath, "data.dat"))
// 		// os.RemoveAll(path.Join(tablePath, "indexes"))
// 	}
// 	logrus.RegisterExitHandler(exitFunc)
// 	defer exitFunc()

// 	// ptr, err := table.Insert(map[string]types.DataType{
// 	// 	"id":        types.Type(types.TYPE_INT).Set(int32(7)),
// 	// 	"firstname": types.Type(types.TYPE_STRING).Set("Vahag"),
// 	// 	"lastname":  types.Type(types.TYPE_STRING).Set("Zargaryan"),
// 	// })
// 	// if err != nil {
// 	// 	logrus.Fatal(err)
// 	// }

// 	// fmt.Printf("%s\n", ptr)
// 	// record, err := table.Get(ptr)
// 	// if err != nil {
// 	// 	logrus.Fatal(err)
// 	// }
// 	// printData(options.ColumnsOrder, [][]types.DataType{record})

// 	// err = table.FullScan(func(ptr *data.RecordPointer, row map[string]types.DataType) (bool, error) {
// 	// 	fmt.Printf("%s, %s", ptr, sprintData(table.Columns(), []map[string]types.DataType{row}))
// 	// 	return false, nil
// 	// })
// 	// if err != nil {
// 	// 	logrus.Fatal(err)
// 	// }

// 	err = table.CreateIndex(nil, []string{"id"}, false)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}
// 	err = table.CreateIndex(nil, []string{"firstname","lastname"}, false)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}

// 	rand.Seed(time.Now().Unix())
// 	ids      := []int{5,6,4,5,7,2,1,9}
// 	names    := []string{"Vahag",     "Sergey",    "Bagrat",   "Mery"}
// 	surnames := []string{"Zargaryan", "Voskanyan", "Galstyan", "Sargsyan"}
// 	for _, id := range ids {
// 		_, err := table.Insert(map[string]types.DataType{
// 			"id":        types.Type(table.ColumnsMap()["id"].Meta).Set(id),
// 			"firstname": types.Type(table.ColumnsMap()["firstname"].Meta).Set(names[rand.Int31n(4)]),
// 			"lastname":  types.Type(table.ColumnsMap()["lastname"].Meta).Set(surnames[rand.Int31n(4)]),
// 		})
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 	}

// 	// err = table.FullScanByIndex("id_1", false, func(row map[string]types.DataType) (bool, error) {
// 	// 	printData(table.Columns(), []map[string]types.DataType{row})
// 	// 	return false, nil
// 	// })
// 	// if err != nil {
// 	// 	logrus.Fatal(err)
// 	// }

// 	// TODO: handle case when count of duplicate entries in node doesn't fit in page
// 	// TODO: add freelist logic
// 	err = table.FullScanByIndex("firstname_lastname_1", false, func(row map[string]types.DataType) (bool, error) {
// 		printData(table.Columns(), []map[string]types.DataType{row})
// 		return false, nil
// 	})
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}

// 	// records, err := table.FindByIndex(
// 	// 	// "id_1",
// 	// 	"firstname_lastname_1",
// 	// 	"<=",
// 	// 	map[string]types.DataType{
// 	// 		// "id": types.Type(table.ColumnsMap()["id"].Meta).Set(5),
// 	// 		"firstname": types.Type(table.ColumnsMap()["firstname"].Meta).Set("Sergey"),
// 	// 		"lastname":  types.Type(table.ColumnsMap()["lastname"].Meta).Set("Zargaryan"),
// 	// 	},
// 	// )
// 	// if err != nil {
// 	// 	logrus.Fatal(err)
// 	// }
// 	// printData(table.Columns(), records)

// 	// for i := 0; i < 10; i++ {
// 	// 	record, err := table.FindByIndex("id_1", false, map[string]types.DataType{
// 	// 		"id": types.Type(types.TYPE_INT, table.ColumnsMap()["id"].Meta).Set(i),
// 	// 	})
// 	// 	if err != nil {
// 	// 		logrus.Error(err)
// 	// 		continue
// 	// 	}
// 	// 	printData(table.Columns(), record)
// 	// }
// }

package main

import (
	"encoding/binary"
	"fmt"
	"go-dbms/pkg/column"
	"go-dbms/pkg/rbtree"
	"go-dbms/pkg/types"
	r "math/rand"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

var rand = r.New(r.NewSource(time.Now().Unix()))

func main() {
	// logrus.SetLevel(logrus.DebugLevel)
	// pwd, _ := os.Getwd()

	// ll, err := freelist.Open(path.Join(pwd, "test", "freelist.bin"), &freelist.LinkedListOptions{
	// 	PageSize: uint16(os.Getpagesize()),
	// 	PreAlloc: 5,
	// 	ValSize:  8,
	// })
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// p, err := pager.Open(path.Join(pwd, "test", "test.dat"), os.Getpagesize(), false, 0664)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// subfl, err := freelist.Open(path.Join(pwd, "test", "freelist.bin"), &freelist.Options{
	// 	PreAlloc:         5,
	// 	TargetPageSize:   uint16(os.Getpagesize()),
	// 	FreelistPageSize: uint16(os.Getpagesize()),
	// })
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// var fl freelist.Freelist
	// tree, err := bptree.Open(path.Join(pwd, "test", "bptree_freelist.idx"), &bptree.Options{
	// 	ReadOnly:     false,
	// 	FileMode:     0664,
	// 	MaxKeySize:   10,
	// 	MaxValueSize: 0,
	// 	PageSize:     os.Getpagesize(),
	// 	PreAlloc:     10,
	// 	FreelistOptions: &freelist.Options{
	// 		Allocator:      p,
	// 		PreAlloc:       5,
	// 		TargetPageSize: uint16(os.Getpagesize()),
	// 	},
	// }, subfl)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// fl = tree

	// start := time.Now()
	// exitFunc := func() {
	// 	fmt.Println("TOTAL DURATION =>", time.Since(start))
	// 	_ = ll.Close()
	// 	// _ = p.Close()
	// }
	// logrus.RegisterExitHandler(exitFunc)
	// defer exitFunc()


	// for i := 1; i <= 10; i++ {
	// 	// p.Alloc(i)
	// 	_, err = fl.AddMem(uint64(i), uint16(rand.Intn(4096)))
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }
	// fmt.Println("ADD DURATION =>", time.Since(start))
	// start = time.Now()
	// if err := fl.WriteAll(); err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println("FLUSH DURATION =>", time.Since(start))

	// fmt.Println(fl.Get(&freelist.Pointer{
	// 	PageId: 4,
	// 	Index:  1,
	// }))

	// pageId, ptr, err := fl.Alloc(30)
	// fmt.Println(pageId, ptr, err)

	// if err = fl.Set(&bptree.Pointer{
	// 	FreeSpace: 1501,
	// 	PageId:    7,
	// }, 1000); err != nil {
	// 	logrus.Fatal(err)
	// }

	// for i := 0; i < 10; i++ {
	// 	val := make([]byte, 8)
	// 	binary.BigEndian.PutUint64(val, uint64(rand.Int63n(100)))
	// 	_, err := ll.Push(val)
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }
	
	// if err = ll.Print(); err != nil {
	// 	logrus.Fatal(err)
	// }

	// _, val, err := ll.Pop(2)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(val)

	// if err = ll.Print(); err != nil {
	// 	logrus.Fatal(err)
	// }


	logrus.SetLevel(logrus.DebugLevel)
	pwd, _ := os.Getwd()

	var err error
	// var arr array.ArrayADS[*array.Number[uint16], array.Number[uint16]]
	// arr, err = array.Open[*array.Number[uint16], array.Number[uint16]](
	// 	path.Join(pwd, "test", "array.bin"),
	// 	&array.ArrayOptions{
	// 		PageSize: uint16(os.Getpagesize()),
	// 		PreAlloc: 5,
	// 	},
	// )
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	
	file := path.Join(pwd, "test", "rbtree.bin")
	// os.Remove(file)
	t, err := rbtree.Open(
		file,
		&rbtree.Options{
			PageSize: uint16(os.Getpagesize()),
			KeySize:  19,
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}
	
	elems := make([]uint16, 0, 10)
	start := time.Now()
	exitFunc := func() {
		fmt.Println("\nTOTAL DURATION =>", time.Since(start))
		fmt.Println(elems)
		// _ = arr.Close()
		_ = t.Close()
	}
	logrus.RegisterExitHandler(exitFunc)
	defer exitFunc()

	// b := make([]byte, 19)
	// els := []uint16{
	// 	10,
	// 	20,
	// 	30,
	// 	100,
	// 	90,
	// 	40,
	// 	50,
	// 	60,
	// 	70,
	// 	80,
	// 	150,
	// 	110,
	// 	120,
	// }
	// for _, elem := range els {
	// 	binary.BigEndian.PutUint16(b, elem)
	// 	if err := t.Insert(b); err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }

	// if err := t.Print(5); err != nil {
	// 	logrus.Fatal(err)
	// }

	
	// db := make([]byte, 19)

	// binary.BigEndian.PutUint16(db, 60)
	// fmt.Println("deleting", 60)
	// if err := t.Delete(db); err != nil {
	// 	logrus.Fatal(err)
	// }

	// binary.BigEndian.PutUint16(db, 120)
	// fmt.Println("deleting", 120)
	// if err := t.Delete(db); err != nil {
	// 	logrus.Fatal(err)
	// }


	b := make([]byte, 19)
	for i := 0; i < 100; i++ {
		elem := uint16(rand.Int31n(256))
		elems = append(elems, elem)
		binary.BigEndian.PutUint16(b, elem)
		if err := t.InsertMem(b); err != nil {
			logrus.Fatal(i, err)
		}
	}
	if err := t.WriteAll(); err != nil {
		logrus.Fatal(err)
	}

	fmt.Println(elems)

	// gb := make([]byte, 19)
	// binary.BigEndian.PutUint16(gb, 242)
	// if v, err := t.Get(gb); err != nil {
	// 	logrus.Error(err)
	// } else {
	// 	fmt.Println(binary.BigEndian.Uint16(v))
	// }

	// elems = []uint16{0, 0, 6, 12, 15, 15, 17, 18, 19, 19, 24, 24, 25, 30, 32, 33, 33, 34, 34, 36, 39, 42, 42, 43, 47, 49, 51, 52, 61, 66, 73, 74, 76, 78, 79, 80, 80, 83, 85, 88, 89, 90, 91, 92, 93, 98, 99, 100, 101, 107, 112, 118, 122, 124, 126, 129, 132, 136, 140, 140, 142, 145, 145, 145, 172, 184, 189, 192, 192, 193, 195, 199, 200, 201, 202, 202, 202, 207, 210, 211, 212, 214, 216, 216, 219, 226, 227, 229, 233, 234, 235, 237, 238, 241, 245, 245, 245, 246, 254, 255}
	db := make([]byte, 19)
	for i := 0; i < 100; i++ {
		// index := rand.Intn(len(elems))
		elem := elems[i]
		// elems = append(elems[:index], elems[index+1:]...)
		binary.BigEndian.PutUint16(db, elem)
		fmt.Println("deleting", elem)
		if err := t.Delete(db); err != nil {
			logrus.Fatal(err)
		}
	}

	err = t.Scan(0, func(key []byte) (bool, error) {
		fmt.Printf("%d, ", binary.BigEndian.Uint16(key))
		return false, nil
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// err = t.Print()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// if err := arr.Truncate(2); err != nil {
	// 	logrus.Fatal(err)
	// }

	// for i := 0; i < 10; i++ {
	// 	n := array.NewNumber(uint16(rand.Int63n(5000)))
	// 	index, err := arr.PushMem(n)
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// 	fmt.Println(index, n.Value())
	// }

	// for arr.Size() > 0 {
	// 	itm, err := arr.PopMem()
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// 	fmt.Println(arr.Size(), itm.Value())
	// }

	// err = arr.Set(666666, num(555))
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// itm, err := arr.Get(0)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// fmt.Println(*itm)

	// if err := arr.Print(); err != nil {
	// 	logrus.Fatal(err)
	// }
}

func num(n uint64) *number {
	a := number(n)
	return &a
}
type number uint64
func (n *number) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(*n))
	return buf, nil
}
func (n *number) UnmarshalBinary(d []byte) error {
	*n = number(binary.BigEndian.Uint64(d[0:8]))
	return nil
}


func sprintData(columns []*column.Column, data []map[string]types.DataType) string {
	str := ""
	for _, d := range data {
		for _, col := range columns {
			str += fmt.Sprintf("'%s' -> '%v', ", col.Name, d[col.Name].Value())
		}
		str += "\n"
	}
	return str
}

func printData(columns []*column.Column, data []map[string]types.DataType) {
	fmt.Print(sprintData(columns, data))
}
