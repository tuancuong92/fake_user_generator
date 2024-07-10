package main

import (
	"log"
	"strconv"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
)

type User struct {
	GamerId          string `fake:"{uuid}" json:"gamerId"`
	Username         string `fake:"{username}" json:"userName"`
	FullName         string `fake:"{name}" json:"fullName"`
	Birthday         string `fake:"{day}/{month}/{year}" format:"31/12/2006" json:"birthday"`
	IdentityNumber   string `fake:"{number:12}" json:"identityNumber"`
	PhoneNumber      string `fake:"{phone}" json:"phoneNumber"`
	Email            string `fake:"{email}" json:"email"`
	Address          string `fake:"{street}, {city}" json:"address"`
	RegistrationTime string `fake:"{day}/{month}/{year} {hour}:{minute}:{second}" format:"31/12/2006 02:34:56" json:"registrationTime"`
	DeviceRegister   string `fake:"{useragent}" json:"deviceRegister"`
	Ip               string `fake:"{ipv4address}" json:"ip"`
	AgreeTerm        string `fake:"{day}/{month}/{year} {hour}:{minute}:{second}" format:"31/12/2006 02:34:56" json:"agreeTerm"`
}

type LoginUser struct {
	Username      string `fake:"{username}" json:"userName"`
	GamerId       string `fake:"{uuid}" json:"gamerId"`
	GameName      string `json:"gameName"`
	PublisherName string `json:"publisherName"`
	Ip            string `fake:"{ipv4address}" json:"ip"`
	DeviceId      string `fake:"{number:12}" json:"deviceId"`
	DeviceName    string `json:"deviceName"`
	Os            string `json:"os"`
	LoginTime     string `fake:"{day}/{month}/{year} {hour}:{minute}:{second}" format:"31/12/2006 02:34:56" json:"loginTime"`
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

		deviceName, os := getUserDevice(loginUser.Username)
		gameInfo := GameNames[gofakeit.Number(0, 4)]
		loginUser.GameName = gameInfo.GameName
		loginUser.PublisherName = gameInfo.PublisherName
		loginUser.DeviceName = deviceName
		loginUser.Os = os
		loginUsers = append(loginUsers, loginUser)
	}
	return loginUsers
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
