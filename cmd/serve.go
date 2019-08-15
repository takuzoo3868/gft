package cmd

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/takuzoo3868/gft/proto"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

type fileTransferService struct {
	root string
}

func (fts *fileTransferService) ListFiles(_ *proto.ListRequestType, stream proto.FileTransferService_ListFilesServer) error {
	err := filepath.Walk(fts.root, func(p string, info os.FileInfo, err error) error {
		name, err := filepath.Rel(fts.root, p)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)
		modTime := new(google_protobuf.Timestamp)
		modTime.Seconds = int64(info.ModTime().Unix())
		modTime.Nanos = int32(info.ModTime().UnixNano())
		f := &proto.ListResponseType{Name: name, Size: info.Size(), Mode: uint32(info.Mode()), ModTime: modTime}
		return stream.Send(f)
	})
	return err
}

func (fts *fileTransferService) Download(r *proto.DownloadRequestType, stream proto.FileTransferService_DownloadServer) error {
	f, err := os.Open(filepath.Join(fts.root, r.Name))
	if err != nil {
		return err
	}
	defer f.Close()

	var b [4096 * 1000]byte
	for {
		n, err := f.Read(b[:])
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		err = stream.Send(&proto.DownloadResponseType{
			Data: b[:n],
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func ServeCommand() cli.Command {
	return cli.Command{
		Name:  "up",
		Usage: "upload files on local server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "a",
				Value: ":8080",
				Usage: "server address",
			},
			cli.StringFlag{
				Name:  "d",
				Value: ".",
				Usage: "base directory to upload files",
			},
		},
		Action: func(c *cli.Context) error {
			lis, err := net.Listen("tcp", c.String("a"))
			if err != nil {
				log.Fatalf("cannot listen: %v", err)
			}
			defer lis.Close()

			options := []grpc.ServerOption{}
			log.Println("server started:", lis.Addr().String())
			server := grpc.NewServer(options...)
			proto.RegisterFileTransferServiceServer(server, &fileTransferService{root: c.String("d")})
			return server.Serve(lis)
		},
	}
}
