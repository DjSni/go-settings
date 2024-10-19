/*
	This module loads a yaml settings file and packs it into a struct.

To load the settings, execute the following in the func init():

func init() {
    var Load Config
    settings.readConfig(&Load)
}

A struct as follows is also required:
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

You can read everything further under "https://github.com/DjSni/go-log"
*/
package settings

import (
    "reflect"
	"errors"
	"os"

	log "github.com/DjSni/go-log"
	yaml "gopkg.in/yaml.v3"
)

var (
	/* NAME var is needed for the CONFIG_PATH_DEB var.
    Is the name of the go App for the path to get the correct config file
    -> "/usr/local/etc/" + NAME + "/settings.yaml"
	*/
	NAME string
	// local config path
	CONFIG_PATH_LOCAL = "./settings.yaml"
	/* config path 
    the NAME var is needed
    -> settings.NAME = MyApp
    */
	CONFIG_PATH_DEB = "/usr/local/etc/" + NAME + "/settings.yaml"
	// config path for docker installations
	CONFIG_PATH_DOCKER = "/config/settings.yaml"
)

// validateConfig überprüft, ob alle Felder in der Struktur gesetzt sind
func validateConfig(config interface{}) error {
	v := reflect.ValueOf(config)

	// Sicherstellen, dass wir einen Zeiger haben und die Struktur referenziert wird
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("config muss ein nicht-nil Zeiger auf eine Struktur sein")
	}

	// Auf die Struktur zugreifen
	v = v.Elem()

	// Gehe durch jedes Feld und prüfe, ob es den Zero-Wert hat (nicht gesetzt)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Überprüfe verschachtelte Strukturen
		if field.Kind() == reflect.Struct {
			err := validateConfig(field.Addr().Interface())
			if err != nil {
				return err
			}
		} else {
			// Prüfe, ob das Feld den Zero-Wert hat
			if isZero(field) {
				return errors.New("Feld " + fieldType.Name + " in der Konfiguration darf nicht leer sein")
			}
		}
	}

	return nil
}

// isZero prüft, ob ein Feld den Zero-Wert hat (bedeutet, dass es nicht gesetzt wurde)
func isZero(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

// test if config files exist and retun the path or exit with error
func testConfigPath() (configPath string) {
	if _, err := os.Stat(CONFIG_PATH_LOCAL); err == nil {
		configPath = CONFIG_PATH_LOCAL
	} else if _, err := os.Stat(CONFIG_PATH_DEB); err == nil {
		configPath = CONFIG_PATH_DEB
	} else if _, err := os.Stat(CONFIG_PATH_DOCKER); err == nil {
		configPath = CONFIG_PATH_DOCKER
	} else {
		log.Error("Please provide a config file")
		log.Error("The paths of the config file can be:")
		log.Error(" ->", CONFIG_PATH_LOCAL)
		log.Error(" ->", CONFIG_PATH_DEB)
		log.Error(" ->", CONFIG_PATH_DOCKER)
		log.Fatal("No config file found!")
	}
    return configPath
}

func readConfig(theConfig interface{}) {
	// Open YAML file
	file, err := os.Open(testConfigPath())
	if err != nil {
		log.Error(err.Error())
	}
	defer file.Close()

	// Decode YAML file to struct
	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(theConfig); err != nil {
			log.Error(err.Error())
		}
	}

	// Validierung der Struktur nach dem Unmarshalling
	err = validateConfig(theConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
}


//func main() {
//    var Load Config
//    readConfig(&Load)

//	log.Info("Type of variable1:",reflect.TypeOf(Load.Remote.Host))
    
    //LogConfig( Load )
//}