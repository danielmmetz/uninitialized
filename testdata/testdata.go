package testfile

type Foo struct {
	Bar      `required:"true"`
	NamedBar Bar `required:"true"`
}

type OptionalFoo struct {
	Bar
	NamedBar Bar
}

type Bar struct {
	Baz Baz
}

type BarWithBaz struct {
	Baz Baz `required:"true"`
}

type Baz struct{}

func compositeUses() {
	_ = Foo{} // want `Foo missing required keys: \[Bar NamedBar\]`
	bar := Bar{}
	_ = Foo{Bar: Bar{}, NamedBar: bar}
	_ = OptionalFoo{}
	_ = BarWithBaz{} // want `BarWithBaz missing required keys: \[Baz\]`
	_ = BarWithBaz{Baz: Baz{}}

}

type basicTypes struct {
	bool               `required:"true"`
	namedBool          bool `required:"true"`
	PublicBool         bool `required:"true"`
	namedOptionalBool  bool
	OptionalPublicBool bool
	boolP              *bool  `required:"true"`
	PublicBoolP        *bool  `required:"true"`
	DoublePointer      **bool `required:"true"`
}

func builtinUses() {
	_ = basicTypes{} // want `basicTypes missing required keys: \[DoublePointer PublicBool PublicBoolP bool boolP namedBool\]`
	var fBool bool
	fBoolP := &fBool
	_ = basicTypes{
		bool:          false,
		namedBool:     false,
		PublicBool:    false,
		boolP:         &fBool,
		PublicBoolP:   &fBool,
		DoublePointer: &fBoolP,
	}
}
