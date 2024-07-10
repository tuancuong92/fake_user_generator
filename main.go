package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"os"
	"strconv"
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

func main() {
	var amount uint
	if len(os.Args) == 2 {
		_amount, err := strconv.ParseUint(os.Args[1], 10, 32)
		if err != nil {
			panic("Xin vui long nhap so luong nguoi dung!")
		}
		amount = uint(_amount)
	} else {
		fmt.Print("Xin vui long nhap so luong nguoi dung: ")

		n, err := fmt.Scanln(&amount)
		if err != nil || n != 1 || amount <= 0 {
			panic("Xin vui long nhap so luong nguoi dung!")
		}
	}

	var user User
	var result []User

	for i := 0; i < int(amount); i++ {
		err := gofakeit.Struct(&user)
		if err != nil {
			panic(err)
		}
		result = append(result, user)
	}

	// Step 1: Open a file
	file, err := os.Create("output.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	jsonData, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}

	_, err = file.Write(jsonData)
	if err != nil {
		panic(fmt.Sprintf("Failed to write JSON to file: %v", err))
	}

	fmt.Println("Xin vui long kiem tra file output.json!")
}
