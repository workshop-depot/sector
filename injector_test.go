package sector

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestInjectorSample(t *testing.T) {
	dj := NewInjector(FactoryFunc(genericFactory))

	var s Sample
	dj.Inject(&s)

	expected := `{Trait:{Val:Aloha Call-0002! Data:{Str:Aloha Call-0004! Num:5 NumList:[1 2 3]} S
I:{Val:SIO-0007}} N:8 S:Aloha Call-0009! D:{Str:Aloha Call-0011! Num:12 NumList:[1 2 3]} NL:[1 2 3] Map:map[a
sset:ready-0015] DataPtr:{Str:Aloha Call-0017! Num:18 NumList:[1 2 3]} SI:{Val:SIO-0020}}`

	expected = strings.Replace(expected, "\n", "", -1)

	actual := fmt.Sprintf("%+v", s)

	if actual != expected {
		t.Log(actual)
		t.Fail()
	}
}

var counter int

func genericFactory(ptr interface{}) bool {
	counter++

	switch x := ptr.(type) {
	case *int:
		*x = counter
	case *string:
		*x = fmt.Sprintf("Aloha Call-%04d!", counter)
	case *Trait:
		buffer := Trait{}
		*x = buffer
	case *Data:
		buffer := Data{}
		*x = buffer
	case *map[string]interface{}:
		*x = make(map[string]interface{})
		(*x)[`asset`] = `ready-` + fmt.Sprintf("%04d", counter)
	case *SI:
		v := SIO{`SIO-` + fmt.Sprintf("%04d", counter)}
		*x = &v
	case *[]int:
		*x = []int{1, 2, 3}
	default:
		return false
	}

	return true
}

//-----------------------------------------------------------------------------
// sample structs

type Sample struct {
	Trait `inject:"*"`

	N  int    `inject:"+"`
	S  string `inject:"+"`
	D  Data   `inject:"*"`
	NL []int  `inject:"+"`

	Map     map[string]interface{} `inject:"+"`
	DataPtr *Data                  `inject:"*"`

	SI SI `inject:"+"`
}

type SI interface {
	Pop()
}

type SIO struct {
	Val string
}

func (x *SIO) String() string {
	if x == nil {
		return NIL
	}

	return fmt.Sprintf("%+v", *x)
}

func (*SIO) Pop() { log.Println(`POPED`) }

type Trait struct {
	Val  string `inject:"+"`
	Data *Data  `inject:"*"`

	SI SI `inject:"*"`
}

func (x *Trait) String() string {
	if x == nil {
		return NIL
	}

	return fmt.Sprintf("%+v", *x)
}

type Data struct {
	Str     string `inject:"+"`
	Num     int    `inject:"+"`
	NumList []int  `inject:"+"`
}

func (x *Data) String() string {
	if x == nil {
		return NIL
	}

	return fmt.Sprintf("%+v", *x)
}

const NIL = `NIL`
