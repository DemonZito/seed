//+build !wasm

package seed

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/qlova/seed/user"
)

//import ua "github.com/avct/uasurfer"

//Live signals if the application is live or not.
var Live bool

//Production signals if the application is in production or not.
var Production bool

var exporting bool

func init() {
	for _, arg := range os.Args {
		if arg == "-live" {
			Live = true
		}
	}
	for _, arg := range os.Args {
		if arg == "-production" {
			Production = true
			user.Production = true
			Live = true
		}
	}
	for _, arg := range os.Args {
		if arg == "-export=static" {
			exporting = true
		}
	}
}

var defers []func()

func cleanup() {
	for _, f := range defers {
		f()
	}
}
func deferFunction(f func()) {
	defers = append(defers, f)
}

func init() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
}

//Launch launches the app listening on the given port.
func (runtime Runtime) Launch(port ...string) {
	if runtime.app.Seed.seed != nil {

		if len(port) > 0 {
			runtime.Listen = port[0]
		}

		//Allow port config from Env
		if port := os.Getenv("PORT"); port != "" {
			runtime.Listen = port
		}

		if runtime.Listen == "" {
			runtime.Listen = ":1234"
		}

		if !Live && !Production {

			if runtime.bootstrapWasm {
				runtime.launchWasm()
				return
			}

			//Launch the app if possible.
			go launch(":10000")

			var Process = exec.Command(os.Args[0], "-live")
			Process.Stdout = os.Stdout
			Process.Stderr = os.Stderr
			Process.Start()

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}
			defer watcher.Close()

			var Compiler *exec.Cmd

			deferFunction(func() {
				if Process.Process != nil {
					Process.Process.Kill()
				}
				if Compiler != nil && Compiler.Process != nil {
					Compiler.Process.Kill()
				}
			})

			var Compiling bool

			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						//log.Println("event:", event)
						if event.Op&fsnotify.Write == fsnotify.Write {

							if path.Ext(event.Name) == ".go" {

								if Compiling {
									continue
								}

								Compiler = exec.Command("go", "build", "-i", "-o", os.Args[0])
								Compiling = true
								go func() {
									err := Compiler.Run()
									if err == nil {
										if Process.Process != nil {
											Process.Process.Kill()
										}
										Process = exec.Command(os.Args[0], "-live")
										Process.Stdout = os.Stdout
										Process.Stderr = os.Stderr
										Process.Start()

										reloading = true
										for _, socket := range localSockets {
											socket.WriteMessage(1, []byte("update();"))
										}
									} else {
										println(err.Error())
									}
									Compiling = false

								}()

							}

						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						log.Println("error:", err)
					}
				}
			}()

			err = watcher.Add(path.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}

			proxy(runtime.Listen, ":10000")

		} else {
			http.Handle("/", runtime.app.Handler())
			http.ListenAndServe(runtime.Listen, nil)
		}
		return
	}
	panic("No seeds!")
}