package ipar

import (
	"context"

	"github.com/chigopher/pathlib"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	coreiface "github.com/ipfs/kubo/core/coreiface"
	coreoptions "github.com/ipfs/kubo/core/coreiface/options"
)

type IpfsAdder struct {
	archivalDir *pathlib.Path
	coreAPI     coreiface.CoreAPI
}

func NewIpfsAdder(archivalDir *pathlib.Path, coreAPI coreiface.CoreAPI) *IpfsAdder {
	return &IpfsAdder{
		archivalDir: archivalDir,
		coreAPI:     coreAPI,
	}
}

func (i *IpfsAdder) Add(ctx context.Context, p *pathlib.Path, pinName string) (*path.ImmutablePath, error) {
	s, err := p.Stat()
	if err != nil {
		return nil, err
	}
	n, err := files.NewSerialFile(p.String(), true, s)
	if err != nil {
		return nil, err
	}

	path, err := i.coreAPI.Unixfs().Add(
		ctx,
		n,
		coreoptions.Unixfs.Pin(true, pinName),
		coreoptions.Unixfs.Progress(true),
		coreoptions.Unixfs.FsCache(true),
		coreoptions.Unixfs.Nocopy(true),
		coreoptions.Unixfs.RawLeaves(true),
		coreoptions.Unixfs.IncludeEmptyDirs(true))
	if err != nil {
		return nil, err
	}
	return &path, nil
}

func (i *IpfsAdder) Remove(ctx context.Context, p *pathlib.Path, pinName string) error {
	s, err := p.Stat()
	if err != nil {
		return err
	}
	_, err = files.NewSerialFile(p.String(), true, s)
	if err != nil {
		return err
	}
	path, err := path.NewPath(path.IPNSNamespace + pinName)
	if err != nil {
		return err
	}

	err = i.coreAPI.Pin().Rm(
		ctx,
		path,
		coreoptions.Pin.RmRecursive(true))
	if err != nil {
		return err
	}
	return nil
}
