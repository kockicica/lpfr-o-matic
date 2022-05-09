package watchdog

import (
	"log"
	"os"
	"os/exec"
	"time"

	"lpfr-o-matic/pkg/lpfr"
	"lpfr-o-matic/pkg/telegram"
)

type WatchdogConfig struct {
	ExePath           string
	CheckUrl          string
	CheckInterval     int `mapstructure:"interval"`
	Pin               string
	NoPin             bool
	MiddlewareApp     string `mapstructure:"middleware"`
	Telegram          bool
	TelegramChannelId string `mapstructure:"telegram-channel-id"`
	TelegramApiKey    string `mapstructure:"telegram-api-key"`
	TelegramSender    string `mapstructure:"telegram-sender"`
}

type Watchdog struct {
	exePath              string
	checkUrl             string
	checkInterval        time.Duration
	exeCmd               *exec.Cmd
	stop                 chan os.Signal
	stopped              chan interface{}
	client               *lpfr.LPFRClient
	pin                  string
	noPin                bool
	middlewareApp        string
	hasMiddlewareStarted bool
	sendTelegramMessages bool
	telegramClient       *telegram.Client
}

func (w *Watchdog) Start() error {

	go func() {

		// initial check if running
		lpfrStatus := w.checkStatus()
		log.Printf("Initial LPFR status: %s", lpfrStatus)
		if err := w.sendPin(lpfrStatus); err != nil {
			_ = w.sendTelegramMessage(telegram.Message{Title: "Fatal", Message: "Unable to set PIN after one attempt"})
			log.Fatalln("Error sending pin:", err)
		}

		for true {
			select {
			case <-w.stop:
				// close
				log.Println("Stopping watchdog")
				err := w.stopExe()
				if err != nil {
					log.Print(err)
				}
				w.stopped <- true
				break
			case <-time.After(time.Second * w.checkInterval):
				// every ten seconds or so
				log.Println("check interval exceeded, do check")
				lpfrStatus = w.checkStatus()
				log.Printf("LPFR status: %s", lpfrStatus)
				if err := w.sendPin(lpfrStatus); err != nil {
					_ = w.sendTelegramMessage(telegram.Message{Title: "Fatal", Message: "Unable to set PIN after one attempt"})
					log.Fatalln("Error sending pin:", err)
				}
				if lpfrStatus == lpfr.Ready && !w.hasMiddlewareStarted && w.middlewareApp != "" {
					log.Println("Trying to start middleware app:", w.middlewareApp)
					// start middleware app
					w.hasMiddlewareStarted = true
					mwcmd := exec.Command("cmd.exe", "/c", "start", w.middlewareApp, "/min")
					err := mwcmd.Start()
					if err != nil {
						log.Println("Error starting middleware app:", err)
					}
				}
			}
		}
	}()

	return nil
}

func (w *Watchdog) Stop() error {
	w.stop <- os.Kill
	<-w.stopped
	return nil
}

func (w *Watchdog) checkStatus() lpfr.LPFRStatus {
	var lpfrStatus lpfr.LPFRStatus
	var err error
	lpfrStatus, err = w.client.EnvironmentStatus()
	if err != nil {
		log.Println("LPFR seems not to be running, trying to start it")
		err := w.runExe()
		if err != nil {
			_ = w.sendTelegramMessage(telegram.Message{Title: "Fatal", Message: "Unable to start lpfr process"})
			log.Fatal(err)
		}
		count := 5
		lpfrStatus, err = w.client.EnvironmentStatus()
		if err != nil {
			for i := 0; i < count; i++ {
				log.Println("Error getting status, wait for 5 seconds than check status again")
				<-time.After(time.Second * 5)
				lpfrStatus, err = w.client.EnvironmentStatus()
				if err == nil {
					_ = w.sendTelegramMessage(telegram.Message{Title: "Info", Message: "LPFR process started"})
					return lpfrStatus
				}
			}
			_ = w.sendTelegramMessage(telegram.Message{Title: "Fatal", Message: "Unable to get lpfr environment status"})
			log.Fatal(err)
		}
		_ = w.sendTelegramMessage(telegram.Message{Title: "Info", Message: "LPFR process started"})
		return lpfrStatus
	}
	return lpfrStatus

}

func (w *Watchdog) runExe() error {
	log.Println("Trying to start external process: ", w.exePath)
	w.exeCmd = exec.Command("cmd.exe", "/C", "start", w.exePath, "/min")
	err := w.exeCmd.Start()
	return err
}

func (w *Watchdog) stopExe() error {
	if w.exeCmd != nil && w.exeCmd.Process != nil {
		err := w.exeCmd.Process.Signal(os.Kill)
		w.exeCmd = nil
		return err
	}
	return nil
}

func (w *Watchdog) sendPin(status lpfr.LPFRStatus) error {
	var err error
	if status != lpfr.NeedsPIN {
		return nil
	}
	if w.noPin {
		log.Println("Skip automatic pin setting")
		return nil
	}
	log.Printf("Wait for 2 seconds before pin attempt")
	<-time.After(2 * time.Second)

	err = w.client.SendPIN(w.pin)
	if err != nil {
		log.Println(err)
	} else {
		_ = w.sendTelegramMessage(telegram.Message{Title: "Info", Message: "PIN successfully set"})
		log.Println("PIN set successfully")
	}
	return err
}

func (w *Watchdog) sendTelegramMessage(message telegram.Message) error {
	if w.sendTelegramMessages {
		return w.telegramClient.SendMessage(message)
	}
	return nil
}

func NewWatchdog(config WatchdogConfig) *Watchdog {
	dg := new(Watchdog)
	dg.exePath = config.ExePath
	dg.checkUrl = config.CheckUrl
	dg.checkInterval = time.Duration(config.CheckInterval)
	dg.stop = make(chan os.Signal)
	dg.stopped = make(chan interface{})
	dg.client = lpfr.NewLPFRClient(config.CheckUrl)
	dg.pin = config.Pin
	dg.noPin = config.NoPin
	dg.middlewareApp = config.MiddlewareApp
	dg.hasMiddlewareStarted = false
	dg.sendTelegramMessages = config.Telegram
	if dg.sendTelegramMessages {
		dg.telegramClient = telegram.NewClient(config.TelegramApiKey, config.TelegramChannelId, config.TelegramSender)
	}
	return dg
}
