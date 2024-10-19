# Golang YAML Settings

This is my Golang settings module.

You have the following 3 options to read the config file:
- ./settings.yaml
- /usr/local/myapp/etc/settings.yaml
- /config/settings.yaml

To load the settings, execute the following in the func init():
```go
func init() {
    var Load Config
    settings.readConfig(&Load)
}
```

A structure like the following is also required:
```go
type Config struct {
    Remote struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
        User string `yaml:"user"`
        Pass string `yaml:"pass"`
        Path struct {
            Movies string `yaml:"movies"`
            Shows  string `yaml:"shows"`
            Animes string `yaml:"animes"`
        } `yaml:"path"`
        Reload int   `yaml:"reload"`
    } `yaml:"remote"`
}
```
