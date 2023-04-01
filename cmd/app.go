package main

import (
	"fmt"
	"github.com/gictorbit/peershare/client"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	HostAddress    string
	ServerPort     uint
	LogRequest     bool
	SendFilePath   string
	ReceiveOutPath string
	SharedCode     string
)

func main() {
	// using the function
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	app := &cli.App{
		Name:  "peershare",
		Usage: "p2p file sharing system",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "host address",
				Value:       "0.0.0.0",
				DefaultText: "0.0.0.0",
				EnvVars:     []string{"HOST_ADDRESS"},
				Destination: &HostAddress,
			},
			&cli.UintFlag{
				Name:        "port",
				Usage:       "server port",
				Value:       3000,
				DefaultText: "3000",
				Aliases:     []string{"p"},
				EnvVars:     []string{"SERVER_PORT"},
				Destination: &ServerPort,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "run file sharing signaling server",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "log-request",
						Usage:       "log incoming requests",
						Aliases:     []string{"lgr"},
						Value:       false,
						DefaultText: "false",
						EnvVars:     []string{"LOG_REQUEST"},
						Destination: &LogRequest,
					},
				},
				Action: func(cliCtx *cli.Context) error {
					serverAddr := net.JoinHostPort(HostAddress, fmt.Sprintf("%d", ServerPort))
					server := RunSignalingGRPCServer(serverAddr)
					stop := make(chan os.Signal, 1)
					signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
					<-stop

					log.Println("shutting down servers")
					server.GracefulStop()
					return nil
				},
			},
			{
				Name:  "client",
				Usage: "run file sharing system client",
				Subcommands: []*cli.Command{
					{
						Name:  "send",
						Usage: "send file to another peer",
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:        "file",
								Usage:       "path of file",
								Aliases:     []string{"f"},
								Required:    true,
								Destination: &SendFilePath,
							},
						},
						Action: func(context *cli.Context) error {
							serverAddr := net.JoinHostPort(HostAddress, fmt.Sprintf("%d", ServerPort))
							opts := []grpc.DialOption{
								grpc.WithTransportCredentials(insecure.NewCredentials()),
							}
							conn, err := grpc.Dial(serverAddr, opts...)
							if err != nil {
								return err
							}
							peerClient := client.NewPeerClient(conn)
							peerClient.SendFile(SendFilePath)
							return nil
						},
					},
					{
						Name:  "receive",
						Usage: "receive file to another peer",
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:        "out",
								Usage:       "output directory for downloaded file",
								Aliases:     []string{"o"},
								DefaultText: pwd,
								Value:       pwd,
								Required:    false,
								Destination: &ReceiveOutPath,
							},
							&cli.StringFlag{
								Name:        "code",
								Usage:       "shared code to receive file from another peer",
								Required:    true,
								Aliases:     []string{"c"},
								Destination: &SharedCode,
							},
						},
						Action: func(context *cli.Context) error {
							serverAddr := net.JoinHostPort(HostAddress, fmt.Sprintf("%d", ServerPort))
							opts := []grpc.DialOption{
								grpc.WithTransportCredentials(insecure.NewCredentials()),
							}
							conn, err := grpc.Dial(serverAddr, opts...)
							if err != nil {
								return err
							}
							peerClient := client.NewPeerClient(conn)
							peerClient.ReceiveFile(SharedCode, ReceiveOutPath)
							return nil
						},
					},
				},
			},
		},
	}
	if e := app.Run(os.Args); e != nil {
		log.Println("failed to run app", e)
	}
}
