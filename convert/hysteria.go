package convert

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/xmdhs/clash2singbox/model/clash"
	"github.com/xmdhs/clash2singbox/model/singbox"
)

func hysteria(p *clash.Proxies, s *singbox.SingBoxOut) error {
	p.Tls = true
	tls(p, s)
	if p.Port == "" {
		return fmt.Errorf("hysteria: %w", ErrNotSupportType)
	}
	if p.AuthStr != "" {
		s.AuthStr = p.AuthStr
	} else {
		s.AuthStr = p.AuthStr1
	}
	if p.Obfs != "" {
		s.Obfs = &singbox.SingObfs{
			Value: p.Obfs,
		}
	}
	s.TLS.Alpn = p.Alpn
	if p.Protocol != "udp" && p.Protocol != "" {
		return fmt.Errorf("hysteria: %w", ErrNotSupportType)
	}
	if up, err := strconv.Atoi(p.Up); err == nil {
		s.UpMbps = up
	} else {
		s.Up = p.Up
	}
	if down, err := strconv.Atoi(p.Down); err == nil {
		s.DownMbps = down
	} else {
		s.Down = p.Down
	}
	if p.RecvWindow != 0 {
		s.RecvWindow = int(p.RecvWindow)
	} else {
		s.RecvWindow = int(p.RecvWindow1)
	}
	if p.RecvWindowConn != 0 {
		s.RecvWindowConn = int(p.RecvWindowConn)
	} else {
		s.RecvWindowConn = int(p.RecvWindowConn1)
	}
	if p.CaStr != "" {
		s.TLS.Certificate = p.CaStr
	} else {
		s.TLS.Certificate = p.CaStr1
	}
	s.DisableMtuDiscovery = bool(p.DisableMtuDiscovery)
	s.HopPorts = p.Ports
	return nil
}

func hysteia2(p *clash.Proxies, s *singbox.SingBoxOut) ([]singbox.SingBoxOut, error) {
	p.Tls = true
	tls(p, s)
	var err error
	s.UpMbps, err = anyToMbps(p.Up)
	if err != nil {
		return nil, fmt.Errorf("hysteia2: %w", err)
	}
	s.DownMbps, err = anyToMbps(p.Down)
	if err != nil {
		return nil, fmt.Errorf("hysteia2: %w", err)
	}
	s.Password = p.Password
	if p.ObfsPassword != "" {
		s.Obfs = &singbox.SingObfs{
			Type:  p.Obfs,
			Value: p.ObfsPassword,
		}
	}
	s.HopPorts = p.Ports
	return []singbox.SingBoxOut{*s}, nil
}

var rateStringRegexp = regexp.MustCompile(`^(\d+)\s*([KMGT]?)([Bb])ps$`)

func anyToMbps(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	if mb, err := strconv.Atoi(s); err == nil {
		return mb, nil
	}

	m := rateStringRegexp.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("anyToMbps: %w", ErrNotSupportType)
	}

	n := 1.0
	switch m[2] {
	case "K":
		n = 1.0 / 1000.0
	case "M":
		n = 1
	case "G":
		n = 1000
	case "T":
		n = 1000 * 1000
	}
	if m[3] == "B" {
		n = n * 8.0
	}
	v, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, fmt.Errorf("anyToMbps: %w", ErrNotSupportType)
	}
	mb := int(float64(v) * n)
	if mb == 0 {
		mb = 1
	}
	return mb, nil
}
