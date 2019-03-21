package idip

import (
	"encoding/xml"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"github.com/clbanning/mxj"
)

func Generate(filename string) {

	logrus.WithFields(logrus.Fields{
		"filename": filename,
	}).Info("Generate")

	xmlValue, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Error(err)
	}

	var v map[string]interface{}
	dec := xml.Unmarshal(xmlValue, &v)

	logrus.WithFields(logrus.Fields{
		"dec": dec,
	}).Info("xml.Unmarshal")

	type Map map[string]interface{}
	//mv := Map(v)
	mv, err := mxj.NewMapXml(xmlValue) // unmarshal
	if err != nil {
		logrus.Error(err)
	}

	//logrus.WithFields(logrus.Fields{
	//	"mv": mv,
	//}).Info("xml.Unmarshal")

	for k, v := range mv["metalib"].(map[string]interface{}) {

		logrus.Warnf("key: %v, value: %v", k, v)
	}

	//dec := xml.NewDecoder(bytes.NewBuffer(xmlData))
	//for {
	//	token, err := dec.Token()
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//		panic(err)
	//	}
	//	switch element := token.(type) {
	//	case xml.StartElement:
	//
	//		logrus.WithFields(logrus.Fields{
	//			"element": element,
	//		}).Info("StartElement")
	//
	//	case xml.EndElement:
	//
	//		logrus.WithFields(logrus.Fields{
	//			"element": element,
	//		}).Info("EndElement")
	//	case xml.CharData:
	//
	//		logrus.WithFields(logrus.Fields{
	//			"element": element,
	//		}).Info("CharData")
	//	case xml.Comment:
	//
	//		logrus.WithFields(logrus.Fields{
	//			"element": element,
	//		}).Info("Comment")
	//	case xml.Directive:
	//
	//		logrus.WithFields(logrus.Fields{
	//			"element": element,
	//		}).Info("Directive")
	//	case xml.ProcInst:
	//
	//		logrus.WithFields(logrus.Fields{
	//			"element": element,
	//		}).Info("ProcInst")
	//	}
	//}
}
