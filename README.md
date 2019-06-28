# Go async worker
The Golang module that provides functions to handle jobs concurrently.

## Getting Started
1) Download the project from github to your host machine.
2) Go to the folder with project

## Prerequisites
For the successful using you should have:
```
go >= 1.12
```
## Running the tests
It's a good practice to run the tests before using the module to make sure everything is OK.
```
go test -v
```
## Sample of using
```go
hwFunc := func() error {
    fmt.Println("Hello world!")
    return nil
}
err := Handle([]JobFunc{hwFunc,hwFunc}, 2, 0)
...
```
## License
This project is licensed under the MIT License.