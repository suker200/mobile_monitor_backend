package main

import (
	"encoding/json"
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-fsnotify/fsnotify"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var Config ConfigInfo
var Data DataReport
var Data2 DataReport
var Update bool
var C = &cron.Cron{}

func readConfig() {
	data, err := ioutil.ReadFile("./config/config.yaml")
	if err == nil {
		err = yaml.Unmarshal(data, &Config)
		if err != nil {
			log.Fatalf("error: %v", err)
		} else {
			fmt.Println(Config)
		}
	} else {
		log.Fatalf("error: %v", err)
	}
}

func Auth(username, password string) bool {
	Url := os.Getenv("SSO_URL")
	client := &http.Client{}

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("group", "monitor_app_user")

	resp, err := client.PostForm(Url, data)

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			return true
		}

		return false
	}
}

func Get_PrometheusMetrics(url string) Resp_PrometheusMetrics {
	var metrics Resp_PrometheusMetrics
	client := &http.Client{}

	resp, err := client.Get(url)

	if err != nil {
		return metrics
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		// fmt.Println(string(body))
		if err != nil {
			fmt.Println(err.Error())
			return metrics
		}
		if resp.StatusCode != 200 {
			fmt.Println("Alarm via telegram directly to devops, info: " + string(body))
			// return string(body), errors.New(string(body))
			return metrics
		}

		if err := json.Unmarshal(body, &metrics); err != nil {
			// AlarmMe(serviceName, nameSpace, "Failed", err.Error())
			return metrics
		}
		return metrics
	}
}

func UpdateDataReport(_data MetricReport) {
	Update = true
	time.Sleep(2)
	// Data.Lock()
	if len(Data.Infos) == 0 {
		if len(_data.Issue) != 0 {
			_data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
		} else {
			_data.Alias = Config.Metrics[_data.Title].Alias
		}

		Data.Infos = append(Data.Infos, _data)
	} else {
		check := false // Case this metric is not exist in array
		for num, v := range Data.Infos {
			if v.Title == _data.Title {
				if len(_data.Issue) != 0 {
					_data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
				} else {
					_data.Alias = Config.Metrics[_data.Title].Alias
				}
				// _data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
				// a = append(a[:i], a[i+1:]...)
				Data.Infos = append(Data.Infos[:num], Data.Infos[num+1:]...)
				Data.Infos = append(Data.Infos, _data)
				check = true // metrics exits in array
				break
			}
		}
		if !check { // metrics not exits in array
			if len(_data.Issue) != 0 {
				_data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
			} else {
				_data.Alias = Config.Metrics[_data.Title].Alias
			}
			// _data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
			Data.Infos = append(Data.Infos, _data)
		}
	}

	// Data.Unlock()
	Update = false
	// time.Sleep(2)
	// Data2.Lock()
	if len(Data2.Infos) == 0 {
		if len(_data.Issue) != 0 {
			_data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
		} else {
			_data.Alias = Config.Metrics[_data.Title].Alias
		}
		// _data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
		Data2.Infos = append(Data2.Infos, _data)
	} else {
		check := false // Case this metric is not exist in array
		for num, v := range Data2.Infos {
			if v.Title == _data.Title {
				// a = append(a[:i], a[i+1:]...)
				if len(_data.Issue) != 0 {
					_data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
				} else {
					_data.Alias = Config.Metrics[_data.Title].Alias
				}
				// _data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
				Data2.Infos = append(Data2.Infos[:num], Data2.Infos[num+1:]...)
				Data2.Infos = append(Data2.Infos, _data)
				check = true // metrics exits in array
				break
			}
		}
		if !check { // metrics not exits in array
			if len(_data.Issue) != 0 {
				_data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
			} else {
				_data.Alias = Config.Metrics[_data.Title].Alias
			}
			// _data.Alias = Config.Metrics[_data.Title].Alias + " (" + strconv.Itoa(len(_data.Issue)) + ")"
			Data2.Infos = append(Data2.Infos, _data)
		}
	}

	// Data2.Unlock()
}

func ScheduleUpdate() {
	// C.Stop()
	// C = cron.New()
	// for _, v := range Config.Metrics {
	// 	fmt.Println("Add cron " + v.Name)
	// 	C.AddFunc("@every 1m", Disk_Usage)
	// }
	for {
		Disk_Usage()
		// fmt.Println(Data)
		Mem_Usage()
		// fmt.Println("============")
		// fmt.Println(Data)
		ServerState()
		PrometheusThrotle()
		Load_Usage()
		PodRestart()
		PodPending()
		UpstreamResponseTime()
		time.Sleep(time.Duration(60) * time.Second)
	}
}

func ConfigUpdate() {
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)
				readConfig()
				// ScheduleUpdate()

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	if err := watcher.Add("./config/config.yaml"); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}

func App(c *gin.Context) {
	// c.BindJSON(&args)

	if !Update {
		c.JSON(200, Data)
	} else {
		c.JSON(200, Data2)
	}
}

func Ping(c *gin.Context) {
	c.JSON(200, "okie")
}

func main() {
	C = cron.New()
	C.Start()
	readConfig()
	go ScheduleUpdate()
	go ConfigUpdate()

	r := gin.Default()

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if ok := Auth(userId, password); ok {
				return userId, true
			}

			return userId, false
		},
		// Authorizator: func(userId string, c *gin.Context) bool {
		// 	if userId == "admin" {
		// 		return true
		// 	}

		// 	return false
		// },
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	r.GET("/ping", Ping)
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/v1/app", App)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	s := &http.Server{
		Addr:           ":12345",
		Handler:        r,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
