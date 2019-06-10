# Go Limit
Rate limit for Go

### Usage

```go
limit := golimit.NewLimiter(3 * time.Minute, 35) // Maximum 35 events for 3 minutes

if limit.IsBanned("user string identifier") {
    w.WriteHeader(http.StatusTooManyRequests)
    w.Write([]byte("429 - Too many requests"))
    return
}

if limit.Allow("user string identifier") == false {
    // User has reached its quota

    limit.Ban("user string identifier", 60 * time.Second) // You can also ban user
}
```
