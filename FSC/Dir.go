package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"context"
	"log"
	"os"
	"sync/atomic"
	"syscall"
	"time"
)

type Dir struct {
	Type      fuse.DirentType
	Attribute fuse.Attr
	Entries   map[string]interface{}
}

func (d *Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	//TODO implement me
	*attr = d.Attribute
	log.Println("Attr permission:", attr.Mode)
	log.Println("Attr:Modified At", attr.Mtime.String())
	return nil
}

func NewDir() *Dir {
	log.Println("NewDir called")
	atomic.AddUint64(&InodeCount, 1)
	return &Dir{
		Type: fuse.DT_Dir,
		Attribute: fuse.Attr{
			Inode: InodeCount,
			Atime: time.Now(),
			Mtime: time.Now(),
			Ctime: time.Now(),
			Mode:  os.ModeDir | 0o777,
		},
		Entries: map[string]interface{}{},
	}
}
func (d *Dir) LookUp(ctx context.Context, name string) (fs.Node, error) {
	node, ok := d.Entries[name]
	if ok {
		return node.(fs.Node), nil
	}
	return nil, syscall.ENOENT
}
func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateRequest) (fs.Node, fs.Handle, error) {
	log.Println("Create called with filename", req.Name)
	f := NewFile(nil)
	log.Println("Create:Modified at", f.Attributes.Mtime.String())
	d.Entries[req.Name] = f
	return f, f, nil
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	log.Println("Mkdir call with name:", req.Name)
	dir := NewDir()
	d.Entries[req.Name] = dir
	return dir, nil
}

func (d *Dir) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	if req.Valid.Atime() {
		d.Attribute.Atime = req.Atime
	}
	if req.Valid.Mtime() {
		d.Attribute.Mtime = req.Mtime
	}
	if req.Valid.Size() {
		d.Attribute.Size = req.Size
	}
	log.Println("Setter called:Attribute ", d.Attribute)
	return nil
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Println("ReadDirAll called")
	var entries []fuse.Dirent
	for k, v := range d.Entries {
		var a fuse.Attr
		v.(fs.Node).Attr(ctx, &a)
		entries = append(entries, fuse.Dirent{
			Inode: a.Inode,
			Type:  v.(EntryGetter).GetDirentType(),
			Name:  k,
		})
	}
	return entries, nil

}
