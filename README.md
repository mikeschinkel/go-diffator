# go-diffator
Diffator is a Go package to provide a difference string for **comparing during testing**.

Diffator does NOT output a standard format diff but is instead is optimized for a developer to recognize the difference between a value they want in their test compared with the value they got in their test, where `want==expected` and `got==actual`.

## Usage
Diffator _(currently)_ offers two (2) types of comparisons:

1. String-to-string comparison
2. Object-to-object comparison


### Usage for String-to-string comparison
```go
result := diffator.CompareStrings(string1, string2, nil)
println(result)
```

```go
result := diffator.CompareStrings(string1,string2,&diffator.StringOpts{
  MinSubstrLen: diffator.Int(3),
})
println(result)
```

```go
c := NewStringComparator(v1, v2, nil)
result := c.Compare()
println(result)
```

```go
c := NewStringComparator(v1, v2, &diffator.StringOpts{
  MinSubstrLen: diffator.Int(3),
})
result := c.Compare()
println(result)
```

To understand `diffator.Int(3)`, see [Nillable Option Values](#nillable-types).

#### Diff Output
```go
// Assuming:
string1 := "ABC"
string2 := ""

// Result: "<(ABC/)>"
```

```go
// Assuming:
string1 := ""
string2 := "ABC"

// Result: "<(/ABC)>"
```

```go
// Assuming:
string1 := "ABC"
string2 := "XYZ"

// Result: "<(ABC/XYZ)>"
```

```go
// Assuming:
string1 := "ABCDEF"
string2 := "ABCDXYZ"

// Result: "ABCD<(EF/XYZ)>"
```
```go
// Assuming:
string1 := "ABCDXYZ"
string2 := "123XYZ"

// Result: "<(ABCD/123)>XYZ"
```
```go
// Assuming:
string1 := "ABCDEF123GHI456JKLMNOP"
string2 := "ABCDEFGHIJKLMNOP"
opts := &StringOpts{
  MatchingPadLen: diffator.Int(5),
  MinSubstrLen:   diffator.Int(2),
}
// Result: "BCDEF<(123/)>GHI<(456/)>JKLMN"
```
```go
// Assuming:
string1 := "Look, it's Batman!!!"
string2 := "Look, it's Superman!!!"

// Result: "Look, it's <(Bat/Super)>man!!!"
```


### Usage for Object-to-object comparison

```go
result := diffator.CompareObjects(value1, value2, nil)
println(result)
```

```go
result := diffator.CompareObjects(value1,value2,&diffator.ObjectOpts{
  OutputFormat: diffator.String("Diff: %s"),
})
println(result)
```

```go
c := NewObjectComparator(v1, v2, nil)
result := c.Compare()
println(result)
```

```go
c := NewObjectComparator(v1, v2, &diffator.ObjectOpts{
  OutputFormat: diffator.String("Diff: %s"),
})
result := c.Compare()
println(result)
```
To understand `diffator.String("Diff: %s")`, see [Nillable Option Values](#nillable-types).


#### Diff Output
```go
// Assuming:
value1 := 100
value2 := 99

// Result: (100!=99)
```

```go
// Assuming:
type TestStruct struct {
	Int    int
	String string
}
value1 := &TestStruct{}
value2 := &TestStruct{
  Int:    1,
  String: "hello",
}

// Result: *TestStruct{Int:(0!=1),String:(!=hello),}
```

```go
// Assuming:
value1 := map[string]int{"Foo": 1, "Bar": 2, "Baz": 3}
value2 := map[string]int{"Foo": 1, "Bar": 20, "Baz": 3}

// Result: map[string]int{Bar:(2!=20),}
```

```go
// Assuming:
value1 := map[string]int{"Foo": 1, "Bar": 2, "Baz": 3, "Superman": 0}
value2 := map[string]int{"Foo": 10, "Bar": 20, "Baz": 30, "Batman": 0}

// Result: map[string]int{Bar:(2!=20),Baz:(3!=30),Foo:(1!=10),Superman:<missing:expected>,Batman:<missing:actual>,}
```

_Note that the above is without `ObjectOps.PrettyPrint := true`._			
### Nillable Option Values
We decided that in order to allow for setting of default values for `StringOpts` and `ObjectOpts` we would use values of `*diffator.IntValue`, `*diffator.BoolValue`, `*diffator.StringValue` instead of `int`, `bool`, and `string`, respectively.

To set the values, use the object constructors `diffator.Int()`,  `diffator.Bool()`, and `diffator.String()`, respectively. 

To see example usage, visit the [Usage](#usage) sections, above.

## Status
In active use, but only addresses those data types that the author has needed to address his use-case.  

If you would like to use this and you find it generates a panic for an unimplemented type, [pull requests](https://github.com/mikeschinkel/go-diffator/compare) are accepted and appreciated.

## License
Apache 2.0
