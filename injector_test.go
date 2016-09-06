package sector

import (
	"fmt"
	"log"
	"testing"
)

func TestInjectorSample(t *testing.T) {
	dj := NewInjector(FactoryFunc(genericFactory))

	var s Sample
	dj.Inject(&s)

	expected := `{Trait:{Val:Aloha Call-0002! Data:{Str: Num:0} SI:{Val:SIO-0004}} N:5 S:Aloha Call-0006! D:{Str:Aloha Call-0008! Num:9} Map:map[asset:ready-0010] DataPtr:{Str:Aloha Call-0012! Num:13} SI:{Val:SIO-0014}}`

	actual := fmt.Sprintf("%+v", s)

	if actual != expected {
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
	default:
		return false
	}

	return true
}

//-----------------------------------------------------------------------------
// sample structs

type Sample struct {
	Trait `inject:"*"`

	N int    `inject:"+"`
	S string `inject:"+"`
	D Data   `inject:"*"`

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
	Data *Data  `inject:"+"`

	SI SI `inject:"*"`
}

func (x *Trait) String() string {
	if x == nil {
		return NIL
	}

	return fmt.Sprintf("%+v", *x)
}

type Data struct {
	Str string `inject:"+"`
	Num int    `inject:"+"`
}

func (x *Data) String() string {
	if x == nil {
		return NIL
	}

	return fmt.Sprintf("%+v", *x)
}

const NIL = `NIL`
