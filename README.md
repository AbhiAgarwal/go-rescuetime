### Go Rescuetime

A simple client to interface with the Rescue Time API. A godoc for this is [here](https://godoc.org/github.com/AbhiAgarwal/go-rescuetime).

Sample use:

```go
var rescue RescueTime
rescue.apiKey = ""
summary, err := rescue.DailySummary()
```