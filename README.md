# sector
*sector* - for Simple Injector - provides a Dependency Injection mechanism for Go.

Put it simply, _sector_ fills pointers with values come from _factories_. So we have the *factory* - the constructor (role) - and the *injector*.

A *factory* implements the `Factory` interface. In it's simplest form, it's just a function with this signature, `func(interface{}) bool` and uses the Go's _type switch_ to fill the pointers:

```go
func myFactory(ptr interface{}) bool {
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
```

Then we use this factory in our _injector_ to fill in the fields of a struct.

```go
dj := NewInjector(FactoryFunc(genericFactory))
```

Now we use this injector to inject desired values into struct's fields.

```go
var s Sample
dj.Inject(&s)
```

Fields that are tagged with `inject:"+"` will get filled with proper value from the _factory_. Fields tagged with `inject:"*"` will get filled recursively down the tree.

```go
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
```