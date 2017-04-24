package vm

import (
	"crypto/md5"
	"errors"
	"io/ioutil"
	"log"
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
	md5     [md5.Size]byte
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

	return nil, errors.New("Package " + pathname + " not found")
}

func (vm *VM) register() {
	registry := require.NewRegistryWithLoader(vm.srcLoader)
	registry.Enable(vm.Runtime)
	console.Enable(vm.Runtime)
	vm.initServer()
	vm.initAlexa()
	vm.initFS()
	vm.initPath()
	vm.initAsync()
	vm.initHomeKit()
	vm.initTimers()
	vm.initUUID()
	vm.initProcess()
	vm.initSSDP()
	vm.initZWay()
}

func (vm *VM) startWatch(path string) error {
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
					log.Printf("%s changed. Reinitializing VM if needed.", path)
					vm.initWithPath(path)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	return watcher.Watch(path)
}

func (vm *VM) init(path, value string) (err error) {
	md5 := md5.Sum([]byte(value))

	if vm.md5 == md5 {
		log.Print("Content didn't change. Doing nothing.\n")
		return
	}

	vm.md5 = md5

	if vm.Runtime != nil {
		vm.Runtime.Interrupt(errors.New("Interreupted due to update"))
	}

	log.Printf("Initializing VM with %s\n", path)
	vm.Runtime = goja.New()
	// Set env
	vm.Set("env", map[string]interface{}{
		"data_dir":      vm.dataDir,
		"huealexa_uuid": vm.server.UUID,
	})

	vm.register()

	_, err = vm.RunScript(path, value)

	log.Printf("Done initializing %s\n", path)

	return err
}

func (vm *VM) initWithPath(path string) (err error) {
	script, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	return vm.init(path, string(script))
}

func (vm *VM) Close() {
	vm.watcher.Close()
}

func NewVM(path string, dataDir string, server *hueserver.Server) (*VM, error) {
	vm := &VM{dataDir: dataDir, server: server}
	vm.startWatch(path)
	err := vm.initWithPath(path)
	return vm, err
}
