package lpfr

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type LPFRStatus int

const (
	Ready LPFRStatus = iota
	Stopped
	NeedsSmartCard
	Initializing
	NeedsPIN
	Unknown
)

func (s LPFRStatus) String() string {
	switch s {
	case Ready:
		return "Ready"
	case Stopped:
		return "Stopped"
	case NeedsSmartCard:
		return "Needs smart card"
	case Initializing:
		return "Initializing"
	case NeedsPIN:
		return "Needs PIN"
	case Unknown:
		return "Unknown"
	}
	return "Some other status"
}

type LPFRClient struct {
	checkUrl string
}

//func (c *LPFRClient) GetStatus() (LPFRStatus, error) {
//	fullUrl := fmt.Sprintf("%s/api/v3/status", c.checkUrl)
//	log.Println("Get status, url: ", fullUrl)
//	rsp, err := http.Get(fullUrl)
//	if err != nil {
//		return Stopped, err
//	}
//	data, err := io.ReadAll(rsp.Body)
//	if err != nil {
//		return Stopped, err
//	}
//	status := StatusResponse{}
//	err = json.Unmarshal(data, &status)
//	if err != nil {
//		return Stopped, err
//	}
//	switch status.IsPinRequired {
//	case true:
//		return NeedsPIN, nil
//	case false:
//		return Ready, nil
//	default:
//		return Unknown, nil
//	}
//}

func (c *LPFRClient) EnvironmentStatus() (LPFRStatus, error) {
	fullUrl := fmt.Sprintf("%s/api/v3/environment-parameters", c.checkUrl)
	log.Println("Get environment, url: ", fullUrl)
	rsp, err := http.Get(fullUrl)
	if err != nil {
		return Stopped, err
	}
	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Fatalln("Error reading environment response:", err)
	}
	strResp := string(data)
	strResp = strings.Trim(strResp, "\"\n")
	switch strResp {
	case "1300":
		return NeedsSmartCard, nil
	case "2400":
		return Initializing, nil
	case "1500":
		return NeedsPIN, nil
	case "1999":
		return Unknown, nil
	}

	return Ready, nil
}

//func (c *LPFRClient) IsAlive() bool {
//	fullUrl := fmt.Sprintf("%s/api/v3/attention", c.checkUrl)
//	log.Println("Check if lpfr is alive, url: ", fullUrl)
//	_, err := http.Get(fullUrl)
//	if err != nil {
//		return false
//	}
//	return true
//}

func (l *LPFRClient) SendPIN(pin string) error {
	fullUrl := fmt.Sprintf("%s/api/v3/pin", l.checkUrl)
	log.Println("Set PIN, url: ", fullUrl)
	rsp, err := http.Post(fullUrl, "application/json", bytes.NewBufferString(pin))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	codeResponse := string(data)
	codeResponse = strings.Trim(codeResponse, "\n")
	codeResponse = strings.Trim(codeResponse, "\"")
	log.Println("PIN set response: ", codeResponse)
	switch codeResponse {
	case "0100":
		return nil
	case "2100":
		return fmt.Errorf("wrong pin specified")
	case "2110":
		return fmt.Errorf("number of allowed PIN entries exceeded")
	case "1300":
		return fmt.Errorf("smart card is not inserted or something wrong")
	default:
		return fmt.Errorf("some other error")
	}
}

func NewLPFRClient(checkUrl string) *LPFRClient {
	return &LPFRClient{
		checkUrl: checkUrl,
	}
}
