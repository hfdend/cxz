package wxpay

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestXMLUnmarshalXML(t *testing.T) {
	xmlString := `<xml>
    <appid><![CDATA[wx6ab211366652b877]]></appid>
</xml>`
	var values DataValues = DataValues{}
	err := xml.Unmarshal([]byte(xmlString), &values)
	fmt.Println(err)
	fmt.Println(values)
}
