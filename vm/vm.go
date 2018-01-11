package vm

import (
	"crypto/md5"
	"errors"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"

	fsnotify "gopkg.in/fsnotify.v0"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	"github.com/orktes/huealexa/hueserver"
)

type VM struct {
	*goja.Runtime
	sync.Mutex
	server  *hueserver.Server
	watcher *fsnotify.Watcher
	dataDir string
	md5Map  map[string][md5.Size]byte
	path    string
}

func (vm *VM) RunString(str string) (goja.Value, error) {
	vm.Lock()
	defer vm.Unlock()

	return vm.Runtime.RunString(str)
}

func (vm *VM) RunScript(name, value string) (goja.Value, error) {
	vm.Lock()
	defer vm.Unlock()

	return vm.Runtime.RunScript(name, value)
}

func (vm *VM) srcLoader(pathname string) ([]byte, error) {
	if !strings.HasSuffix(pathname, ".js") {
		pathname += ".js"
	}

	asset, err := Asset(pathname)
	if err == nil {
		return asset, nil
	}

	pathToFile := path.Join(path.Dir(vm.path), pathname)
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	md5 := md5.Sum(data)
	vm.md5Map[pathToFile] = md5

	return data, vm.watcher.Watch(pathToFile)
}

func (vm *VM) register() {
	registry := require.NewRegistryWithLoader(vm.srcLoader)
	registry.Enable(vm.Runtime)
	console.Enable(vm.Runtime)
	vm.initServer()
	vm.initAlexa()
	vm.initBase64()
	vm.initFS()
	vm.initPath()
	vm.initAsync()
	vm.initHomeKit()
	vm.initTimers()
	vm.initUUID()
	vm.initProcess()
	vm.initSSDP()
	vm.initZWay()
	vm.initWebSocket()
	vm.initWOL()
}

func (vm *VM) startWatch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	vm.watcher = watcher

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					log.Printf("%s changed. Reinitializing VM if needed.", ev.Name)

					data, err := ioutil.ReadFile(ev.Name)
					if err != nil {
						log.Printf("Unable to read file %s: %s", ev.Name, err.Error())
						continue
					}

					md5 := md5.Sum(data)

					if vm.md5Map[ev.Name] == md5 {
						log.Print("Content didn't change. Doing nothing.\n")
						continue
					}

					vm.md5Map[ev.Name] = md5
					vm.init()
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	return watcher.Watch(vm.path)
}

func (vm *VM) init() (err error) {
	value, err := ioutil.ReadFile(vm.path)
	if err != nil {
		return
	}

	md5 := md5.Sum(value)
	vm.md5Map[vm.path] = md5

	if vm.Runtime != nil {
		vm.Runtime.Interrupt(errors.New("Interreupted due to update"))
	}

	log.Printf("Initializing VM with %s\n", vm.path)
	vm.Runtime = goja.New()
	// Set env
	vm.Set("env", map[string]interface{}{
		"data_dir":      vm.dataDir,
		"huealexa_uuid": vm.server.UUID,
	})

	vm.register()

	_, err = vm.RunScript(vm.path, string(value))

	log.Printf("Done initializing %s\n", vm.path)

	return err
}

func (vm *VM) Close() {
	vm.watcher.Close()
}

func NewVM(path string, dataDir string, server *hueserver.Server) (*VM, error) {
	vm := &VM{
		dataDir: dataDir,
		server:  server,
		path:    path,
		md5Map:  map[string][md5.Size]byte{},
	}
	vm.startWatch()
	err := vm.init()
	return vm, err
}
