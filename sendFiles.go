package main

import (
	"fmt"
	"flag"
	"os"
	"encoding/json"
	"log"
	"path/filepath"
	"github.com/fatih/color"
	"github.com/scorredoira/email"
	"github.com/mattn/go-zglob"
	"io"
	"path"
	"time"
	"archive/zip"
	"strings"
	"net/mail"
	"net/smtp"
	"net"
	//	"strconv"
)


type Config struct {
	ApplicationLogs []ApplicationLog
	Global
}

type ApplicationLog struct {
	AppName       string `json:"appName"`
	LogPath       string `json:"logPath"`
	Pattern       string `json:"pattern"`

}

type Global struct {
	SmtpHostPort string `json:"smtpHostPort"`
	LogFile string `json:"logFile"`
	EmailFromName string `json:"emailFromName"`
	EmailFromAddr string `json:"emailFromAddr"`
	EmailTo string `json:"emailTo"`
	AdminEmail string `json:"adminEmail"`
	TmpDir string `json:"tmpDir"`
	TimeToSend int64 `json:"timeToSend"`
}

var command string

//var pinfo color.New(color.FgYellow,color.BgBlack)
//var panounce color.New(color.FgHiCyan,color.BgBlack)
//var perror color.New(color.FgHiRed,color.BgBlack)


func init() {

	flag.StringVar(&command, "command", command, "" +
		"Commands:\n " +
		"\t\t showConfig - Show configuration file \n" +
		"\t\t Run - Run Trigger \n")
}

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		//baseDir = filepath.Base(source)
		baseDir = ""
	}

	filepath.Walk(source, func(path1 string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path1, source))
		}

		if info.IsDir() {
			return nil
			//header.Name += "/"

		} else {
			header.Method = zip.Deflate

		}
		//fmt.Print("path: "+path1+"\n")
		//fmt.Print("header: "+header.Name+"\n")


		writer, err := archive.CreateHeader(header)

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path1)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err

	})

	return err
}

func CopyFile(src,dst string) {

	pinfo := color.New(color.FgYellow,color.BgBlack)
	//panounce := color.New(color.FgHiCyan,color.BgBlack)
	perror := color.New(color.FgHiRed,color.BgBlack)


	r, err := os.Open(src)
	if err != nil {
		perror.Printf("[ERROR] Open src %s %s\n",src,err)
		log.Fatalf("[ERROR] %s",err)

		//panic(err)
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		perror.Printf("[ERROR] Create dst %s %s\n",dst,err)
		log.Fatalf("[ERROR] %s",err)

		//panic(err)

	}
	defer w.Close()

	// do the actual work
	n,err := io.Copy(w, r)
	if err != nil {
		perror.Printf("[ERROR] Copy %s %s\n",src,err)
		log.Fatalf("[ERROR] Copy %s %s",src,err)

		//panic(err)
	} else {
		pinfo.Printf("[INFO] Copy file %s - %v bytes\n",src,n)
		log.Printf("[INFO] Copy file %s - %v bytes",src,n)
	}


	//fmt.Printf("Copied %v bytes\n",n)
	//log.Printf("Copied %v bytes",n)

	//return ret

}


func main() {
	flag.Parse()

	pinfo := color.New(color.FgYellow,color.BgBlack)
	panounce := color.New(color.FgHiCyan,color.BgBlack)
	//panounce := color.New(color.Bold, color.FgGreen).PrintlnFunc()
	//warning := color.New(color.FgYellow)
	perror := color.New(color.FgHiRed,color.BgBlack)
	appdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		perror.Printf("[ERROR] %v", err)
		log.Fatalf("[ERROR] %s",err)
	}

	//Парсим файл конфигурации
	file, _ := os.Open(appdir+"/config.json")
	//file, _ := os.Open("c:/goSendFiles/config.json")
	decoder := json.NewDecoder(file)
	config := new(Config)
	err = decoder.Decode(&config)
	
	if err != nil {
		perror.Printf("[ERROR] Error read configuration file(config.json) %v\n", err)
		//perror.Printf("[ERROR] error opening log file: %v\n", appdir+"/config.json")
		log.Fatalf("[ERROR] Error read configuration file(config.json) %v\n", err)

	}

	defer file.Close()
	pinfo.Printf("[INFO] application dir %v\n", appdir)
	//lll, err := os.Create(config.LogFile+"sfsdfsdf")
	//defer lll.Close()
	f, err1 := os.OpenFile(config.LogFile,  os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if err1 != nil {
		perror.Printf("[ERROR] error opening log file: %v\n", err1)
		perror.Printf("[ERROR] error opening log file: %v\n", config.LogFile)
		perror.Printf("[ERROR] error opening log file: %v\n", appdir+"/config.json")
		log.Fatalf("[ERROR] error opening log file: %v\n", err1)
	}
	defer f.Close()
	log.SetOutput(f)


	switch {
	default:
		//log.Fatal("Invalid or undefined command, type -h to help \n")
		perror.Printf("%s", "[ERROR] Invalid or undefined command, type -h to help \n")
		log.Fatalf("%s", "[ERROR] Invalid or undefined command, type -h to help \n")

	case command == "showConfig":
		color.Green("Application Logs")
		for _, ALogs := range config.ApplicationLogs {
			color.Cyan("appName: " + ALogs.AppName)
			fmt.Printf("logPath: %s\n", ALogs.LogPath)
			fmt.Printf("pattern: %s\n", ALogs.Pattern)
			fmt.Printf("%s", "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
		}
		color.Green("Global config")
		fmt.Printf("smtpHostPort: %s\n", config.SmtpHostPort)
		fmt.Printf("emailFromName: %s\n", config.EmailFromName)
		fmt.Printf("emailFromAddr: %s\n", config.EmailFromAddr)
		fmt.Printf("emailTo: %s\n", config.EmailTo)
		fmt.Printf("adminEmail: %s\n", config.AdminEmail)
		fmt.Printf("logFile: %s\n", config.LogFile)
		fmt.Printf("tmpDir: %s\n", config.TmpDir)
		fmt.Printf("timeToSend: %v\n", config.TimeToSend)

	case command == "Run":
		//color.Red("Application Logs")
		ntoday := time.Now()
		today := ntoday.Format("2006-01-02")
		//tnow := ntoday.Format("20060102150405")

		for _, ALogs := range config.ApplicationLogs {
			panounce.Printf("[INFO] Application Name: %s\n", ALogs.AppName)
			mailbody := "-="+ALogs.AppName+"=-\n"
			//fmt.Printf("logPath: %s\n", ALogs.LogPath)

			matches, err := zglob.Glob(ALogs.LogPath + `/`+ ALogs.Pattern)
			if err != nil {
				perror.Printf("[ERROR] Search logs: %s",err)
				log.Fatalf("[ERROR] Search logs: %s",err)
				os.Exit(1)
			}

			for _, match := range matches {
				//fmt.Println(match)
				file, err := os.Stat(match)
				if err != nil {
					perror.Printf("[ERROR] %s",err)
					log.Fatalf("[ERROR] %s",err)
				}

				modifiedtime := file.ModTime().Unix()
				//fmt.Println("---Last modified time : ", modifiedtime)
				nowtime := time.Now().Unix()
				//fmt.Println("===Now time : ", nowtime)

				if (nowtime-modifiedtime) <= config.TimeToSend {
					//info.Printf("File match %s\n",match)
					os.MkdirAll(config.TmpDir+"/"+ALogs.AppName,0666)
					//pinfo.Printf("[INFO] Copy file %s\n",path.Base(match))
					mailbody = mailbody+path.Base(match)+"\n"
					CopyFile(match,config.TmpDir+"/"+ALogs.AppName+"/"+path.Base(match))
					//log.Printf("[INFO] Copy file %s",match)
					os.Chtimes(config.TmpDir+"/"+ALogs.AppName+"/"+path.Base(match),file.ModTime(),file.ModTime())
				}
			}
			//fmt.Println(today)
			// archive files
			os.MkdirAll(config.TmpDir+"/output",0666)
			err = zipit(config.TmpDir+"/"+ALogs.AppName+"/", config.TmpDir+"/output/"+ALogs.AppName+"_"+today+".zip")
			if err != nil {
				perror.Printf("[ERROR] Packing: %s\n",err)
				log.Fatalf("[ERROR] Packing: %s\n",err)
			} else {
				pinfo.Printf("[INFO] Packing %s \n",config.TmpDir+"/"+ALogs.AppName+"/")
				log.Printf("[INFO] Packing %s \n",config.TmpDir+"/"+ALogs.AppName+"/")
			}
			//delete tmpfiles
			err = os.RemoveAll(config.TmpDir+"/"+ALogs.AppName+"/")
			if err != nil {
				perror.Printf("[ERROR] Remove tmpdirs: %s\n",err)
				log.Fatalf("[ERROR] Remove tmpdirs: %s\n",err)
			} else {
				pinfo.Printf("[INFO] Deleted tmp %s \n",config.TmpDir+"/"+ALogs.AppName+"/")
				log.Printf("[INFO] Deleted tmp %s \n",config.TmpDir+"/"+ALogs.AppName+"/")
			}


			//--SEND EMAIL
			m := email.NewMessage("Logs "+ALogs.AppName+"_"+today, mailbody)
			//mailFrom := mail.Address{Name: config.EmailFromName, Address: config.EmailFromAddr}
			mailFrom := mail.Address{Address: config.EmailFromAddr}
			m.From = mailFrom

			//m.To = []string{config.EmailTo}
			m.To = strings.Split(config.EmailTo,",")
			attach := config.TmpDir+"/output/"+ALogs.AppName+"_"+today+".zip"
			if err := m.Attach(attach); err != nil {
				perror.Printf("[ERROR] Attach file: %s\n",err)
				log.Fatalf("[ERROR] Attach file: %s\n",err)
			}

			//host, _, _ := net.SplitHostPort(config.SmtpHostPort)
			//auth := smtp.PlainAuth("", "", "", host)
			//if err := email.Send(config.SmtpHostPort, auth, m); err != nil {
			if err := email.Send(config.SmtpHostPort, nil, m); err != nil {
				perror.Printf("[ERROR] Send mail: %s\n",err)
				log.Fatalf("[ERROR] Send mail: %s\n",err)
			} else {
				pinfo.Printf("[INFO] Sended %s to %s\n",path.Base(attach),config.EmailTo)
				log.Printf("[INFO] Sended %s to %s\n",path.Base(attach),config.EmailTo)
			}

			//delete archives
			err = os.RemoveAll(config.TmpDir+"/output/"+ALogs.AppName+"_"+today+".zip")
			if err != nil {
				perror.Printf("[ERROR] Remove tmp archive: %s\n",err)
				log.Fatalf("[ERROR] Remove tmp archive: %s\n",err)
			} else {
				pinfo.Printf("[INFO] Remove tmp archive %s\n",config.TmpDir+"/output/"+ALogs.AppName+"_"+today+".zip")
				log.Printf("[INFO] Remove tmp archive %s\n",config.TmpDir+"/output/"+ALogs.AppName+"_"+today+".zip")
			}

		}
		//rotate sendlog logfile
		ntoday = time.Now()
		tnow := ntoday.Format("20060102150405")
		today = ntoday.Format("2006-01-02")

		err = zipit(config.LogFile, config.LogFile+"_"+tnow+".zip")
		if err != nil {
			perror.Printf("[ERROR] Rotate log: %s\n",err)
			log.Fatalf("[ERROR] Rotate log: %s\n",err)
		}

		f.Close()
		err = os.Remove(config.LogFile)
		if err != nil {
			perror.Printf("[ERROR] Remove oldlog:  %s\n",err)
			//log.Fatalf("[ERROR] Remove oldlog: %s\n",err)
		}

		//--SEND EMAIL
		mailbody := "SendLogs LogFile"
		m := email.NewMessage("SendLogs LogFile "+today, mailbody)
		mailFrom := mail.Address{Name: config.EmailFromName, Address: config.EmailFromAddr}
		m.From = mailFrom

		//m.To = []string{config.EmailTo}
		m.To = strings.Split(config.AdminEmail,",")
		attach := config.LogFile+"_"+tnow+".zip"
		if err := m.Attach(attach); err != nil {
			perror.Printf("[ERROR] Attach file: %s\n",err)
			//log.Fatalf("[ERROR] Attach file: %s\n",err)
		}
		//host, _, _ := net.SplitHostPort(config.SmtpHostPort)
		//auth := smtp.PlainAuth("", "", "", host)
		//if err := email.Send(config.SmtpHostPort, auth, m); err != nil {
		if err := email.Send(config.SmtpHostPort, nil, m); err != nil {
			perror.Printf("[ERROR] Send mail: %s\n",err)
			//log.Fatalf("[ERROR] Send mail: %s\n",err)
		} else {
			pinfo.Printf("[INFO] Sended %s to %s\n",path.Base(attach),config.EmailTo)
			///log.Printf("[INFO] Sended %s to %s\n",path.Base(attach),config.EmailTo)
		}
	}





}
