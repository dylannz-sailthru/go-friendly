# go-friendly

go-friendly is a go package that lets you accompany your internal error
messages with more friendly messages, that are suitable for and end-user to
see.

[Documentation](https://godoc.org/github.com/dylannz-sailthru/go-friendly)

License: MIT

## Example usage

For example, you may have encountered a database connection error that you want
to log out, but don't want to return to the user. That message might look like:

> error connecting to postgres db with host 127.0.0.1 and username 'my_user'

This might be useful for debugging, but might be too much information to
display to your end user. So when generating the error, you can also choose to
generate a distinct user-friendly error that might tell them a little about the
error and nothing that they don't need to know, e.g. you might show them
something like:

> We encountered an internal error while trying to process your request. Please try again later :(

Using go-friendly, this would look something like:

```go
import (
  "github.com/dylannz-sailthru/go-friendly"
)

func fetchSomethingFromTheDatabase() error {
  return errors.New("error connecting to postgres db with host 127.0.0.1 and username 'my_user'")
}

func fetchTime() error {
  err := fetchSomethingFromTheDatabase()

  return friendly.New().
    WithCause(err).
    WithFriendly("We encountered an internal error while trying to process your request. Please try again later :(")
}
```

Now, when you call `err.Error()`, you'll get see the internal error:

```go
err := fetchTime()
fmt.Println(err.Error()) // error connecting to postgres db with host 127.0.0.1 and username 'my_user'
```

But if you call friendly.Friendly() on the error, it'll return your more
user-friendly error:

```go
err := fetchTime()
if f := friendly.Friendly(err); f != nil {
  fmt.Println(f.Error()) // We encountered an internal error while trying to process your request. Please try again later :(
}
```