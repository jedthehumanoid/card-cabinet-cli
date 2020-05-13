module card-cabinet-cli

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/jedthehumanoid/card-cabinet v0.0.0-20200424050516-14e85b67161c
	golang.org/x/text v0.3.2
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/jedthehumanoid/card-cabinet => ../card-cabinet
