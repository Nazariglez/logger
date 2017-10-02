logger
====

Logger is just a simple log system for golang applications. Allow print in the terminal and in a 
file at the same time.  

## Install
```go get -u github.com/Nazariglez/logger```

## How to use
See `example` folder.

## Why an Async mode?
Some apps needs a way to execute logs asynchronously without block the main thread, in go you can do this easily with the 
`go` keyword, but the execution order are unpredictable. Using `logger.NewAsync()` the order how the logs 
are printed in the terminal and on the file are always the same.

## Formatting text
To format a text just use one of the log functions with `f` in the end, for example, for `Debug` will be `Debugf`.

## Is thread-safe?
Yes.

## How are saved the files?
You must define the path, and the logs will be sorted in files by days. Each day will have their own file. The day 2017-11-01 will have a file like this: `logname.20171101.log`. 