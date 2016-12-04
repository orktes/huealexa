package vm

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	fsnotify "gopkg.in/fsnotify.v0"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

type VM struct {
	*goja.Runtime
	sync.Mutex
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
	initAsync(vm)
	initTimers(vm)
	initUUID(vm)
	initProcess(vm)
	initSSDP(vm)
}

func (vm *VM) startWatch(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					log.Printf("%s changed. Reinitializing VM", path)
					vm.initWithPath(path)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	return watcher.Watch(path)
}

func (vm *VM) initWithPath(path string) (err error) {
	if vm.Runtime != nil {
		vm.Runtime.Interrupt(errors.New("Interreupted due to update"))
	}

	log.Printf("Initializing VM with %s\n", path)
	vm.Runtime = goja.New()
	vm.register()

	script, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	_, err = vm.RunScript(path, string(script))

	log.Printf("Done initializing %s\n", path)
	return
}

func NewVM(path string) (*VM, error) {
	vm := &VM{}
	vm.startWatch(path)
	err := vm.initWithPath(path)
	return vm, err
}
