package transaction

import (
	"container/list"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var priv ecdsa.PrivateKey



func TestHashTx(t *testing.T) {
	/*tx := NewTx(
		[]byte("444444333asdf"),
		[][]byte{[]byte("sdfkjlwke"), []byte("sdfkjlwke")},
		[]byte("12E9tH5tL2td8jgkrK8WCtkEJpCRgHDXpX"),
		[]byte("This road is fine"),
		[]byte("luoyu"),
		)*/

	tx := NewTx(
		[]byte("444444333asdfsdfsdf"),
		[][]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
	)

	tx2 := NewTx(
		[]byte("sdfsdfsdf"),
		[][]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
	)

	tx3 := &Transaction{
		[]byte("sdfsdfsdf"),
		[]byte{},
		[][]byte{},
		[]byte{},
		[]byte{},
		nil,
		12,
		1231232312321,
		nil,
	}

	hash := tx.HashTx()
	hash2 := tx2.HashTx()
	hash3 := tx3.HashTx()

	t.Log("the raw byte of hash: ",hash)


	// not good
	/*hashSpf := fmt.Sprintf("hash of tx is: %x", string(hash[:]))

	t.Log("sprintf: ", hashSpf)

	hashByteSpf := []byte(hashSpf)

	t.Log("hashByteSpf: ", hashByteSpf)*/


	// good
	hashEnc := hex.EncodeToString(hash)
	hashEnc2 := hex.EncodeToString(hash2)
	hashEnc3 := hex.EncodeToString(hash3)

	t.Log("EncodeToString: ", hashEnc)
	t.Log("EncodeToString: ", hashEnc2)
	t.Log("EncodeToString: ", hashEnc3)

	hashByteDec, _ := hex.DecodeString(hashEnc)

	t.Log("DecodeString: ", hashByteDec)
}

func TestTransaction_Serialize(t *testing.T) {
	/*tx := NewTx(
		[]byte("444444333asdf"),
		[][]byte{[]byte("sdfkjlwke"), []byte("sdfkjlwke")},
		[]byte("12E9tH5tL2td8jgkrK8WCtkEJpCRgHDXpX"),
		[]byte("fine"),
		[]byte("luoyu"),
	)*/

	// TODO bugs deep copy, the difference of twice hash, run twice

	tx := NewTx(
		[]byte("444444333asdf"),
		[][]byte{[]byte("4df"),[]byte("df"),},
		[]byte{},
		[]byte{},
		[]byte{},
	)

	t.Log(tx.Serialize())

	byt := tx.Serialize()

	fmt.Println("tx below")
	fmt.Println(fmt.Sprintf("%+v", tx))
	fmt.Println(fmt.Sprintf("%+v", DeserializeTx(byt)))

	var byt2 []byte

	copy(byt2, byt)

	t.Log(byt)

	//b := []byte{127, 255, 129, 3, 1, 1, 11, 84, 114, 97, 110}

	t.Log("sha256 below")
	//t.Log(sha256.Sum256(tx.Serialize()))
	t.Log(sha256.Sum256(byt2))
	t.Log(sha256.Sum256(byt2))

	t.Log("hash below")
	t.Log(tx.HashTx())

	/*valid := make([]byte, 2)

	var encode bytes.Buffer

	enc := gob.NewEncoder(&encode)
	err := enc.Encode(valid)

	if err != nil {
		log.Panic("tx encode fail:", err)
	}

	t.Log("[][]byte")
	//t.Log(sha256.Sum256(encode.Bytes()))
	t.Log(sha256.Sum256(valid))*/


}

func TestTransaction_Sign(t *testing.T) {
	/*tx := NewTx(
		[]byte("444444333asdf"),
		[][]byte{[]byte("4df"),[]byte("df"),},
		[]byte{},
		[]byte{},
		[]byte{},
	)*/

	//tx.Sign()

	bt := []byte("The road is fine")

	fmt.Printf("Transaction Info: %s", bt)
}

func TestDeserializeTx(t *testing.T) {
	tx := []int{1,2,3,4,5,6,7,8}
	//tx = append(tx[:2], tx[3:]...)
	i := 9
	tx = append(tx, i)
	tx = append(tx, i)
	fmt.Println(tx)
	fmt.Println(len(tx))
}

func TestPrune(t *testing.T) {

	tx := NewTx(
		[]byte("444444333asdfsdfsdf"),
		[][]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
	)

	tx2 := NewTx(
		[]byte("sdfsdfsdf"),
		[][]byte{},
		[]byte{},
		[]byte{},
		[]byte{},
	)

	hash := tx.HashTx()
	hash2 := tx2.HashTx()

	slideWindow := list.New()

	slideWindow.PushBack(hash)
	slideWindow.PushBack(hash2)

	first := slideWindow.Front()
	if first == nil {
		t.Fatal("err")
	}

	rmHashInf := slideWindow.Remove(first)

	//rmHash, errBy := GetBytes(rmHashInf)
	rmHash := rmHashInf.([]byte)

	//if errBy != nil {
	//	t.Fatal("err")
	//}

	fmt.Println(hash)
	fmt.Println(rmHashInf)
	fmt.Println(rmHash)

	assert.Equal(t, hash, rmHashInf, "The two bytes should be equal")
	assert.Equal(t, hash, rmHash, "The two bytes should be equal")
	assert.Equal(t, rmHashInf, rmHash, "The two bytes should be equal")

}

/*func GetBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}*/

func TestSlide(t *testing.T) {
	s := list.New()

	for i := 0; i < 100; i++ {
		PushSlide(s, i)
		if i == 48 {
			fmt.Println(s.Len())
		}

		if i == 49 {
			fmt.Println(s.Len())
		}

		if i == 50 {
			fmt.Println(s.Len())
		}

		if i == 51 {
			fmt.Println(s.Len())
		}

		if i == 98 {
			fmt.Println("length: ", s.Len())
			fmt.Println("value: ")
			for e := s.Front(); e != nil; e = e.Next() {
				fmt.Println(e.Value.(int))
			}
		}
	}

}

func PushSlide(slideWindow *list.List, v interface{}) {
	if slideWindow.Len() < 50 {
		slideWindow.PushBack(v)
	} else {
		first := slideWindow.Front()
		if first == nil {

		}

		slideWindow.Remove(first)

		//rmHash := rmHashInf.([]byte)



		slideWindow.PushBack(v)
	}
}