package sector

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInterfaces(t *testing.T) {
	// it is nice & important that FactoryRepo can be used as just Factory
	var fac Factory
	fac = FactoryFunc(func(interface{}) bool {
		return false
	})
	fac = NewFactoryRepo()
	_ = fac

	var j Injector
	j = NewInjector(fac)
	_ = j
}

func TestFill(t *testing.T) {
	fp := NewFactoryRepo()
	fp.Register(FactoryFunc(SomethingerFactory))

	var buffer Somethinger
	slot := &buffer
	ok := fp.Fill(slot)
	if !ok {
		t.Fail()
	}

	(*slot).Something(`OK`)
	str := fmt.Sprintf("%v", *slot)
	if !strings.Contains(str, `OK`) {
		t.Fail()
	}
}

func TestFillStruct(t *testing.T) {
	fp := NewFactoryRepo()
	fp.Register(FactoryFunc(SomethingerFactory))

	a := someAction{}
	slot := &a.sm

	ok := fp.Fill(slot)
	if !ok {
		t.Fail()
	}

	(*slot).Something(`OK 1`)
	str := fmt.Sprintf("%v", *slot)
	if !strings.Contains(str, `OK 1`) {
		t.Fail()
	}

	str = fmt.Sprintf("%v", a.sm)
	if !strings.Contains(str, `OK 1`) {
		t.Fail()
	}

	slot = &a.sm2

	ok = fp.Fill(slot)
	if !ok {
		t.Fail()
	}

	(*slot).Something(`OK 2`)
	str = fmt.Sprintf("%v", *slot)
	if !strings.Contains(str, `OK 2`) {
		t.Fail()
	}

	str = fmt.Sprintf("%v", a.sm2)
	if !strings.Contains(str, `OK 2`) {
		t.Fail()
	}

	a.sm.Something(`+I)`)
	a.sm2.Something(`II)`)
	str = fmt.Sprintf("%v %v", a.sm, a.sm2)
	if !strings.Contains(str, `+I)`) || !strings.Contains(str, `II)`) {
		t.Fail()
	}
}

func TestFillALL(t *testing.T) {
	fp := NewFactoryRepo()
	fp.Register(FactoryFunc(SomethingerFactory))

	var buffer1 Somethinger
	slot1 := &buffer1
	var buffer2 Somethinger
	slot2 := &buffer2
	var buffer3 Somethinger
	slot3 := &buffer3

	fp.FillAll(slot1, slot2, slot3)

	all := []*Somethinger{slot1, slot2, slot3}

	for _, v := range all {
		if v == nil {
			t.FailNow()
		}
	}

	for k, v := range all {
		ix := strconv.Itoa(k)
		slot := v
		s := `OK,` + ix
		(*slot).Something(s)
		str := fmt.Sprintf("%v", *slot)
		if !strings.Contains(str, s) {
			t.FailNow()
		}
	}
}

type someAction struct {
	sm  Somethinger
	sm2 Somethinger
}

//-----------------------------------------------------------------------------

// Somethinger /
type Somethinger interface {
	Something(string)
}

// ConcreteSomethinger /
type ConcreteSomethinger struct {
	State string
}

// Something /
func (v *ConcreteSomethinger) Something(s string) {
	v.State = s
}

// SomethingerFactory template factory, v is pointer to interface
func SomethingerFactory(v interface{}) bool {
	switch x := v.(type) {
	case *Somethinger:
		*x = new(ConcreteSomethinger)
		return true
	}

	return false
}

func T() struct {
	N int
} {
	type F struct {
		N int
	}

	return F{}
}

//-----------------------------------------------------------------------------
// another complete sample (test)

func TestFactorySample(t *testing.T) {
	t.Skip()

	factory := NewFactoryRepo()
	factory.Register(FactoryFunc(NotifierFactory))
	factory.Register(FactoryFunc(ConfigFactory))
	// factory.Register(FactoryFunc(GenericFactory))

	x := useCase{}
	factory.FillAll(&x.config, &x.notifier)

	log.Printf("%+v", x)
	x.notifier.Notify(`OK`)
}

//-----------------------------------------------------------------------------
// our struct

type useCase struct {
	notifier Notifier
	config   Config
}

//-----------------------------------------------------------------------------
// factory methods

func NotifierFactory(ptr interface{}) bool {
	switch x := ptr.(type) {
	case *Notifier:
		buffer := &DefaultNotifier{}
		ConfigFactory(&buffer.config)
		*x = buffer
		return true
	}
	return false
}

func ConfigFactory(ptr interface{}) bool {
	switch x := ptr.(type) {
	case *Config:
		*x = Config{}
		x.Delay = time.Second * 30
		x.Home = "/"
		x.UserPass.User = `usr`
		x.UserPass.Pass = `1234`
		return true
	}
	return false
}

func GenericFactory(ptr interface{}) (interface{}, bool) {
	switch x := ptr.(type) {
	case *Config:
		*x = Config{}
		x.Delay = time.Second * 30
		x.Home = "/"
		x.UserPass.User = `usr`
		x.UserPass.Pass = `1234`
		return x, true
	case *Notifier:
		*x = &DefaultNotifier{}
		return x, true
	}

	return nil, false
}

//-----------------------------------------------------------------------------
// we want to create Notifier and Config

type Notifier interface {
	Notify(interface{})
}

type DefaultNotifier struct {
	config Config
}

func (x *DefaultNotifier) Notify(msg interface{}) {
	log.Println(`alert:`, msg, `in`, x.config.Home)
}

type Config struct {
	Home     string
	Delay    time.Duration
	UserPass struct {
		User, Pass string
	}
}

//-----------------------------------------------------------------------------
