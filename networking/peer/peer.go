package peer

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type Peer struct {
	Name string
	Group string
	peerDiscovery
	isListening bool
}

type DiscoveryOption func(*Settings)

func New(name string, group string, options ...DiscoveryOption) (*Peer, error) {
	p := &Payload{
		Name:  name,
		Group: group,
	}

	payload, err := encodePayload(p)
	if err != nil {
		return nil, err
	}

	s := Settings{}
	for _, opt := range options {
		opt(&s)
	}

	peer, err := initialize(s)
	if err != nil {
		return nil, err
	}

	peer.Name = name
	peer.Group = group
	peer.settings.Payload = payload

	return peer, nil
}

func (p *Peer) StopListening() {
	p.Lock()
	p.isListening = false
	p.Unlock()
	return
}

func (p *Peer) Listen() error {
	address := net.JoinHostPort(p.settings.MulticastAddress, p.settings.Port)
	portNum := p.settings.portNum
	allowSelf := p.settings.AllowSelf
	notify := p.settings.Notify

	localIPs := getLocalIPs()

	// get interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// Open up a connection
	c, err := net.ListenPacket(fmt.Sprintf("udp%d", p.settings.IPVersion), address)
	if err != nil {
		return err
	}
	defer c.Close()

	group := p.settings.multicastAddressNumbers
	var p2 NetPacketConn
	if p.settings.IPVersion == IPv4 {
		p2 = PacketConn4{ipv4.NewPacketConn(c)}
	} else {
		p2 = PacketConn6{ipv6.NewPacketConn(c)}
	}

	for i := range ifaces {
		p2.JoinGroup(&ifaces[i], &net.UDPAddr{IP: group, Port: portNum})
	}

	p.Lock()
	p.isListening = true
	p.Unlock()

	go func() {
		for {
			if !p.isListening {
				return
			}

			buffer := make([]byte, maxDatagramSize)
			var (
				n       int
				src     net.Addr
				errRead error
			)
			n, src, errRead = p2.ReadFrom(buffer)
			if errRead != nil {
				err = errRead
				return
			}

			srcHost, _, _ := net.SplitHostPort(src.String())

			if _, ok := localIPs[srcHost]; ok && !allowSelf {
				continue
			}

			// log.Println(src, hex.Dump(buffer[:n]))

			p.Lock()
			if _, ok := p.received[srcHost]; !ok {
				p.received[srcHost] = buffer[:n]
			}
			p.Unlock()

			payload, err := decodePayload(buffer[:n])
			if err != nil {
				continue
			}

			discovered := Discovered{
				Address: srcHost,
				Payload: *payload,
			}

			if payload.Group != p.Group {
				continue
			}

			if notify != nil {
				notify(discovered)
			}

			_ = p.Broadcast()
		}
	}()

	return nil
}

func (p *Peer) Broadcast() error {
	address := net.JoinHostPort(p.settings.MulticastAddress, p.settings.Port)
	portNum := p.settings.portNum
	payload := p.settings.Payload

	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// Open up a connection
	c, err := net.ListenPacket(fmt.Sprintf("udp%d", p.settings.IPVersion), address)
	if err != nil {
		return err
	}
	defer c.Close()

	group := p.settings.multicastAddressNumbers

	// ipv{4,6} have an own PacketConn, which does not implement net.PacketConn
	var p2 NetPacketConn
	if p.settings.IPVersion == IPv4 {
		p2 = PacketConn4{ipv4.NewPacketConn(c)}
	} else {
		p2 = PacketConn6{ipv6.NewPacketConn(c)}
	}

	for i := range ifaces {
		p2.JoinGroup(&ifaces[i], &net.UDPAddr{IP: group, Port: portNum})
	}

	for i := range ifaces {
		if errMulticast := p2.SetMulticastInterface(&ifaces[i]); errMulticast != nil {
			continue
		}
		p2.SetMulticastTTL(2)
		if _, errMulticast := p2.WriteTo([]byte(payload), &net.UDPAddr{IP: group, Port: portNum}); errMulticast != nil {
			continue
		}
	}

	return nil
}

func WithSettings(settings Settings) DiscoveryOption {
	return func(s *Settings) {
		s = &settings
	}
}

func SetLimit(val int) DiscoveryOption {
	return func(s *Settings) {
		s.Limit = val
	}
}

func SetPort(val int) DiscoveryOption {
	return func(s *Settings) {
		s.Port = strconv.Itoa(val)
	}
}

func SetMulticastAddress(val string) DiscoveryOption {
	return func(s *Settings) {
		s.MulticastAddress = val
	}
}

//SetPayload to set the payload content for the discoveries. If omitted, then the peer name will be used.
func SetPayload(val []byte) DiscoveryOption {
	return func(s *Settings) {
		s.Payload = val
	}
}

func SetDelay(val time.Duration) DiscoveryOption {
	return func(s *Settings) {
		s.Delay = val
	}
}

func SetTimeLimit(val time.Duration) DiscoveryOption {
	return func(s *Settings) {
		s.TimeLimit = val
	}
}

func SetStopChan(val chan struct{}) DiscoveryOption {
	return func(s *Settings) {
		s.StopChan = val
	}
}

func SetAllowSelf(val bool) DiscoveryOption {
	return func(s *Settings) {
		s.AllowSelf = val
	}
}

func SetDisableBroadcast(val bool) DiscoveryOption {
	return func(s *Settings) {
		s.DisableBroadcast = val
	}
}

func SetIPVersion(val IPVersion) DiscoveryOption {
	return func(s *Settings) {
		s.IPVersion = val
	}
}

func SetNotify(val func(Discovered)) DiscoveryOption {
	return func(s *Settings) {
		s.Notify = val
	}
}
