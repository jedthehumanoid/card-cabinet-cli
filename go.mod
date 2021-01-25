module card-cabinet-cli

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/JohannesKaufmann/html-to-markdown v1.2.0
	github.com/jedthehumanoid/cardcabinet v0.0.0-20210125195512-0df2e8317189
	github.com/spf13/cobra v1.1.1
	golang.org/x/text v0.3.2
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/jedthehumanoid/cardcabinet => ../cardcabinet
