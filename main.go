package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
)

type User struct {
	GamerId          int    `fake:"{number:7}" json:"gamerId"`
	Username         string `fake:"{username}" json:"userName"`
	FullName         string `fake:"{name}" json:"fullName"`
	Birthday         string `json:"birthday"`
	IdentityNumber   string `fake:"{number:12}" json:"identityNumber"`
	PhoneNumber      string `fake:"{phone}" json:"phoneNumber"`
	Email            string `fake:"{email}" json:"email"`
	Address          string `fake:"{street}, {city}" json:"address"`
	RegistrationTime string `json:"registrationTime"`
	DeviceRegister   string `json:"deviceRegister"`
	Ip               string `fake:"{ipv4address}" json:"ip"`
	AgreeTerm        string `json:"agreeTerm"`
	OtpTime          string `json:"otpTime"`
}

type LoginUser struct {
	Username      string `fake:"{username}" json:"userName"`
	GamerId       int    `fake:"{number:7}" json:"gamerId"`
	GameName      string `json:"gameName"`
	PublisherName string `json:"publisherName"`
	Ip            string `fake:"{ipv4address}" json:"ip"`
	DeviceId      string `fake:"{number:12}" json:"deviceId"`
	DeviceName    string `json:"deviceName"`
	Os            string `json:"os"`
	LoginTime     string `json:"loginTime"`
}

type GameInfo struct {
	GameName      string
	PublisherName string
}

var Platforms = [3]string{"Android", "iOS", "Web"}
var AndroidDevices = [5]string{"Galaxy", "Oppo", "Vivo", "Pixel", "Redmi"}
var IOSDevices = [2]string{"iPhone", "iPad"}
var WebDevices = [3]string{"Chrome", "Firefox", "Safari"}
var GameNames = [5]GameInfo{
	{GameName: "PUBG", PublisherName: "Tencent"},
	{GameName: "Free Fire", PublisherName: "Garena"},
	{GameName: "Call of Duty", PublisherName: "Activision"},
	{GameName: "Mobile Legends", PublisherName: "Moonton"},
	{GameName: "Among Us", PublisherName: "InnerSloth"},
}

func getAndroidDevice(username string) string {
	return username + " " + AndroidDevices[gofakeit.Number(0, 4)]
}

func getIOSDevice(username string) string {
	return username + " " + IOSDevices[gofakeit.Number(0, 1)]
}

func getWebDevice() string {
	return WebDevices[gofakeit.Number(0, 2)]
}

func getUserDevice(username string) (string, string) {
	platform := Platforms[gofakeit.Number(0, 2)]

	switch platform {
	case "Android":
		return getAndroidDevice(username), platform
	case "iOS":
		return getIOSDevice(username), platform
	case "Web":
		return getWebDevice(), platform
	}
	return "", platform
}

func getSyncUsers(amount int) []User {
	var users []User
	for i := 0; i < amount; i++ {
		var user User
		err := gofakeit.Struct(&user)
		if err != nil {
			panic(err)
		}

		registrationTime := getRandomRegisterTime()

		user.Birthday = getRandomAdultBirthday()
		user.RegistrationTime = registrationTime
		user.DeviceRegister, _ = getUserDevice(user.Username)
		user.AgreeTerm = registrationTime
		user.OtpTime = registrationTime

		users = append(users, user)
	}
	return users

}

func getLoginUsers(amount int) []LoginUser {
	var loginUsers []LoginUser
	for i := 0; i < amount; i++ {
		var loginUser LoginUser
		err := gofakeit.Struct(&loginUser)
		if err != nil {
			panic(err)
		}

		loginTime := getRandomRegisterTime()
		deviceName, os := getUserDevice(loginUser.Username)
		gameInfo := GameNames[gofakeit.Number(0, 4)]

		loginUser.GameName = gameInfo.GameName
		loginUser.PublisherName = gameInfo.PublisherName
		loginUser.DeviceName = deviceName
		loginUser.Os = os
		loginUser.LoginTime = loginTime

		loginUsers = append(loginUsers, loginUser)
	}
	return loginUsers
}

func getRandomDate(_startDate, _endDate time.Time) time.Time {
	var startDate time.Time
	var endDate time.Time

	if _startDate.IsZero() {
		startDate = time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC)
	} else {
		startDate = _startDate
	}

	if _endDate.IsZero() {
		endDate = time.Now()
	} else {
		endDate = _endDate
	}

	// Calculate the difference between start and end date
	diff := endDate.Sub(startDate).Hours() / 24

	log.Println(int(diff))

	// Generate a random number of days to add to the start date
	randomDays := rand.Intn(int(diff))

	// Add the random number of days to the start date
	randomDate := startDate.AddDate(0, 0, randomDays)

	return randomDate
}

func toVnDate(time time.Time) string {
	return time.Format("02/01/2006")
}

func toVnDateTime(time time.Time) string {
	return time.Format("02/01/2006 15:04:05")
}

func getRandomAdultBirthday() string {
	// Get random date from 18 years ago to now
	endDate := time.Now()
	startDate := endDate.AddDate(-18, 0, 0)
	return toVnDate(getRandomDate(time.Time{}, startDate))
}

func getRandomRegisterTime() string {
	// Get random date from 1 year ago to now
	endDate := time.Now()
	startDate := endDate.AddDate(-4, 0, 0)
	return toVnDateTime(getRandomDate(startDate, endDate))
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	const ADDRESS = "localhost:30556"
	log.Println("Starting server...")
	log.Println("Bấm Ctrl + C để dừng server")
	log.Println("Để lấy danh sách tài khoản đồng bộ, truy cập vào đường dẫn: http://" + ADDRESS + "/sync-users?amount=10")
	log.Println("Để lấy danh sách tài khoản đăng nhập, truy cập vào đường dẫn: http://" + ADDRESS + "/login-users?amount=10")

	r.GET("/sync-users", func(c *gin.Context) {
		amountStr := c.Query("amount")
		amount, _ := strconv.ParseInt(amountStr, 10, 64)
		users := getSyncUsers(int(amount))
		c.JSON(200, users)
	})

	r.GET("/login-users", func(c *gin.Context) {
		amountStr := c.Query("amount")
		amount, _ := strconv.ParseInt(amountStr, 10, 64)
		loginUsers := getLoginUsers(int(amount))
		c.JSON(200, loginUsers)
	})

	err := r.Run(ADDRESS)
	if err == nil {
		log.Println("Server is running on address: " + ADDRESS)
	}
}
