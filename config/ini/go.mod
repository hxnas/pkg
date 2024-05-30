module github.com/cnk3x/pkg/config/ini

go 1.22.3

require (
	github.com/cnk3x/pkg/config v0.0.0-00010101000000-000000000000
	gopkg.in/ini.v1 v1.67.0
)

require github.com/stretchr/testify v1.9.0 // indirect

replace github.com/cnk3x/pkg/config => ../
