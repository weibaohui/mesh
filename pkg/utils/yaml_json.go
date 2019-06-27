package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/rand"
	"os"
	"sigs.k8s.io/yaml"
)

func YamlToFile(obj interface{}) {

	metadata, _ := meta.Accessor(obj)
	num := rand.IntnRange(1000, 9999)
	name := fmt.Sprintf("%s-%s-%d", metadata.GetNamespace(), metadata.GetName(), num)

	bytes, err := yaml.Marshal(obj)
	if err != nil {
		fmt.Println("输出YAML文件时出错 ", name, err.Error())
		fmt.Println(err)
	}

	err = ioutil.WriteFile("/Users/baohui/Desktop/mesh/"+name+".yaml", bytes, os.ModePerm)
	if err != nil {
		fmt.Println("输出YAML文件时出错 ", name, err.Error())
	}
}

func ObjToJson(obj interface{}) {

	bytes, err := json.Marshal(obj)
	if err != nil {
		logrus.WithField("obj", "true").Error(err.Error())

	}
	logrus.WithField("obj", "true").Error(string(bytes))

}
func YamlToJson(obj interface{}) {

	metadata, _ := meta.Accessor(obj)
	name := fmt.Sprintf("%s-%s", metadata.GetNamespace(), metadata.GetName())

	bytes, err := yaml.Marshal(obj)
	if err != nil {
		logrus.WithField("ns", name).Error(err.Error())

	}
	logrus.WithField("ns", name).Error(string(bytes))

}
