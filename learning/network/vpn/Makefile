PASS=8811

setup:
#	echo ${PASS} | sudo -S ip netns add mock_internet_ns
#	sudo ip netns exec mock_internet_ns ip link set lo up

#	sudo ip link add vpn_source_veth type veth peer name vpn_dest_veth netns mock_internet_ns

#	sudo ip addr add 192.168.2.1/28 dev vpn_source_veth
#	sudo ip netns exec mock_internet_ns ip addr add 192.168.2.8/28 dev vpn_dest_veth
#	sudo ip netns exec mock_internet_ns ip addr add 192.168.2.11/28 dev vpn_dest_veth
#	sudo ip link set vpn_source_veth up
#	sudo ip netns exec mock_internet_ns ip link set vpn_dest_veth up


run-server-client:
	sudo ip netns exec mock_internet_ns /home/erema/h/go_sandbox/learning/network/vpn/server/server 192.168.2.8:33033
	sudo /home/erema/h/go_sandbox/learning/network/vpn/client/client 192.168.2.8:33033 tun_test 192.168.2.11


reset:
	echo ${PASS} | sudo -S ip netns delete mock_internet_ns

get:
	ip route get ${SERVER_IP}

tun5:
	echo ${PASS} | sudo -S ip tuntap add dev tun5 mode tun
	sudo ip link set dev tun5 up
	sudo ip route add ${SERVER_IP} dev tun5
	sudo ip addr add 10.0.5.1 peer 10.0.5.2 dev tun5

tun5-down:
	sudo ip tuntap del dev tun5 mode tun


reset-test:
	echo ${PASS} |sudo -S ip link del tun_test
	sudo ip link del vpn_source_veth
	sudo ip netns delete mock_internet_ns
	sudo ip route delete 192.168.2.11
	sudo ip route delete 192.168.2.1/28