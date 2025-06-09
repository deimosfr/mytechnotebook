---
title: "Rust: learn basics"
slug: rust-basics/
description: "Learn the basics of Rust programming language."
categories: ["Rust"]
tags: ["Rust", "Programming", "Language"]
---

Rust is a systems programming language that emphasizes safety, concurrency, and performance. It achieves memory safety without garbage collection through its unique ownership system, while still maintaining high performance comparable to C/C++. This guide focuses on practical knowledge to get you productive quickly.

## Core Philosophy

Rust's design principles make it unique among programming languages:

- **Memory Safety**: Prevents common programming errors like null pointer dereferencing, dangling pointers, and data races. This is achieved through the compiler's strict checking rather than runtime checks.
- **Zero-Cost Abstractions**: High-level abstractions compile to efficient machine code. You can write expressive code without performance penalties.
- **Fearless Concurrency**: Safe concurrent programming without data races. The compiler prevents common concurrency bugs at compile time.
- **No Runtime**: No garbage collector or runtime, making it suitable for embedded systems and performance-critical applications.

## Basic Syntax and Data Types

### Variables and Mutability

Rust's variable system is designed to prevent bugs and make code more predictable:

```rust
// Variables are immutable by default - this is a key safety feature
// Once a value is bound to a name, it cannot be changed
let x = 5;
// x = 6;  // This would cause a compile error

// To make a variable mutable, use the 'mut' keyword
// This explicitly shows your intent to change the value
let mut y = 5;
y = 6; // This is allowed because of 'mut'

// Constants are always immutable and must have type annotation
// They can be declared in any scope, including global
const MAX_POINTS: u32 = 100_000;

// Shadowing allows you to reuse a variable name
// This creates a new variable, effectively hiding the previous one
// Useful when you want to transform a value but keep it immutable
let x = 5;
let x = x + 1;  // x is now 6, but still immutable
let x = x * 2;  // x is now 12, still immutable
```

### Basic Types

Rust is a statically typed language, which means the compiler must know the type of every variable. However, it can often infer types, so you don't always need to write them:

```rust
// Scalar Types - represent a single value
let integer: i32 = 42;      // Signed integer (can be negative)
let unsigned: u32 = 42;     // Unsigned integer (only positive)
let float: f64 = 3.14;      // Floating point (64-bit precision)
let boolean: bool = true;   // Boolean (true or false)
let character: char = 'A';  // Character (Unicode scalar value)

// Compound Types - group multiple values
// Tuples: fixed-length collection of values of different types
let tuple: (i32, f64, char) = (1, 2.0, 'a');
let (x, y, z) = tuple;  // Destructuring a tuple

// Arrays: fixed-length collection of values of the same type
// Size must be known at compile time
let array: [i32; 5] = [1, 2, 3, 4, 5];  // Type is [i32; 5]
let first = array[0];  // Accessing elements (zero-based)

// Slices: reference to a portion of an array or string
// They don't have ownership, just a view into the data
let slice: &[i32] = &array[1..3];  // References elements 1 and 2
```

## Ownership and Borrowing

Ownership is Rust's most unique feature and the key to its memory safety guarantees. It's a set of rules that the compiler checks at compile time.

### Ownership Rules

1. Each value has an owner (a variable that owns it)
2. Only one owner at a time
3. When owner goes out of scope, value is dropped (memory is freed)

```rust
// Ownership transfer (move semantics)
let s1 = String::from("hello");  // s1 owns the string
let s2 = s1;  // Ownership moves from s1 to s2
// println!("{}", s1);  // This would fail - s1 no longer owns the string
// This prevents double-free errors that can occur in other languages

// Clone for deep copy
let s1 = String::from("hello");
let s2 = s1.clone();  // Creates a complete copy of the data
println!("{} {}", s1, s2);  // Both valid because they own different data
// Note: clone can be expensive for large data

// Stack-only data (like integers) implements Copy trait
let x = 5;
let y = x;  // x is still valid because integers are copied, not moved
println!("{} {}", x, y);  // Both valid
```

### Borrowing

Borrowing allows you to access data without taking ownership. This is crucial for performance and flexibility:

```rust
// Immutable borrows (multiple allowed)
let s = String::from("hello");
let r1 = &s;  // r1 borrows s
let r2 = &s;  // r2 also borrows s
println!("{} {}", r1, r2);  // Both references are valid
// The data can't be modified through these references

// Mutable borrows (only one allowed)
let mut s = String::from("hello");
let r1 = &mut s;  // r1 mutably borrows s
// let r2 = &mut s;  // This would fail - can't have multiple mutable borrows
// let r3 = &s;      // This would fail - can't have immutable borrow while mutable exists
// These rules prevent data races at compile time

// Borrowing rules in practice
fn main() {
    let mut s = String::from("hello");

    {
        let r1 = &s;  // First borrow
        let r2 = &s;  // Second borrow
        println!("{} {}", r1, r2);
    } // r1 and r2 go out of scope here

    let r3 = &mut s;  // Now we can mutably borrow
    r3.push_str(" world");
}
```

## Pattern Matching and Control Flow

Rust's pattern matching is powerful and exhaustive, ensuring you handle all possible cases:

### Match Expressions

```rust
// Match is like a switch statement but more powerful
// It must be exhaustive (cover all possible cases)
let number = 13;

match number {
    1 => println!("One!"),  // Match exact value
    2 | 3 | 5 | 7 | 11 => println!("Prime!"),  // Match multiple values
    13..=19 => println!("Teen!"),  // Match range
    _ => println!("Something else"),  // Catch-all pattern
}

// Pattern matching with destructuring
// This is useful for extracting values from complex types
let point = (3, 4);
match point {
    (0, 0) => println!("Origin"),  // Match exact tuple
    (0, y) => println!("Y-axis at {}", y),  // Extract y value
    (x, 0) => println!("X-axis at {}", x),  // Extract x value
    (x, y) => println!("Point at ({}, {})", x, y),  // Extract both values
}

// Match with guards (additional conditions)
match number {
    n if n < 0 => println!("Negative"),
    n if n > 0 => println!("Positive"),
    _ => println!("Zero"),
}
```

### If Let and While Let

These are convenient syntax for when you only care about one pattern:

```rust
// if let for single pattern matching
// More concise than match when you only care about one case
let some_value = Some(3);
if let Some(x) = some_value {
    println!("Got: {}", x);
} else {
    println!("Got nothing");
}

// while let for pattern matching in loops
// Useful for iterating until a pattern doesn't match
let mut stack = Vec::new();
stack.push(1);
stack.push(2);
while let Some(top) = stack.pop() {  // Continues while pop returns Some
    println!("{}", top);
}
```

## Error Handling

Rust uses the type system for error handling, making it explicit and forcing you to handle errors:

### Result and Option

```rust
// Option<T> represents a value that might not exist
// This replaces null/undefined from other languages
fn find_user(id: u32) -> Option<String> {
    if id == 1 {
        Some(String::from("John"))  // Value exists
    } else {
        None  // No value
    }
}

// Result<T, E> represents an operation that might fail
// T is the success type, E is the error type
fn divide(a: f64, b: f64) -> Result<f64, String> {
    if b == 0.0 {
        Err(String::from("Division by zero"))  // Operation failed
    } else {
        Ok(a / b)  // Operation succeeded
    }
}

// Using ? operator for error propagation
// This is a shorthand for handling errors
fn process_division() -> Result<f64, String> {
    let result = divide(10.0, 2.0)?;  // Returns early if Err
    Ok(result * 2.0)
}

// Common pattern for handling Option/Result
match find_user(1) {
    Some(name) => println!("Found user: {}", name),
    None => println!("User not found"),
}

// Using if let for cleaner code when you only care about success
if let Some(name) = find_user(1) {
    println!("Found user: {}", name);
}
```

## Common Patterns and Best Practices

### Structs and Enums

```rust
// Structs are custom data types that group related data
struct User {
    name: String,    // Field with type
    age: u32,        // Another field
    active: bool,    // And another
}

// Tuple structs are useful for simple groupings
struct Point(i32, i32);  // Similar to a tuple but with a name

// Enums can hold different types of data
// This is more powerful than enums in other languages
enum Message {
    Quit,                          // No data
    Move { x: i32, y: i32 },      // Struct-like
    Write(String),                // Tuple-like
    ChangeColor(i32, i32, i32),   // Tuple-like with multiple values
}

// Implementing methods for types
impl User {
    // Constructor (convention in Rust)
    fn new(name: String, age: u32) -> User {
        User {
            name,
            age,
            active: true,
        }
    }

    // Method that takes &self (reference to the instance)
    fn is_adult(&self) -> bool {
        self.age >= 18
    }

    // Method that takes &mut self (mutable reference)
    fn deactivate(&mut self) {
        self.active = false;
    }
}
```

### Traits

Traits are similar to interfaces in other languages but more powerful:

```rust
// Defining a trait (interface)
trait Drawable {
    fn draw(&self);  // Method that implementors must define
    fn area(&self) -> f64;  // Another required method
}

// Implementing a trait for a type
struct Circle {
    radius: f64,
}

impl Drawable for Circle {
    fn draw(&self) {
        println!("Drawing circle with radius {}", self.radius);
    }

    fn area(&self) -> f64 {
        std::f64::consts::PI * self.radius * self.radius
    }
}

// Using trait bounds to constrain generic types
fn draw_shape<T: Drawable>(shape: T) {
    shape.draw();
    println!("Area: {}", shape.area());
}

// Default implementations in traits
trait Printable {
    fn print(&self) {
        println!("Default implementation");
    }
}
```

### Common Collections

```rust
// Vector: growable array
let mut vec = Vec::new();  // Create empty vector
vec.push(1);              // Add elements
vec.push(2);
println!("First: {}", vec[0]);  // Access by index
for x in &vec {           // Iterate
    println!("{}", x);
}

// HashMap: key-value store
use std::collections::HashMap;
let mut map = HashMap::new();
map.insert(String::from("key"), 1);  // Insert
if let Some(value) = map.get("key") {  // Get with pattern matching
    println!("Value: {}", value);
}

// String: growable UTF-8 text
let mut s = String::from("hello");
s.push_str(" world");     // Append
s.push('!');             // Append single character
println!("{}", s);       // Print: "hello world!"
```

## Best Practices

1. **Use `Option` and `Result`** instead of null or exceptions

   - Makes error handling explicit
   - Forces you to handle all cases
   - Prevents null pointer dereferencing

2. **Prefer references** over cloning when possible

   - More efficient
   - Clearer ownership semantics
   - Use `clone()` only when necessary

3. **Use the type system** to enforce invariants

   - Create new types for different concepts
   - Use enums to represent state
   - Leverage the compiler to catch errors

4. **Leverage pattern matching** for control flow

   - Exhaustive checking
   - Clear intent
   - Destructuring for clean code

5. **Use iterators** instead of manual loops when possible

   - More expressive
   - Less error-prone
   - Often more efficient

6. **Document public APIs** with `///` comments

   - Generates documentation
   - Shows examples
   - Documents behavior

7. **Use `cargo clippy`** for additional linting

   - Catches common mistakes
   - Suggests improvements
   - Enforces style

8. **Write tests** using `#[test]` attributes
   - Unit tests
   - Integration tests
   - Documentation tests

## Common Gotchas

1. **Move semantics** can be surprising when working with non-Copy types

   - String, Vec, and other heap-allocated types are moved
   - Primitives and simple types are copied
   - Use references to avoid moves

2. **Lifetime annotations** might be needed for complex references

   - Required when returning references
   - Needed for structs containing references
   - Helps compiler verify reference validity

3. **Trait bounds** can become complex with multiple constraints

   - Use where clauses for clarity
   - Consider creating new traits
   - Use trait objects when appropriate

4. **String vs &str** - know when to use each

   - String: owned, growable
   - &str: borrowed, fixed-size
   - Use &str for function parameters when possible

5. **Clone vs Copy** - understand the difference
   - Copy: implicit copying (primitives)
   - Clone: explicit deep copy
   - Implement Copy only for small types

## Concurrency in Rust

Rust provides powerful tools for concurrent programming, with both traditional threading and modern async/await patterns.

### Multithreading

Rust's standard library provides thread-safe primitives and a threading model that prevents data races at compile time:

```rust
use std::thread;
use std::time::Duration;
use std::sync::{Arc, Mutex, mpsc};

// Basic thread creation
fn basic_threading() {
    // Spawn a new thread
    let handle = thread::spawn(|| {
        for i in 1..=5 {
            println!("Thread: count {}", i);
            thread::sleep(Duration::from_millis(100));
        }
    });

    // Main thread continues executing
    for i in 1..=3 {
        println!("Main: count {}", i);
        thread::sleep(Duration::from_millis(200));
    }

    // Wait for the spawned thread to finish
    handle.join().unwrap();
}

// Sharing data between threads using Arc (Atomic Reference Counting)
fn shared_state_threading() {
    // Arc allows multiple ownership across threads
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];

    for i in 0..3 {
        // Clone the Arc to create a new reference
        let counter = Arc::clone(&counter);

        // Spawn a new thread
        let handle = thread::spawn(move || {
            // Lock the mutex to access the data
            let mut num = counter.lock().unwrap();
            *num += 1;
            println!("Thread {}: counter is now {}", i, *num);
        });

        handles.push(handle);
    }

    // Wait for all threads to complete
    for handle in handles {
        handle.join().unwrap();
    }

    println!("Final counter value: {}", *counter.lock().unwrap());
}

// Message passing between threads using channels
fn message_passing() {
    // Create a channel (multiple producer, single consumer)
    let (tx, rx) = mpsc::channel();

    // Spawn a thread that sends messages
    thread::spawn(move || {
        let messages = vec!["Hello", "from", "the", "thread"];
        for msg in messages {
            tx.send(msg).unwrap();
            thread::sleep(Duration::from_millis(100));
        }
    });

    // Receive messages in the main thread
    for received in rx {
        println!("Got: {}", received);
    }
}

// Example of thread-safe data structures
fn thread_safe_collections() {
    use std::sync::RwLock;
    use std::collections::HashMap;

    // RwLock allows multiple readers or a single writer
    let map = Arc::new(RwLock::new(HashMap::new()));
    let mut handles = vec![];

    // Spawn writer threads
    for i in 0..3 {
        let map = Arc::clone(&map);
        handles.push(thread::spawn(move || {
            let mut map = map.write().unwrap();
            map.insert(i, i * i);
            println!("Thread {} wrote to map", i);
        }));
    }

    // Spawn reader threads
    for i in 0..3 {
        let map = Arc::clone(&map);
        handles.push(thread::spawn(move || {
            let map = map.read().unwrap();
            println!("Thread {} reading map: {:?}", i, *map);
        }));
    }

    for handle in handles {
        handle.join().unwrap();
    }
}
```

### Async/Await

Rust's async/await syntax provides a more efficient way to handle concurrent operations, especially I/O-bound tasks:

```rust
use std::time::Duration;
use tokio; // Popular async runtime

// Basic async function
async fn say_hello() {
    println!("Hello");
    // Simulate some async work
    tokio::time::sleep(Duration::from_secs(1)).await;
    println!("World");
}

// Async function that returns a value
async fn fetch_data(id: u32) -> String {
    // Simulate network request
    tokio::time::sleep(Duration::from_millis(100)).await;
    format!("Data for id {}", id)
}

// Running multiple async tasks concurrently
async fn concurrent_tasks() {
    // Join multiple futures
    let (result1, result2) = tokio::join!(
        fetch_data(1),
        fetch_data(2)
    );
    println!("Results: {}, {}", result1, result2);

    // Spawn multiple tasks
    let mut handles = vec![];
    for i in 0..3 {
        let handle = tokio::spawn(async move {
            let result = fetch_data(i).await;
            println!("Task {}: {}", i, result);
        });
        handles.push(handle);
    }

    // Wait for all tasks to complete
    for handle in handles {
        handle.await.unwrap();
    }
}

// Async streams
async fn process_stream() {
    use tokio_stream::StreamExt;
    use futures::stream;

    // Create a stream of numbers
    let mut stream = stream::iter(1..=5);

    // Process stream items
    while let Some(num) = stream.next().await {
        println!("Processing number: {}", num);
        tokio::time::sleep(Duration::from_millis(100)).await;
    }
}

// Error handling in async code
async fn async_error_handling() -> Result<String, Box<dyn std::error::Error>> {
    // Simulate a fallible operation
    let result = tokio::time::timeout(
        Duration::from_secs(1),
        fetch_data(1)
    ).await??;  // Note the double ? for both timeout and fetch_data errors

    Ok(result)
}

// Example of async main
#[tokio::main]
async fn main() {
    // Run basic async function
    say_hello().await;

    // Run concurrent tasks
    concurrent_tasks().await;

    // Process stream
    process_stream().await;

    // Handle errors
    match async_error_handling().await {
        Ok(data) => println!("Success: {}", data),
        Err(e) => println!("Error: {}", e),
    }
}
```

### Choosing Between Threads and Async

Use threads when:

- You have CPU-bound tasks
- You need to run blocking operations
- You want to utilize multiple CPU cores
- You need to run code that can't be made async

Use async when:

- You have I/O-bound tasks (network, file operations)
- You need to handle many concurrent operations
- You want to minimize resource usage
- You're building a web server or other I/O-heavy application

### Common Concurrency Patterns

1. **Thread Pool**

```rust
use rayon::prelude::*;  // Popular parallel iterator library

fn parallel_processing() {
    let numbers = vec![1, 2, 3, 4, 5];

    // Process items in parallel
    let sum: i32 = numbers.par_iter()
        .map(|&n| n * n)
        .sum();

    println!("Sum of squares: {}", sum);
}
```

2. **Async Mutex**

```rust
use tokio::sync::Mutex;

async fn async_shared_state() {
    let counter = Arc::new(Mutex::new(0));
    let mut handles = vec![];

    for i in 0..3 {
        let counter = Arc::clone(&counter);
        handles.push(tokio::spawn(async move {
            let mut lock = counter.lock().await;
            *lock += 1;
            println!("Task {}: counter is now {}", i, *lock);
        }));
    }

    for handle in handles {
        handle.await.unwrap();
    }
}
```

3. **Select Macro for Async Operations**

```rust
use tokio::select;

async fn select_example() {
    let mut interval = tokio::time::interval(Duration::from_secs(1));
    let mut timeout = tokio::time::sleep(Duration::from_secs(5));

    select! {
        _ = interval.tick() => {
            println!("Interval ticked");
        }
        _ = timeout => {
            println!("Timeout reached");
        }
    }
}
```

Remember:

- Always use appropriate synchronization primitives (Mutex, RwLock, etc.)
- Be careful with shared mutable state
- Consider using message passing for thread communication
- Use async for I/O-bound operations
- Use threads for CPU-bound operations
- Profile your application to ensure you're using the right approach

## Next Steps

- Read the [Rust Book](https://doc.rust-lang.org/book/) for in-depth coverage
- Practice with [Rustlings](https://github.com/rust-lang/rustlings) - small exercises to get you used to reading and writing Rust code
