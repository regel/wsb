module github.com/regel/tinkerbell

go 1.15

replace github.com/regel/tinkerbell/pkg/finance => ./pkg/finance

replace github.com/regel/tinkerbell/pkg/config => ./pkg/config

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
