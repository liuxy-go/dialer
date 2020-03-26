package dialer

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strings"

	"golang.org/x/crypto/ssh"
)

type sshDialer struct {
	addr string
	cfg  *ssh.ClientConfig
	cli  *ssh.Client
}

type badParamError struct {
	p, v string
}

func (e badParamError) Error() string {
	return fmt.Sprintf("bad param %s: %s", e.p, e.v)
}

// SSH 根据 ssh 连接和密钥文件路径, 创建 dialer
//	登录模式:
//		password - uri: "ssh://[user]:pass@host:[port]"
//          user 默认为 `root`, port 默认为 `22`
//      keypair  - uri: "ssh://[user]:[pass]@host:[port]"
//          user 默认为 `root`, pass 是密钥的密码, port 默认为 `22`
func SSH(uri, keyfile string) (dl ContextDialer, er error) {
	var (
		user, pass, addr string
	)
	if !strings.HasPrefix(uri, "ssh:") {
		uri = "ssh://" + uri
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	user = u.User.Username()
	pass, _ = u.User.Password()
	addr = u.Host
	return newSSHDialer(user, addr, keyfile, pass)
}

func newSSHDialer(user, addr, keyfile, pass string) (dl ContextDialer, er error) {
	if len(addr) < 1 {
		return nil, badParamError{"addr", addr}
	}
	if len(keyfile) < 1 {
		if len(pass) < 1 {
			return nil, badParamError{"pass", pass}
		}
	}
	if _, p := splitHostPort(addr); len(p) < 0 {
		addr = addr + ":22" // default port
	}

	var (
		err  error
		auth []ssh.AuthMethod
		key  []byte
	)

	if len(keyfile) < 1 {
		auth = append(auth, ssh.Password(pass))
	} else {
		var sig ssh.Signer
		key, err = ioutil.ReadFile(keyfile)
		if err != nil {
			return nil, err
		}
		if len(pass) > 0 {
			sig, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(pass))
		} else {
			sig, err = ssh.ParsePrivateKey(key)
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(sig))
	}
	cfg := &ssh.ClientConfig{
		User: user, Auth: auth, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cli, err := ssh.Dial("tcp", addr, cfg)
	return &sshDialer{addr, cfg, cli}, err
}

func (s *sshDialer) Dial(addr string) (net.Conn, error) {

	return s.cli.Dial("tcp", addr)
}

func (s *sshDialer) DialContext(_ context.Context, addr string) (net.Conn, error) {
	return s.cli.Dial("tcp", addr)
}
