package vpn

type addrinfo struct {
	ai_socktype int
	ai_protocol int
}

func udpBind() int {
	// unix.SOCK_DGRAM
	return 0
}
