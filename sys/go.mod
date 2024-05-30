module github.com/hxnas/pkg/sys

go 1.22.3

replace github.com/hxnas/pkg/lod => ../lod/

require (
	github.com/hxnas/pkg/lod v0.0.0-00010101000000-000000000000
	github.com/moby/sys/mount v0.3.3
	github.com/moby/sys/mountinfo v0.7.1
	github.com/moby/sys/symlink v0.2.0
	golang.org/x/sync v0.7.0
	golang.org/x/sys v0.20.0
)
