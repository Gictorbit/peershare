package main

import (
	"fmt"
	logcharm "github.com/charmbracelet/log"
	"github.com/gictorbit/peershare/api"
	"github.com/gictorbit/peershare/client"
	"github.com/gictorbit/peershare/sigserver"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	HostAddress    string
	ServerPort     uint
	LogStacktrace  bool
	SendFilePath   string
	ReceiveOutPath string
	SharedCode     string
)

func main() {
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
				Value:       "127.0.0.1",
				DefaultText: "127.0.0.1",
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
						Name:        "stacktrace",
						Usage:       "enable/disable stacktrace",
						Aliases:     []string{"str"},
						Value:       false,
						DefaultText: "false",
						EnvVars:     []string{"LOG_STACK_TRACE"},
						Destination: &LogStacktrace,
					},
				},
				Action: func(cliCtx *cli.Context) error {
					serverAddr := net.JoinHostPort(HostAddress, fmt.Sprintf("%d", ServerPort))
					logCfg := zap.NewProductionConfig()
					logCfg.DisableStacktrace = !LogStacktrace
					logger, err := logCfg.Build()
					if err != nil {
						log.Fatalf("create new logger failed:%v\n", err)
					}
					server := sigserver.NewPeerShareServer(serverAddr, logger)
					go server.Start()

					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					server.Stop()
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
							logger := logcharm.New(os.Stderr)
							logger.SetReportCaller(true)
							serverAddr := net.JoinHostPort(HostAddress, fmt.Sprintf("%d", ServerPort))
							logger.Info("server address is ", "addr", serverAddr)
							peerClient := client.NewPeerClient(serverAddr, api.SenderClient, logger)
							if e := peerClient.Connect(); e != nil {
								return e
							}
							if e := peerClient.SendFile(SendFilePath); e != nil {
								return e
							}
							peerClient.Stop()
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
							logger := logcharm.New(os.Stderr)
							logger.SetReportCaller(true)
							if _, e := os.Stat(ReceiveOutPath); os.IsNotExist(e) {
								if mkdirErr := os.MkdirAll(ReceiveOutPath, os.ModePerm); mkdirErr != nil {
									return mkdirErr
								}
								logger.Info("output directory created", "path", ReceiveOutPath)
							}
							logger.Info("server address is ", "addr", serverAddr)
							peerClient := client.NewPeerClient(serverAddr, api.ReceiverClient, logger)
							if e := peerClient.Connect(); e != nil {
								return e
							}
							if e := peerClient.ReceiveFile(SharedCode, ReceiveOutPath); e != nil {
								return e
							}
							peerClient.Stop()
							return nil
						},
					},
				},
			},
		},
	}
	if e := app.Run(os.Args); e != nil {
		logger := logcharm.New(os.Stderr)
		logger.Error("failed to run app", "error", e)
	}
}
