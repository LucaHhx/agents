package core

import (
	"fmt"
	"os"
	"path/filepath"
	"server/config"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	// 支持 "5s" 这种写法
	dd, err := time.ParseDuration(value.Value)
	if err != nil {
		return err
	}
	*d = Duration(dd)
	return nil
}

type Manager struct {
	path  string
	value atomic.Value // 存 *Config
}

func NewManager(path string) (*Manager, error) {
	m := &Manager{path: path}
	cfg, err := loadOnce(path)
	if err != nil {
		return nil, err
	}
	m.value.Store(cfg)
	return m, nil
}

func (m *Manager) Get() *config.Config {
	return m.value.Load().(*config.Config)
}

func loadOnce(path string) (*config.Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Watch 会阻塞运行；建议放 goroutine；stop 关闭即可退出
func (m *Manager) Watch(stop <-chan struct{}) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	// 监听目录更稳：很多编辑器是"写临时文件 + rename"
	dir := filepath.Dir(m.path)
	base := filepath.Base(m.path)

	if err := w.Add(dir); err != nil {
		return err
	}

	// 简单防抖：编辑器可能触发多次事件
	var (
		timer   *time.Timer
		timerCh <-chan time.Time
	)

	scheduleReload := func() {
		if timer == nil {
			timer = time.NewTimer(200 * time.Millisecond)
			timerCh = timer.C
			return
		}
		if !timer.Stop() {
			select { // 清空可能残留的 tick
			case <-timer.C:
			default:
			}
		}
		timer.Reset(200 * time.Millisecond)
	}

	reload := func() {
		cfg, err := loadOnce(m.path)
		if err != nil {
			// 重载失败：保留旧配置，打印/上报即可
			fmt.Printf("[config] reload failed: %v\n", err)
			return
		}
		m.value.Store(cfg)
		fmt.Printf("[config] reloaded OK\n")
	}

	for {
		select {
		case <-stop:
			return nil

		case err := <-w.Errors:
			// watcher 自身错误一般不致命，但你可以选择 return
			fmt.Printf("[config] watcher error: %v\n", err)

		case ev := <-w.Events:
			// 只关心目标文件
			if filepath.Base(ev.Name) != base {
				continue
			}
			// 常见触发：Write/Create/Rename
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				scheduleReload()
			}

		case <-timerCh:
			timerCh = nil
			reload()
		}
	}
}
