module github.com/regel/wsb

go 1.20

replace github.com/regel/wsb/pkg/finance/types => ./pkg/finance/types

replace github.com/regel/wsb/pkg/finance/yahoo => ./pkg/finance/yahoo

replace github.com/regel/wsb/pkg/common => ./pkg/common

replace github.com/regel/wsb/pkg/finance => ./pkg/finance

replace github.com/regel/wsb/pkg/config => ./pkg/config

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20210324051636-2c4c8ecb7826
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20210315160823-c6e025ad8005 // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
