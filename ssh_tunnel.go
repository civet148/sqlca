package sqlca

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/tcpsock" //required (register socket instance)
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2/types"
	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
)

const (
	sshTunnelPort = 6033 //
)

type SSH struct {
	User       string               //SSH tunnel server login account
	Password   string               //SSH tunnel server login password
	PrivateKey string               //SSH tunnel server private key, eg. "/home/test/.ssh/private-key.pem"
	Host       string               //SSH tunnel server host [ip or domain], default port 22 if not specified
	listenPort int                  //SSH transfer port of localhost
	listenIP   string               //SSH transfer ip of localhost
	tunnel     *sshtunnel.SSHTunnel //SSH tunnel instance
	authMethod ssh.AuthMethod       //SSH tunnel auth method
}

func (s *SSH) GoString() string {
	return s.String()
}

func (s *SSH) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}

func (s *SSH) setDefaultPort() {
	if -1 == strings.Index(s.Host, ":") {
		s.Host += ":22"
	}
}
func (s *SSH) openSSHTunnel(dsn *dsnDriver) (d *dsnDriver) {
	var adapter = types.GetAdapterType(dsn.strDriverName)

	if ok := s.startSSHTunnel(dsn); !ok {
		err := fmt.Errorf("start SSH tunnel service failed, please make sure your tunnel server config [%+v] is correct", s)
		log.Errorf(err.Error())
		panic(err.Error())
	}
	switch adapter {
	case types.AdapterSqlx_MySQL:
		d = s.buildMysqlTunnelDSN(dsn)
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		d = s.buildPostgresTunnelDSN(dsn)
	case types.AdapterSqlx_Mssql:
		d = s.buildMssqlTunnelDSN(dsn)
	default:
		d = dsn
	}
	log.Debugf("openSSHTunnel [%s] ssh [%+v] tunnel DSN driver [%+v]", adapter.String(), s, d)
	return d
}

func (s *SSH) startSSHTunnel(dsn *dsnDriver) (ok bool) {

	s.setDefaultPort()

	log.Debugf("try connect to SSH tunnel server [%s] ok", s.Host)
	if s.listenPort, ok = s.tryListenRandomPort(); !ok {
		log.Errorf("try listen random port for SSH tunnel failed")
		return
	}

	if s.PrivateKey != "" {
		s.authMethod = sshtunnel.PrivateKeyFile(s.PrivateKey)
	} else {
		s.authMethod = ssh.Password(s.Password)
	}

	var strTunnelHost = s.User + "@" + s.Host
	s.tunnel = sshtunnel.NewSSHTunnel(
		// User and host of tunnel server, it will default to port 22 if not specified.
		strTunnelHost,

		// Pick ONE of the following authentication methods:
		// 1. ssh.Password("123456")
		// 2. sshtunnel.PrivateKeyFile("path/to/private/key.pem")
		s.authMethod,

		// The destination host and port of the actual server.
		dsn.parameter.host,

		// The local port you want to bind the remote port to. specifying "0" will lead to a random port.
		fmt.Sprintf("%d", s.listenPort),
	)

	//make sure tunnel server is reachable
	if err := s.tryConnectSSH(); err != nil {
		log.Errorf("try connect to SSH tunnel server [%s] failed", s.Host)
		return false
	}

	//start tunnel service
	go s.start()

	//make sure local tunnel transfer port is ready
	var strLocalAddr = fmt.Sprintf("%s:%d", s.listenIP, s.listenPort)
	if err := s.tryConnect(strLocalAddr); err != nil {
		log.Errorf("connect to transfer address [%s] error", strLocalAddr)
		return false
	}
	log.Debugf("try connect to transfer address [%s] ok", strLocalAddr)
	return true
}

func (s *SSH) start() (err error) {
	//s.tunnel.Log = logger.New(os.Stdout, "", logger.Ldate|logger.Lmicroseconds)
	if err = s.tunnel.Start(); err != nil {
		log.Errorf("start tunnel service error [%s]", err)
		panic(err.Error())
	}
	return
}

func (s *SSH) tryConnectSSH() (err error) {
	var c *ssh.Client
	if c, err = ssh.Dial("tcp", s.Host, s.tunnel.Config); err != nil {
		log.Errorf("try connect ssh [%s] config [%+v] error [%s]", s.Host, s.tunnel.Config, err.Error())
		return
	}
	defer c.Close()
	return
}

func (s *SSH) tryConnect(strHost string) (err error) {

	var strConnect = fmt.Sprintf("tcp://%s", strHost)
	c := wss.NewClient()
	if err = c.Connect(strConnect); err != nil {
		log.Errorf("connect to [%s] error", strHost)
		return
	}
	_ = c.Close()
	return
}

func (s *SSH) tryListenRandomPort() (port int, ok bool) {

	s.listenIP = "127.0.0.1"
	port = sshTunnelPort

	for i := 0; i < 1000; i++ {
		port += i
		var strListenAddr = fmt.Sprintf("%s:%d", s.listenIP, port)
		if listener, err := net.Listen("tcp", strListenAddr); err == nil {
			log.Debugf("tunnel service try listen [%s] ok", strListenAddr)
			_ = listener.Close()
			ok = true
			break
		}
	}
	return
}

func (s *SSH) buildMysqlTunnelDSN(dsn *dsnDriver) *dsnDriver {
	log.Debugf("buildMysqlTunnelDSN dsn [%+v] ssh [%+v]", dsn, s)
	var kvs []string
	var parameter = &dsn.parameter
	parameter.strDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?", parameter.user, parameter.password, s.listenIP, s.listenPort, parameter.db)
	for k, v := range parameter.queries {
		kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
	}
	if len(kvs) > 0 {
		parameter.strDSN += strings.Join(kvs, "&")
	}
	return dsn
}

func (s *SSH) buildPostgresTunnelDSN(dsn *dsnDriver) *dsnDriver {
	log.Debugf("buildPostgresTunnelDSN dsn [%+v] ssh [%+v]", dsn, s)
	var strListenPort = fmt.Sprintf("%d", s.listenPort)
	var parameter = &dsn.parameter
	parameter.strDSN = buildPostgresDSN(s.listenIP, strListenPort, parameter.user, parameter.password, parameter.db, parameter.queries)
	return dsn
}

func (s *SSH) buildMssqlTunnelDSN(dsn *dsnDriver) *dsnDriver {
	log.Debugf("buildMssqlTunnelDSN dsn [%+v] ssh [%+v]", dsn, s)
	var strListenPort = fmt.Sprintf("%d", s.listenPort)
	var parameter = &dsn.parameter
	parameter.strDSN = buildMssqlDSN(s.listenIP, strListenPort, parameter.user, parameter.password, parameter.db, parameter.queries)
	return dsn
}

