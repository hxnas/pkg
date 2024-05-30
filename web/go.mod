module github.com/hxnas/pkg/web

go 1.22.3

replace github.com/hxnas/pkg/lod => ../lod/

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/hxnas/pkg/lod v0.0.0-00010101000000-000000000000
)
