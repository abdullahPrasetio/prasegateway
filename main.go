package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdullahPrasetio/prasegateway/config"
	"github.com/abdullahPrasetio/prasegateway/routers"
	"github.com/fsnotify/fsnotify"
)

var server *http.Server

func main() {

	// var myConfig entity.MyConfig
	// file, err := os.ReadFile("prase.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := json.Unmarshal(file, &myConfig); err != nil {
	// 	log.Fatal(err)
	// }

	myConfig := config.GetMyConfig()

	router := routers.Setup(myConfig)

	port := "8080"
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
		Handler:      router,
	}

	go watchConfig()

	server = srv

	go func() {

		fmt.Println("myconfig", myConfig)
		fmt.Println("Server run in url : http://localhost:" + port)
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("panic", err)
			// panic(err)
		}
	}()

	// ====== Mematikan server secara gracefull ======

	// Membuat channel untuk menerima sinyal shutdown
	quit := make(chan os.Signal)

	// Menerima sinyal untuk mematikan server dari os
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Memberikan waktu timeout 5 detik untuk mematikan server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Jika lebih dari 5 detik matikan paksa
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

}

func watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	// Awasi perubahan pada file konfigurasi
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Write {
					fmt.Println("Config file changed:", event.Name)
					// Hentikan server
					server.Shutdown(nil)

					// Tunggu sebentar agar server berhenti sepenuhnya
					time.Sleep(2 * time.Second)

					// Jalankan kembali server
					// if err := server.ListenAndServe(); err != nil {
					// 	panic(err)
					// }
					config.InitializeConfig()
					main()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Error watching file:", err)
			}
		}
	}()

	// Tambahkan file konfigurasi ke watcher
	err = watcher.Add("prase.json")
	if err != nil {
		fmt.Println("Error adding file to watcher:", err)
		return
	}

	// Tunggu goroutine memantau perubahan
	select {}
}
