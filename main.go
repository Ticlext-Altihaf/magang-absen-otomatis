package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/viper"
)

func main() {
	// Try to run the real main function
	// If it panics, recover and print the error
	err := beeep.Alert("Absen Magang", "ikuzoo!", "assets/warning.png")
	if err != nil {
		fmt.Println("Failed to send notification:", err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
			err := beeep.Alert("Absen Magang Error", fmt.Sprint(r), "")
			if err != nil {
				fmt.Println("Failed to send alert:", err)
			}
		}
	}()
	init_w()
	loop()
}

type Config struct {
	Email         string    `yaml:"email"`
	Password      string    `yaml:"password"`
	Presence_type string    `yaml:"presence_type"`
	Jam_absen     time.Time `yaml:"jam_absen"`
	Jam_pulang    time.Time `yaml:"jam_pulang"`
}

var config Config

func init_w() {
	// Load the configuration file
	viper.SetConfigFile("config.yaml")
	//check if config file exists
	if err := viper.ReadInConfig(); err != nil {
		//make config file if not exists
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set("email", "")
			viper.Set("password", "")
			viper.Set("presence_type", "1")
			viper.Set("jam_absen", "08:00")
			viper.Set("jam_pulang", "16:00")
			viper.WriteConfig()
			panic("Mohon isi config.yaml dengan username dan password")
		} else {
			//throw error if other error
			panic(err)
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic("Gagal membaca config.yaml: " + err.Error())
	}

	if viper.GetString("email") == "" || viper.GetString("password") == "" {
		panic("Mohon isi config.yaml dengan username dan password")
	}

	if viper.GetString("jam_absen") == "" || viper.GetString("jam_pulang") == "" {
		panic("Mohon isi config.yaml dengan jam_absen dan jam_pulang")
	}

	config.Email = viper.GetString("email")
	config.Password = viper.GetString("password")
	config.Presence_type = viper.GetString("presence_type")
	config.Jam_absen, err = time.Parse("15:04", viper.GetString("jam_absen"))
	if err != nil {
		panic("Gagal parsing jam_absen: " + err.Error())
	}
	config.Jam_pulang, err = time.Parse("15:04", viper.GetString("jam_pulang"))
	if err != nil {
		panic("Gagal parsing jam_pulang: " + err.Error())
	}
	println("Jam absen: ", config.Jam_absen.Format("15:04"))
	println("Jam pulang: ", config.Jam_pulang.Format("15:04"))
	res, err := absen_login()
	if err != nil {
		panic("Gagal login: " + err.Error())
	}
	if res == "" {
		panic("Gagal mendapatkan cookies: " + err.Error())
	}
}
func loop() {
	//resume loop regardless of panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
			err := beeep.Alert("Absen Magang Error", fmt.Sprint(r), "")
			if err != nil {
				fmt.Println("Failed to send alert:", err)
			}
			//wait a moment
			time.Sleep(5 * time.Second)
			loop()
		}
	}()
	sudah_absen_masuk := false

	for {
		//now date (23/2/20 00:00) + hours (08:00) = jam absen (23/2/20 08:00)
		jam_pulang := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.Jam_pulang.Hour(), config.Jam_pulang.Minute(), config.Jam_pulang.Second(), config.Jam_pulang.Nanosecond(), time.Now().Location())
		jam_absen := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.Jam_absen.Hour(), config.Jam_absen.Minute(), config.Jam_absen.Second(), config.Jam_absen.Nanosecond(), time.Now().Location())
		now := time.Now()
		// Sekarang harus apa ?
		if sudah_absen_masuk && now.Before(jam_pulang) {
			//sudah absen masuk dan belum jam pulang
			durasi := jam_pulang.Sub(now)
			println("Menunggu jam pulang: ", fmt.Sprintf("%02d:%02d:%02d", int(durasi.Hours()), int(durasi.Minutes())%60, int(durasi.Seconds())%60), " (", jam_pulang.Format("15:04"), ")")
		} else if !sudah_absen_masuk && now.Before(jam_absen) {
			//belum absen masuk dan belum jam absen
			durasi := jam_absen.Sub(now)
			println("Menunggu jam absen: ", fmt.Sprintf("%02d:%02d:%02d", int(durasi.Hours()), int(durasi.Minutes())%60, int(durasi.Seconds())%60), " (", jam_absen.Format("15:04"), ")")
		} else if sudah_absen_masuk && now.After(jam_pulang) {
			println("Sudah jam pulang")
		} else if !sudah_absen_masuk && now.After(jam_pulang) {
			println("Sudah jam pulang, menunggu sampai besok")
		}
		if now.After(jam_absen) && now.Before(jam_pulang) && !sudah_absen_masuk {

			cookie, err := absen_login()
			if err != nil {
				panic("Gagal login: " + err.Error())
			}
			// Absen masuk
			res, err := absen_masuk(cookie)
			if err != nil {
				panic("Gagal absen masuk: " + err.Error())
			}
			fmt.Println("Absen Masuk: ", res)
			beeep.Alert("Absen Masuk", res, "assets/info.png")
			if res == "Kehadiran berhasil disimpan" || res == "Anda telah melakukan Absen Mulai" {
				sudah_absen_masuk = true
			}
		} else if now.After(jam_pulang) && sudah_absen_masuk {

			cookie, err := absen_login()
			if err != nil {
				panic("Gagal login: " + err.Error())
			}
			// Absen pulang
			res, err := absen_pulang(cookie)
			if err != nil {
				panic("Gagal absen pulang: " + err.Error())
			}
			fmt.Println("Absen Pulang: ", res)
			beeep.Alert("Absen Pulang", res, "assets/info.png")
			if res == "Kehadiran berhasil disimpan" {
				sudah_absen_masuk = false
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func get_start_form_data() (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	image := "absen.png"
	file, err := os.Open(image)
	if err != nil {
		panic("Gagal membuka file: " + err.Error())
	}
	presence_type := viper.GetString("presence_type")
	if presence_type == "" {
		panic("Mohon isi config.yaml dengan presence_type")
	}

	part, err := writer.CreateFormFile("image", image)
	if err != nil {
		panic("Gagal membuat part: " + err.Error())
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic("Gagal menulis part: " + err.Error())
	}
	err = writer.WriteField("presence_type", presence_type)
	if err != nil {
		panic("Gagal menulis field: " + err.Error())
	}
	err = writer.Close()
	if err != nil {
		panic("Gagal menutup writer: " + err.Error())
	}
	return body, writer.FormDataContentType()
}

func absen_login() (string, error) {
	fmt.Println("[" + time.Now().Format("15:04") + "] Login...")
	// Create an HTTP client
	client := &http.Client{}

	// Login to the website
	loginData := url.Values{
		"email":    {viper.GetString("email")},
		"password": {viper.GetString("password")},
	}

	loginURL := "https://pkl.smknegeri1garut.sch.id/login"
	//get cookies first
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return "", fmt.Errorf("Gagal membuat permintaan login: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gagal melakukan permintaan login: %v", err)
	}
	defer resp.Body.Close()
	cookie := resp.Header.Get("Set-Cookie")
	fmt.Println("Cookie:", len(cookie))
	if cookie == "" {
		return "", fmt.Errorf("Gagal mendapatkan cookie")
	}

	req, err = http.NewRequest("POST", loginURL, strings.NewReader(loginData.Encode()))
	req.Header.Add("Cookie", cookie)

	if err != nil {
		return "", fmt.Errorf("Gagal membuat permintaan login: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Gagal melakukan permintaan login: %v", err)
	}
	defer resp.Body.Close()

	// Check the login response
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Login gagal. Status code: %v", resp.Status)
	}
	// Read the fetched data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read data: %v", err)
	}

	// Write the data to a file
	err = ioutil.WriteFile("login.html", data, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write data to file: %v", err)
	}

	// Read danger alert
	messages, err := parse_html_text(string(data), ".alert-danger")
	if messages != nil {
		return "", fmt.Errorf("%v", strings.Join(messages, ""))
	}
	if err != nil {
		return "", fmt.Errorf("Failed to parse html: %v", err)
	}

	actualDataURL := "https://pkl.smknegeri1garut.sch.id/partisipant"
	currentURL := resp.Request.URL.String()
	if currentURL != actualDataURL {
		//try get error message
		messages, err := parse_html_text(string(data), ".alert-danger")
		if err == nil {
			return "", fmt.Errorf("%v", strings.Join(messages, ""))
		} else {
			println(err)
		}

		return "", fmt.Errorf("Login gagal, redirect ke halaman: %v", currentURL)
	}
	fmt.Println("Login berhasil:", resp.Status)
	if resp.Header.Get("Set-Cookie") != "" {
		cookie = resp.Header.Get("Set-Cookie")
	}

	return cookie, nil
}

func absen_masuk(cookie string) (string, error) {
	fmt.Println("[" + time.Now().Format("15:04") + "] Absen Masuk...")
	// Create an HTTP client
	client := &http.Client{}

	postURL := "https://pkl.smknegeri1garut.sch.id/partisipant/mulai"
	postData, contentType := get_start_form_data()
	req, err := http.NewRequest("POST", postURL, postData)
	req.Header.Add("Cookie", cookie)
	if err != nil {
		return "", fmt.Errorf("Failed to create post request: %v", err)
	}
	req.Header.Add("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to perform post request: %v", err)
	}
	defer resp.Body.Close()

	// Check the post response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Post failed. Status code: %v", resp.StatusCode)
	}
	// Write data to file
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read data: %v", err)
	}
	err = ioutil.WriteFile("post-mulai.html", data, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write data to file: %v", err)
	}

	// Read alert
	messages, err := parse_html_text(string(data), ".alert")
	if messages != nil {
		return strings.Join(messages, ""), nil
	}
	if err != nil {
		return "", fmt.Errorf("Failed to parse html: %v", err)
	}
	return "", nil
}

func absen_pulang(cookie string) (string, error) {
	fmt.Println("Absen pulang...")
	client := &http.Client{}

	postURL := "https://pkl.smknegeri1garut.sch.id/partisipant/selesai"
	postData, contentType := get_start_form_data() //same as start
	req, err := http.NewRequest("POST", postURL, postData)
	req.Header.Add("Cookie", cookie)
	if err != nil {
		return "", fmt.Errorf("Failed to create post request: %v", err)
	}
	req.Header.Add("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to perform post request: %v", err)
	}
	defer resp.Body.Close()

	// Check the post response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Post failed. Status code: %v", resp.StatusCode)
	}
	// Write data to file
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read data: %v", err)
	}
	err = ioutil.WriteFile("post-selesai.html", data, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write data to file: %v", err)
	}

	// Read alert
	messages, err := parse_html_text(string(data), ".alert")
	if messages != nil {
		return strings.Join(messages, ""), nil
	}
	if err != nil {
		return "", fmt.Errorf("Failed to parse html: %v", err)
	}
	return "", nil
}
