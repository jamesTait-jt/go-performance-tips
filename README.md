# go-performance

## How to Benchmark

Benchmarking in go is extremely easy due to the [testing](https://pkg.go.dev/testing) package provided in the runtime. This section will get us familiar with this tool so we can start to see how our functions perform.

### Benchmarking the Fibonacci Number Generation

There are two simple ways to calculate the n'th number in a fibonacci sequence. The first one is the recursive method, and the second is the sequential method. We'll implement both functions and use the benchmarking package to compare their performance.

Let's create a go module and add a fibonacci package with a `fibonacci.go` and `fibonacci_bench_test.go`. You should have something like this

```
├── cmd
│   ├── main.go
├── internal
│   ├── fibonacci
│   │   ├── fibonacci.go
│   │   ├── fibonacci_bench_test.go
├── go.mod
```

> :warning: Note that the `_test` suffix is required as it allows the go tool to recognise the benchmark as a test file

The first step is to implement the recursive function inside `fibonacci.go`:

```go
package fibonacci

func fib(n int) int {
    if n <= 1 {
        return n
    }

    return fib(n-1) + fib(n-2)
}
```

Now to benchmark it! We will use a table test so we can see how the function performs as the input `n` changes.

```go
package fibonacci

import "testing"

func Benchmark_Fib(b *testing.B) {

    scenarios := []struct{
        name string
        input int
    }{
        {   
            name: "base_0",
            input: 0,
        }, {
            name: "base_1",
            input: 1,
        }, {
            name: "recursive",
            input: 10,
        }, {
            name: "recursive_big",
            input: 25,
        },  
    }

    for _, bs := range scenarios {
        b.Run(bs.name, func(b *testing.B) {
            for i := 0 ; i < b.N ; i++ {
                fib(bs.input)
            }
        })
    }
}
```

Now we can implement the sequential method

```go
func fibSequential(n int) int {
    if n <= 1 {
        return n
    }

    var n2, n1 int = 0, 1

    for i := 2 ; i < n ; i++ {
        n2, n1 = n1, n1 + n2
    }
 
    return n2 + n1
}
```

and add some benchmarks

```go
func Benchmark_FibSequential(b *testing.B) {

    scenarios := []struct{
        name string
        input int
    }{
        {
            name: "base_0",
            input: 0,
        }, {
            name: "base_1",
            input: 1,
        }, {
            name: "recursive",
          input: 10,
        }, {
            name: "recursive_big",
            input: 25,
        },
    }

    for _, bs := range scenarios {
        b.Run(bs.name, func(b *testing.B) {
            for i := 0 ; i < b.N ; i++ {
                fibSequential(bs.input)
            }
        })
    }
}
```

Now we're set up ready to run both benchmarking functions to see how the functions stack up. Navigate to `internal/fibonacci` in the terminal and run the following command

```sh
go test -bench .
```

And you should see something like this

```sh
Benchmark_Fib/base_0-10         543429344                2.034 ns/op
Benchmark_Fib/base_1-10         591273541                2.028 ns/op
Benchmark_Fib/recursive-10       6637140               180.6 ns/op
Benchmark_Fib/recursive_big-10              4675            250465 ns/op
Benchmark_FibSequential/base_0-10       1000000000               0.3133 ns/op
Benchmark_FibSequential/base_1-10       1000000000               0.3119 ns/op
Benchmark_FibSequential/recursive-10    321476074                3.742 ns/op
Benchmark_FibSequential/recursive_big-10                142730412                8.417 ns/op
```

The first column is obviously the name we gave the benchmark scenario, the second column is the number of times the function was ran by the benchmarking function, and the third column is the number of nanoseconds per function call.

This third column is the one we are interested in.

As we can see, the recursive function is outperformed by the sequential version by quite some margin.

### Benchmarking heap allocations in string parsing

Now that we are able to benchmark the runtime of a function, how about benchmarking the number of times a function allocates to the heap? When writing Go at any kind of scale, garbage collection often becomes a bottleneck. This means, if we are trying to write high performance code at a large scale, we must be vigilant in trying to avoid needless heap allocations.

We often receive delimited strings and need to parse them into some sort of structure. We won't go into too much detail yet, but it serves as a nice example for how we can benchmark heap allocations.

So, let's create another package in our module called `parse` where we will add a `parse.go` and `parse_bench_test.go`

```go
package parse

import "strings"

type person struct {
    firstName string
    lastName string
}

func Person(s string) person {
    split := strings.Split(s, " ")

    n := person{
        firstName: split[0],
        lastName: split[1],
    }

    return n
}
```

Here we are simply taking a string, splitting it on a delimiter and populating two fields on a struct. A very common pattern, and one where we might not even think twice about the fact that we might be haemorrhaging memory to the heap unnecessarily. So, let's write the benchmarks to see if we have a problem!

```go
package parse

import "testing"

func Benchmark_Person(b *testing.B) {

    in := "joe bloggs"

    b.ReportAllocs()
    for i := 0 ; i < b.N ; i++ {
        Person(in)
    }
}
```

This time we have used the `b.ReportAllocs()` method which will let the benchmarker know that we want to see how much memory we're allocating to the heap.

Let's run the benchmarks using the same command as before and we should get something like this

```sh
Benchmark_Person-10       33384687                35.78 ns/op           32 B/op          1 allocs/op
```

So now we have two new columns. We can wee the number of bytes allocated to the heap in each function call, as well as the number of times we allocated to the heap. So in this example, we have one heap allocation of 32B.

This is useful information, but what do we do with it? How do we know which line allocated the memory?

### Using pprof

[pprof](https://pkg.go.dev/runtime/pprof) is a very handy tool provided by the go runtime that allows us to profile our applications in order to find cpu/memory hotspots. This is often the first step in fine tuning our applications as we must know which areas to optimise before we start benchmarking.

There are some great blog posts about this online, easily found with some googling, so I will just introduce the basics that we need to get started!

Navigate to the `parse` package and run the following command from the terminal

```sh
go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out
```

which should generate two files for you; `cpu.out` and `mem.out`. These are not human readable, but we have another way of reading them using `pprof`

```sh
go tool pprof mem.out
```

This will take you into an interactive tool that will help us to understand some of the stats. Type: alloc_space signifies that the data we have collected is regarding all memory allocated to the heap in the lifecycle of the benchmarks.

```sh
Type: alloc_space
Time: May 11, 2023 at 5:22pm (BST)
Entering interactive mode (type "help" for commands, "o" for options)
```

 Since we're going for the basics here, we can type `top 10` to see the areas we allocate the most memory from, and we should see something like this:

```sh
(pprof) top 10
Showing nodes accounting for 1724.11MB, 99.79% of 1727.75MB total
Dropped 26 nodes (cum <= 8.64MB)
      flat  flat%   sum%        cum   cum%
 1203.59MB 69.66% 69.66%  1724.11MB 99.79%  github.com/jamesTait-jt/go-performance/internal/parse.Person
  520.52MB 30.13% 99.79%   520.52MB 30.13%  strings.genSplit
         0     0% 99.79%  1724.11MB 99.79%  github.com/jamesTait-jt/go-performance/internal/parse.Benchmark_Person
         0     0% 99.79%   520.52MB 30.13%  strings.Split (inline)
         0     0% 99.79%  1724.11MB 99.79%  testing.(*B).launch
         0     0% 99.79%  1724.11MB 99.79%  testing.(*B).runN
```

Here we can see `flat` and `cum`

- `flat` means the memory was allocated by this function call
- `cum` means the memory was allocated by this function call or a function it called down the stack

So we can see that the `strings.genSplit` function seems to allocate a lot of memory which highlights that it might be an area we want to look at when we come to making optimisations.

The exact same method can be used to analyse the cpu profile to find areas that the application spends a lot of time in.

To exit the interactive pprof tool, type `q` and press enter!

##  Heap vs. Stack Allocations

The Go runtime takes care of memory management for us, meaning we are able to focus on correctness and maintainability while writing applications. However, as mentioned before, for large programs, the garbage collection step can quickly become a bottleneck. Before we dive into some common memory issues, we should discuss the two main areas of memory. The stack and the heap.

The stack refers to the callstack of a thread. It is a LIFO data structure that stores data as a thread executes functions. Each function call pushes a new frame to the stack, and each returning function pops (removes) from the stack.

This means that when a function returns, we must be able to safely free the memory of the most recent stack frame, which means we must be sure that nothing is stored on the stack that must be used later.

The amount of memory available to a thread's stack is limited, and normally quite small (this will be important later). Storing memory on the stack is very quick, and it is simple to free it as the memory is simply cleared.

The heap is an area of memory that has no relation to any specific thread and it can be accessed from anywhere within the program. This is where data is stored if it must be accessed after a function exits. Storing data on the heap is slow, and freeing the memory is not straightforward as we need to be sure the data is no longer needed. This happens during garbage collection and is a very expensive operation.

Go manages all of this for us by determining whether arguments and variables need to be stored on the stack or on the heap using something called [escape analysis](https://tip.golang.org/src/cmd/compile/internal/escape/escape.go).

I've barely touched the surface of this topic, and if you want to learn more, you can check out [this](https://medium.com/@ankur_anand/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed) blog post. The main thing we need to take away is that storing data on the stack is preferable to the heap as we do not need to garbage collect it.

## A Better Way to Split

As we saw in the first section, splitting a string into a slice will have to allocate to the heap. This is because the compiler cannot determine how big this slice will be, so it must go on the heap to avoid a potential stack overflow.

Since the slice is only an intermediate for us, i.e. we want to build the `person` object and do not need to keep the slice, we can avoid this heap allocation if we handle the string as a stream.

The approach will be to do the following:

   1. parse the string until the first delimiter
   2. populate the `person.first` field with what we've parsed so far
   3. parse until the next delimiter
   4. populate the `person.last` field with what we parsed since the last delimiter

This way, we can avoid using a slice and save on some memory!

First let's write a function to parse until a delimiter.

```go
func ParseUntil(s string, sep rune) (string, string) {
    if len(s) == 0 {
        return "", ""
    }

    indexOfNext := strings.Index(s, " ")
    if indexOfNext == -1 {
        return s, ""
    }

    return s[:indexOfNext], s[indexOfNext + 1:]
}
```

Now we can use this in our person parsing and hopefully we'll see an improvement in performance

```go
func PersonEfficient(s string) person {
    first, last := ParseUntil(s, ' ')

    n := person{
        firstName: first,
        lastName: last,
    }

    return n
}
```

As you can see, we're removed the slice, so now, nothing should escape to the heap. Let's benchmark it to be sure

```go
func Benchmark_PersonEfficient(b *testing.B) {
    s := "joe bloggs"

    b.ReportAllocs()
    for i := 0 ; i < b.N ; i++ {
        _ = PersonEfficient(s)
    }
}
```

Let's run both benchmarks together so we can compare

```sh
Benchmark_Person
Benchmark_Person-10               32530510                36.55 ns/op           32 B/op          1 allocs/op
Benchmark_PersonEfficient
Benchmark_PersonEfficient-10      261291204                4.577 ns/op           0 B/op          0 allocs/op
```

As we can see, the allocation has gone and the function is also nearly 10x faster! So, not only have we saved on the garbage collection time, but by avoiding writing to the heap we have saved ourselves an expensive operation.

This method is very useful when we want to parse a message from another service and convert it into a structure that is known by our service. I have shown it with strings, but the same method can be used with any slice, like a slice of bytes, for example.

## string <-> []byte conversions

> :warning: This section contains some unsafe operations. Be very careful using any code written here without understanding exactly what you are doing

Parsing between `string` and `[]byte` types is very common in go. For example, if we are parsing a payload sent over by tcp, we will receive it as a `[]byte` type. But if we want this to be human readable, then we will need to convert it into `string`.

Using our `person` object as an example, let's suppose we receive a `[]byte` where the byte slice represents the person. i.e. each byte represents part of the string representation of the person, or a delimiter. For example, the output of `[]byte("joe bloggs")` where `" "` is our delimiter.

Something like this is likely to be on the boundaries of your service, and thus may be extremely high in traffic, so any optimisations here are likely to have a decent impact on overall performance of your application.

Now that the motivation is out of the way, let's see why these conversions are pretty expensive.

The first thing we need to understand is that strings in go are immutable. This means that once they are defined, they cannot be changed.

The second thing is that strings and byte slices are both backed by [arrays](https://go.dev/tour/moretypes/6) under the hood.

So, let's look at some code.

```go
func PersonBytes(bs []byte) person {
    first, last := ParseUntilBytes(bs, '\x1f')

    p := person{
        firstName: string(first),
        lastName: string(last),
    }

     return p
}
```

We are using a `[]byte` version of the `ParseUntil` function which we proved to be fast and memory efficient. And after that, we simply convert each field to a string so we can read the names. Now we can benchmark it as see how it performs.

```go
func Benchmark_PersonBytes(b *testing.B) {
    bs := []byte("joe\x1fbloggs")

    b.ReportAllocs()
    for i := 0 ; i < b.N ; i++ {
        PersonBytes(bs)
    }
}
```

```sh
Benchmark_PersonBytes
Benchmark_PersonBytes-10        51941656                22.68 ns/op           10 B/op          2 allocs/op
```

We can see there are two allocations, and if we use escape analysis, we can see that both the string conversions escape to the heap:

```sh
$ go build -gcflags="-m"
./parse.go:39:20: string(first) escapes to heap
./parse.go:40:19: string(last) escapes to heap
```

So why is this? Well, since strings are immutable, we can't just point the string at the backing array of the byte slice. If the byte slice were to be changed, the changes would reflect in the string, breaking go's promise of immutability. But since the string escapes the function, we must put it on the heap so that we don't delete it from memory once the stackframe has been popped. This means that every time we convert to and from byte slices, if the string is used afterwards, we will be performing costly heap allocations.

> :warning: Danger starts here!!!!

We can avoid these allocations by breaking the string immutability promise. This is dangerous so we must be certain that we are careful here. In particular, we need to ensure that the byte slice will come out of scope as soon as we have converted it to a string. This will make sure we avoid any changes to the array that is now backing both the string and the byte slice.

To begin, let's create a new package called `unsafe` and add the following file:

```go
// unsafe.go

package unsafe

import "unsafe"

func ToString(bytes []byte) string {
    return *(*string)(unsafe.Pointer(&bytes))
}
```

What this function essentially does is create a string that points at the same backing array as the byte slice that was passed in. Let's see how this can cause issues.

```go
func example() {
    bs := []byte{'1', '1', '1'}
    s := ToString(bs)

    fmt.Println(s)

    bs[0] = '0'

    fmt.Println(s)
}
```

And here is the output

```go
111
011
```

Now that we understand the issues that can be caused by this method, let's see what happens when we use it in our person parsing

```go
func PersonBytesUnsafe(bs []byte) person {
    first, last := ParseUntilBytes(bs, '\x1f')

    p := person{
        firstName: unsafe.ToString(first),
        lastName: unsafe.ToString(last),
    }

    return p
}
```

Now write a benchmark

```go
func Benchmark_PersonBytesUnsafe(b *testing.B) {
    bs := []byte("joe bloggs")

    b.ReportAllocs()
    for i := 0 ; i < b.N ; i++ {
        _ = PersonBytesUnsafe(bs)
    }
}
```

```sh
Benchmark_PersonBytesUnsafe
Benchmark_PersonBytesUnsafe-10          285726729                4.164 ns/op           0 B/op          0 allocs/op
```

So we have removed the allocations, but at the cost of making or application slightly unsafe. I wouldn't really recommend using this method if you intend on keeping the string, as the risk is too high. However, the solution really comes into its own when the string is an intermediate type when converting to something else, like an integer or a float.

### []byte -> string -> other

Let's extend our `person` struct to contain `age` and `height` fields with `int` and `float64` types respectively.

```go
type person struct {
    firstName string
    lastName string
    age int
    height float64
}
```

And now when we receive a new person represented as a byte slice, we will also receive the two new fields in their string representation. e.g. the output of `[]byte("joe\x1fbloggs\x1f25\x1f182.3")`. So, let's modify our parsing function to see how it performs with the new requirements.

```go
func PersonBytes(bs []byte) (person, error) {
    first, rest := ParseUntilBytes(bs, '\x1f')
    last, rest := ParseUntilBytes(rest, '\x1f')
    ageBytes, rest := ParseUntilBytes(rest, '\x1f')
    heightBytes, _ := ParseUntilBytes(rest, '\x1f')

    ageStr := string(ageBytes)
    age, err := strconv.Atoi(ageStr)
    if err != nil {
        return person{}, err
    }

    heightStr := string(heightBytes)
    height, err := strconv.ParseFloat(heightStr, 64)
    if err != nil {
        return person{}, err
    }

    p := person{
        firstName: string(first),
        lastName: string(last),
        age: age,
        height: height,
    }

    return p, nil
}
```

And updating our benchmarking function we get

```go
func Benchmark_PersonBytes(b *testing.B) {
    bs := []byte("joe\x1fbloggs\x1f25\x1f182.3")

    b.ReportAllocs()
    for i := 0 ; i < b.N ; i++ {
        PersonBytes(bs)
    }
}
```

```sh
Benchmark_PersonBytes
Benchmark_PersonBytes-10        12817989                91.76 ns/op           16 B/op          4 allocs/op
```

So we are allocating even more memory to the heap due to the string conversions of the `age` and `height` fields. Now, let's introduce some unsafe functions to help us out (these ones are actually pretty safe!)

```go
func ToInt(bytes []byte) (int, bool) {
    s := ToString(bytes)

    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, false
    }

    return n, true
}

func ToFloat64(bytes []byte) (float64, bool) {
    s := ToString(bytes)

    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        return 0, false
    }

    return f, true
}
```

Using our `ToString` function, we should be able to avoid the allocations to the heap. However, in this case we do not break the promise of string immutability as the string is contained within the `ToInt` and `ToFloat64` functions, meaning it is popped from the stack when this function returns anyway. This leaves us with an integer/float64 conversion, from a byte slice without any heap allocations! Let's prove it!

```go
func Benchmark_PersonBytesUnsafe(b *testing.B) {
    bs := []byte("joe\x1fbloggs\x1f25\x1f182.3")

    b.ReportAllocs()
    for i := 0 ; i < b.N ; i++ {
        PersonBytesUnsafe(bs)
    }
}
```

```sh
Benchmark_PersonBytesUnsafe
Benchmark_PersonBytesUnsafe-10          19736895                59.81 ns/op            0 B/op          0 allocs/op
```

### 32 byte optimisation

There is an interesting compiler optimisation when converting `[]byte` to `string` when the slice has less than 32 bytes. There is a hardcoded buffer to avoid the heap allocations when the string does not escape.

```go
func BenchmarkString(b *testing.B) {
    b.Run("LessThan32", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0 ; i < b.N ; i++ {
            in := make([]byte, 31)
            _ = string(in)
        }
    })

    b.Run("EqualTo32", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0 ; i < b.N ; i++ {
            in := make([]byte, 32)
            _ = string(in)
        }
    })

    b.Run("GreaterThan32", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0 ; i < b.N ; i++ {
            in := make([]byte, 33)
            _ = string(in)
        }
    })
}
```

```sh
BenchmarkString
BenchmarkString/LessThan32
BenchmarkString/LessThan32-10           383141762                3.009 ns/op           0 B/op          0 allocs/op
BenchmarkString/EqualTo32
BenchmarkString/EqualTo32-10            378646593                3.005 ns/op           0 B/op          0 allocs/op
BenchmarkString/GreaterThan32
BenchmarkString/GreaterThan32-10        72316562                15.73 ns/op           48 B/op          1 allocs/op
```

However, I have not been able to get this to work when parsing the resultant string to an integer/float as we saw in the above examples. I have looked into this a bit to figure out where the string might escape to the heap, but in the words of a very accomplished architect: "It's easier to make a stone bleed than understanding the Go compiler". If anyone knows the reason, I would love to know it!

## Pre-allocating Slices

As we have discussed, allocating memory at runtime is expensive. So if we know at compile time the upper bound of a slice's size, it is worth pre-allocating the capacity with a constant. Here is an example of a growing slice without pre allocating memory

```go
func growSlice() {
    toAdd := []byte("123456")
    bs := []byte{}

    bs = append(bs, toAdd...)
    bs = append(bs, toAdd...)
    bs = append(bs, toAdd...)
}
```

Looking at the benchmarks it is easy to see that a lot of heap allocations are taking place.

```sh
Benchmark_GrowSlice
Benchmark_GrowSlice-10                  23496968                48.22 ns/op           56 B/op          3 allocs/op
```

every time the slice grows past its capacity, the entire slice needs to be reallocated on the heap with a larger capacity. The method is useful in a large number of use cases where we do not know the capacity of the slice, however, if we do know it then we can do a lot better and save on those costly runtime allocations.

```go
func growSlicePreAllocate() {
    toAdd := []byte("123456")
    bs := make([]byte, 0, 18)

    bs = append(bs, toAdd...)
    bs = append(bs, toAdd...)
    bs = append(bs, toAdd...)
}
```

Here, we know that the `toAdd` bytes will be appended to the `bs` slice 3 times, giving us a final slice of length `6 * 3 = 18` bytes. So we pre-allocate this much memory which allows the compiler to allocate memory on the stack rather than allocating on the heap at runtime, and then having to handle the garbage collection as well.

```sh
Benchmark_GrowSlicePreAllocate
Benchmark_GrowSlicePreAllocate-10       659889508                1.797 ns/op           0 B/op          0 allocs/op
```

## Interfaces

### Channels and function parameters

Sometimes it is preferable to use channels that accept the blank interface `interface{}` type, or functions that accept or return an `interface{}` type. However this can come with some performance issues. Whenever we do this, the compiler has to cast the type of the object passed into an `interface{}` type. If the object cannot fit in a single machine word (int, bool, etc) then a pointer to the object value object will be allocated to the heap, and the object value may be moved to the heap (this is not guaranteed behaviour). The same happens with typed interfaces.

This is often not a problem, but if you have noticed that the garbage collector is causing you issues then it might be worth redesigning some areas to avoid interfaces and use concrete types instead.

Because of this behaviour, logging tends to be very expensive as the functions accept interface types to allow them to be flexible. I would recommend using a third party logger such as [zerolog](https://github.com/rs/zerolog).

### Interface method calls

Sometimes when using interface methods, we can unknowingly allocate some of the method arguments to the heap. The situation in which this occurs is when the arguments to the function are pointers (remember this includes slices and strings!).

The reason is that the compiler has no knowledge about the implementation of an interface method at compile time. So it is possible for the implementation to store the value somewhere (maybe a package level variable) and the compiler will have no idea.

Let's say this is the case, and the method argument was a pointer. So the package level variable is pointing to the same location as the original pointer that we passes into the method. However, When the original pointer goes out of scope, the value pointed to by the pointer will be trashed. And thus, the package level variable will be pointing to garbage. Because of this, the pointer value is moved to the heap so that it is not trashed when the stackframe is popped.

Here is an example with some benchmarks

```go
package interfaces

type NoOpper interface {
    Int(n int)
    IntPtr(n *int)
}

type nothing struct {}

func (no nothing) Int(n int) {}

func (no nothing) IntPtr(n *int) {}

//go:noinline
func NoOpInt(no NoOpper, n int) {
    no.Int(n)
}

//go:noinline
func NoOpIntPtr(no NoOpper, n int) {
    no.IntPtr(&n)
}
```

Here we have two trivial interface methods, one accepts an `int` and the other accepts an `*int`. We use `//go:noinline` to prevent the compiler from optimising out the function call as we can't rely on this in a production setting with ever-changing code. 

Even though these methods don't do anything, the expectation is that with the `*int` version, `n` will be allocated to the heap.

```go
package interfaces

import (
    "testing"
)

func BenchmarkNoOpInt(b *testing.B) {
    b.ReportAllocs()

    no := nothing{}
    n := 10

    for i := 0 ; i < b.N ; i++ {
        NoOpInt(no, n)
    }
}

func BenchmarkNoOpIntPtr(b *testing.B) {
    b.ReportAllocs()

    no := nothing{}
    n := 10

    for i := 0 ; i < b.N ; i++ {
        NoOpIntPtr(no, n)
    }
}
```

```sh
BenchmarkNoOpInt
BenchmarkNoOpInt-10             578409886                2.068 ns/op           0 B/op          0 allocs/op
BenchmarkNoOpIntPtr
BenchmarkNoOpIntPtr-10          131864689                9.051 ns/op           8 B/op          1 allocs/op
```

As expected, the variable is allocated to the heap. Of course, the tradeoff is that if we did not pass a pointer, we would have to take a copy of the parameter to pass it into the called function's stackframe. The recommendation here is to benchmark both options and see which is better for your application.

## for i vs for range

There are a couple of ways to iterate through slices in go. The two most common methods are what I will refer to as `for i` and `for range` 

```go
// for i
for i := 0 ; i < len(sl) ; i++ {
    // do something
}
```

```go
// for range
for _, elem := range sl {
    // do something
}
```

for range is commonly used as it gives you a variable `elem` to use rather than having to index the slice. However, this comes at a cost when the elements in the slice are large due to the fact that go will copy the element into the range variable.

```go
type bigObj struct {
    id int
    lotsOfStuff [10 * 1024]byte
}

type otherObj struct {
    id int
}

func forI(sl []bigObj) otherObj {
    oo := otherObj{}
    for i := 0 ; i < len(sl) ; i++ {
        if oo.id == 5 {
            oo.id = sl[i].id
        }
    }
    return oo 
}

func forRange(sl[]bigObj) otherObj {
    oo := otherObj{}
    for _, elem := range sl {
        if oo.id == 5 {
            oo.id = elem.id
        }
    }
    return oo
}
```

These are just simple functions to iterate through a list of big objects and then copy a field from one object to another given a certain condition. If the `if` clause is not included, the performance cost is not realised in the benchmarks. I'm not sure why but I'm assuming some kind of compiler optimisation. Again, if anyone has the answer I would like to hear it!

Here are the benchmarks

```go
func BenchmarkForI(b *testing.B) {
    sl := []bigObj{
        {id: 0},
        {id: 1},
        {id: 2},
        {id: 3},
        {id: 4},
        {id: 5},
    }

    for i := 0 ; i < b.N ; i++ {
        forI(sl)
    }
}

func BenchmarkFoRange(b *testing.B) {
    sl := []bigObj{
        {id: 0},
        {id: 1},
        {id: 2},
        {id: 3},
        {id: 4},
        {id: 5},
    }

    for i := 0 ; i < b.N ; i++ {
        forRange(sl)
    }
}
```

and the results:

```sh
BenchmarkForI
BenchmarkForI-10        231138032                5.171 ns/op           0 B/op          0 allocs/op
BenchmarkFoRange
BenchmarkFoRange-10       467359              2543 ns/op               0 B/op          0 allocs/op
```

As we can se, the forRange version performs far worse in this scenario! As a consequence of this, I would recommend sticking with the forI syntax.
