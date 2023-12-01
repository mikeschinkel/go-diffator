# go-diffator
Diffator is a Go package to provide a difference string for comparing two Go value during testing. 

Diffator does NOT output a standard format diff but is instead is optimized for a developer to recognize the difference between a value they want in their test compared with the value they got in their test, where `want==expected` and `got==actual`.

## Usage

```go
diff := diffator.Diff(value1, value2)
println(diff)
```

```go
diff := diffator.DiffWithFormat(value1,value2,"Diff: %s")
println(diff)
```

```go
d := diffator.NewDiffator()
diff := d.Diff(value1, value2)
println(diff)
```

```go
d := diffator.NewDiffator()
diff := d.DiffWithFormat(value1,value2,"Diff: %s")
println(diff)
```

## Diff Output
```go
value1 := 100
value2 := 99

// Diff: (100!=99)
```

```go
type TestStruct struct {
	Int    int
	String string
}
value1 := &TestStruct{}
value2 := &TestStruct{
  Int:    1,
  String: "hello",
}

// Diff: *TestStruct{Int:(0!=1),String:(!=hello),}
```

```go
value1 := map[string]int{"Foo": 1, "Bar": 2, "Baz": 3}
value2 := map[string]int{"Foo": 1, "Bar": 20, "Baz": 3}

// Diff: map[string]int{Bar:(2!=20),}
```

```go
value1 := map[string]int{"Foo": 1, "Bar": 2, "Baz": 3, "Superman": 0}
value2 := map[string]int{"Foo": 10, "Bar": 20, "Baz": 30, "Batman": 0}

// Diff: map[string]int{Bar:(2!=20),Baz:(3!=30),Foo:(1!=10),Superman:<missing:expected>,Batman:<missing:actual>,}
```

_Note that the above is without `diffator.Pretty := true`. To set that you must use the instance syntax by first calling `diffator.NewDiffator()`._			

## Status
In active use, but only addresses those data types that the author has needed to address his use-case.  

If you would like to use this and you find it generates a panic for an unimplemented type, [pull requests](https://github.com/mikeschinkel/go-diffator/compare) are accepted and appreciated.

## License
Apache 2.0
