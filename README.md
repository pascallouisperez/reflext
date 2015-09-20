## Reflext — Regular Expressions for Types

### tl;dr

A simple example for matching

    r := reflext.MustCompile("func({_}) error")
    if !r.Match(myFunction) {
        return fmt.Errorf("expected %s, got %T", r, myFunction)
    }

Building an expression with concrete types

    r := reflext.MustCompile("func(%T) error", Context{})

Capturing

    types, ok := r.FindAll(myFunction)

### Finding Good Matches

Let's say we want to match functions of the form

    func(Context, *struct, *struct) error

where the function arguments are an interface `Context`, followed by two `struct`s passed by reference, and returning an `error`.

Or we need to match a slice of structs (passed by value or reference)

    []struct | []*struct

Or any signed integer

    int | int8 | int16 | int32 | int64

Or any type alias of a string

    alias[string]

Here, we are introducing a type selector `alias[T]` which matches types aliases of `T` (but not `T` itself).

Another type selector is `kind[K]` to match types with kind of `K`. We saw `[]struct` above which is syntactic sugar for

    []kind[struct]

Wildcards are also supported

    map[string]*_

here for maps of `string` to any pointer to a value.

### Capturing

Matching is a good first step, yet in most cases we want to do something with the sub-types. To capture, we place sub-types between brackets such as

    func(Context, *{struct}, *{struct}) error

which would capture the type of the struct (passed by reference). Or for type aliases of `string` we would write

    {kind[string]}

### Limitations

The following are not yet implemented

* `chan` matching
* Variadic argument types (e.g `int...`)
* matching variable number of return types

## Details

### Formal Grammar

The grammar of type expressions is as follows

    E := B
       | [n]E
       | []E
       | *E
       | map[E]E
       | chan E | chan <- E | <- chan E
       | func (E, ...) R
       | kind[K]
       | alias[T]
       | _
       | %T
       | E "|" E
       | { E }

    R := E
       | (E, ...)

    B := bool | uint | int | float | complex | byte | ...

    K := B
       | struct | array | chan | func | interface | map | slice

All base types `B` e.g. `uint8`, or `float64` are supported. They are simply elided here for bervity.

### AST

The grammar transalates naturally into the following decomposition

* Exact(B)
* Slice(E)
* Ptr(E)
* Map(E, E)
* Chan(E, opt)
* Func([]E, []E)
* Kind(K)
* Alias(T)
* Any
* FirstOf([]E)
* Capture(E, index)

For captures, the index represents the location of the capturing group, starting at `0` and sequentially increasing from left to right. This handles sub-capture groups, such as

    func({{int} | {uint}}) {_}

Where

* The function's argument is group `0`
* The function's argument (if `int`) is group `1`
* The function's argument (if `uint`) is group `2`
* The function's (only) return type is group `3`