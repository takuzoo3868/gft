package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/takuzoo3868/gft/proto"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func listFiles(ctx context.Context, client proto.FileTransferServiceClient) error {
	slist, err := client.ListFiles(ctx, new(proto.ListRequestType))
	if err != nil {
		return err
	}
	fmt.Println("name,size,mode,modtime")
	for {
		file, err := slist.Recv()
		if err != nil {
			break
		}
		fmt.Printf("%q,\"%v\",\"%v\",\"%v\"\n",
			file.Name, file.Size, os.FileMode(file.Mode), time.Unix(file.ModTime.Seconds, 0).Format(time.RFC3339))
	}
	slist.CloseSend()
	return err
}

func ListCommand() cli.Command {
	return cli.Command{
		Name:  "ls",
		Usage: "list files from server by CSV format",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "a",
				Value: ":8080",
				Usage: "server address",
			},
		},
		Action: func(c *cli.Context) error {
			options := []grpc.DialOption{}
			addr := c.String("a")
			if !strings.Contains(addr, ":") {
				addr += ":8080"
			}
			conn, err := grpc.Dial(addr, options...)
			if err != nil {
				log.Fatalf("cannot connect: %v", err)
			}
			defer conn.Close()

			return listFiles(context.Background(), proto.NewFileTransferServiceClient(conn))
		},
	}
}
