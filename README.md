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

## Status
In active use, but only addresses those data types that the author has needed to address his use-case.  If you would like to use this and you find it generated a panic for an unimplemented type, [pull requests](https://github.com/mikeschinkel/go-diffator/compare) are accepted and appreciated.

## License
Apache 2.0
