package vpn_test

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"github.com/r-erema/go_sendbox/utils"
	"github.com/r-erema/go_sendbox/utils/net"
)

func TestVPN(t *testing.T) { //nolint:paralleltest
	RegisterFailHandler(Fail)

	RunSpecs(t, "todo")
}

const (
	vethSourceName      = "vpn_source_veth"
	vethDestinationName = "vpn_dest_veth"
	destinationNs       = "mock_internet_ns"
)

//nolint:gochecknoglobals
var (
	clientCmd    *exec.Cmd
	serverStdErr io.ReadCloser
)

var _ = BeforeSuite(func() {
	err := net.CreateNS(destinationNs)
	Expect(err).NotTo(HaveOccurred())

	_, err = net.SetupVeth(vethSourceName, "192.168.2.1/28", vethDestinationName, "192.168.2.8/28", destinationNs)
	Expect(err).NotTo(HaveOccurred())

	port, err := freeport.GetFreePort()
	Expect(err).NotTo(HaveOccurred())

	serverAddr := fmt.Sprintf("192.168.2.8:%d", port)

	ctx := context.Background()

	_, _, serverStdErr, err = utils.GoCompileAndStartInNetNs(
		ctx,
		"./server/main.go",
		destinationNs,
		serverAddr,
	)
	Expect(err).NotTo(HaveOccurred())

	_, err = fmt.Fprintf(GinkgoWriter, "VPN Server Started at: %s\n", serverAddr)
	Expect(err).NotTo(HaveOccurred())

	clientCmd, _, _, err = utils.GoCompileAndStart(
		ctx,
		"./client/main.go",
		serverAddr,
		"tun_test",
		"192.168.2.11",
	)
	Expect(err).NotTo(HaveOccurred())

	_, err = fmt.Fprintf(GinkgoWriter, "VPN Server Started at: %s\n", serverAddr)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := clientCmd.Process.Signal(syscall.SIGTERM)
	Expect(err).NotTo(HaveOccurred())

	err = clientCmd.Wait()
	Expect(err).NotTo(HaveOccurred())

	err = net.DeleteLink(vethSourceName)
	Expect(err).NotTo(HaveOccurred())

	err = net.DeleteNS(destinationNs)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("VPN", func() {
	It("Packets from ping reach the Server", func() {
		pingCmd := exec.Command("ping", "192.168.2.11")
		err := pingCmd.Start()
		Expect(err).NotTo(HaveOccurred())
		GinkgoT().Cleanup(func() {
			err = pingCmd.Process.Kill()
			Expect(err).NotTo(HaveOccurred())
		})

		Eventually(func(g Gomega) {
			buf := make([]byte, 1024)
			n, err := serverStdErr.Read(buf)
			g.Expect(err).NotTo(HaveOccurred())
			buf = buf[:n]
			g.Expect(buf).To(ContainSubstring("Destination IP: 192.168.2.11"))
		}, time.Second*5, time.Millisecond*500).Should(Succeed())
	})
})
