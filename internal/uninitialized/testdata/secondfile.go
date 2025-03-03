package testdata

type SecondFileStruct struct { // want SecondFileStruct:"[Bool]"
	Bool bool `required:"true"`
}

func inSecondFileReferencingFirst() {
	_ = BarWithBaz{} // want `BarWithBaz missing required keys: \[Baz\]`
}
