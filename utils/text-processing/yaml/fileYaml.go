package yaml

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

//ReadFromYAML reads the YAML file and pass to the object
//args:
//	path: file path location
//	target: object which will hold the value
//returns:
//	error: operation state error
func ReadFromYAML(path string, target interface{}) error {
	yf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	yf = []byte(os.ExpandEnv(string(yf)))
	return yaml.Unmarshal(yf, target)
}
