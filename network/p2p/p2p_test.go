package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"
	"os"
	"testing"
)

func parseFlags() *config {
	c := &config{}

	flag.StringVar(&c.Rendezvous, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.listenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.StringVar(&c.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")
	flag.StringVar(&c.listenPort, "port", "4001", "node listen port")

	flag.Parse()
	return c
}

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream")
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[36m%s\x1b[0m> ", str)
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			panic(err)
		}
	}
}

func Test(t *testing.T) {
	cfg := parseFlags()
	fmt.Printf("[*] Listening on: %s with port: %s\n", cfg.listenHost, cfg.listenPort)

	ctx := context.Background()
	r := rand.Reader

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, cfg.listenPort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(host)
	//
	//// Set a function as stream handler.
	//// This function is called when a peer initiates a connection and starts a stream with this peer.
	//host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)
	//fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%s/p2p/%s\n", cfg.listenHost, cfg.listenPort, host.ID().Pretty())
	//peerChan := initMDNS(ctx, host, cfg.Rendezvous)
	//peer := <- peerChan  // will block, until discover a peer
	//fmt.Println("Found peer:", peer, ", connecting")
	//if err := host.Connect(ctx, peer); err != nil {
	//	fmt.Println("Connection failed", err)
	//}
	//// open a stream, this stream will be handled by handleStream other end
	//stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))
	//if err != nil {
	//	fmt.Println("Stream open failed", err)
	//} else {
	//	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	//
	//	go writeData(rw)
	//	go readData(rw)
	//	fmt.Println("Connected to:", peer)
	//}
	//
	//select {
	//
	//}
}