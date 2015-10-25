[![GoDoc](https://godoc.org/github.com/AbhiAgarwal/go-rescuetime?status.svg)](https://godoc.org/github.com/AbhiAgarwal/go-rescuetime)

### Go RescueTime

A simple client to interface with the RescueTime API. For usage, please see the GoDocs [here](https://godoc.org/github.com/AbhiAgarwal/go-rescuetime).

Sample use:

```go
var rescue RescueTime
rescue.APIKey = ""
summary, err := rescue.GetDailySummary()
```