## dialer

一些常用的 Dialer (内部实现确定了是 tcp 还是 udp)

接口:

```go
type Dialer interface {
	Dial(string) (net.Conn, error)
}

type ContextDialer interface {
	Dialer
	DialContext(context.Context, string) (net.Conn, error)
}
```