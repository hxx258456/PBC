package configs

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"log"
	"os"
	"path/filepath"
	"pyramid/pyramid-manage/backend/app/backend/internal/conf"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-kratos/kratos/v2/config"
)

var (
	confPath string
	Conf     conf.Bootstrap
)
var (
	DeployEnv  string
	AppID      string
	RemoteUrl  string
	RemotePort int
	RemoteUser string
	RemotePW   string
)

func init() {
	env()
	if err := Init(); err != nil {
		panic(err)
	}
}

func Init() error {
	return local()
}

func flagInit() {
	flag.StringVar(&confPath, "conf", confPath, "config path, eg: -conf config.yaml")
	flag.Parse()
	fmt.Println("=====", confPath)
	if confPath == "" {
		confPath = "../../configs/config.yaml"
	}
}

func env() {
	flagInit()
	AppID = os.Getenv("AppID")
}

func local() error {
	c := config.New(
		config.WithSource(
			file.NewSource(confPath),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		return err
	}
	return c.Scan(&Conf)
}

func remote() error {
	c := config.New(
		config.WithSource(
			remoteLoad(),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		return err
	}
	return c.Scan(&Conf)
}

func remoteLoad() config.Source {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(RemoteUrl, uint64(RemotePort)),
	}
	cc := &constant.ClientConfig{
		NamespaceId:         AppID, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              fmt.Sprintf("/tmp/%s/log", AppID),
		CacheDir:            fmt.Sprintf("/tmp/%s/cache", AppID),
		LogLevel:            "debug",
		Username:            RemoteUser,
		Password:            RemotePW,
	}
	// a more graceful way to create naming client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Panic(err)
	}
	source := NewConfigSource(client, WithGroup(AppID), WithDataID(AppID))
	return source
}

type Option func(*options)

type options struct {
	group  string
	dataID string
}

// WithGroup With nacos config group.
func WithGroup(group string) Option {
	return func(o *options) {
		o.group = group
	}
}

// WithDataID With nacos config data id.
func WithDataID(dataID string) Option {
	return func(o *options) {
		o.dataID = fmt.Sprintf("%s.yaml", dataID)
	}
}

type Config struct {
	opts   options
	client config_client.IConfigClient
}

func NewConfigSource(client config_client.IConfigClient, opts ...Option) config.Source {
	_options := options{}
	for _, o := range opts {
		o(&_options)
	}
	return &Config{client: client, opts: _options}
}

func (c *Config) Load() ([]*config.KeyValue, error) {
	content, err := c.client.GetConfig(vo.ConfigParam{
		DataId: c.opts.dataID,
		Group:  c.opts.group,
	})
	if err != nil {
		return nil, err
	}
	k := c.opts.dataID
	return []*config.KeyValue{
		{
			Key:    k,
			Value:  []byte(content),
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		},
	}, nil
}

func (c *Config) Watch() (config.Watcher, error) {
	watcher := newWatcher(context.Background(), c.opts.dataID, c.opts.group, c.client.CancelListenConfig)
	err := c.client.ListenConfig(vo.ConfigParam{
		DataId: c.opts.dataID,
		Group:  c.opts.group,
		OnChange: func(_, group, dataId, data string) {
			if dataId == watcher.dataID && group == watcher.group {
				watcher.content <- data
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

type Watcher struct {
	dataID             string
	group              string
	content            chan string
	cancelListenConfig cancelListenConfigFunc

	ctx    context.Context
	cancel context.CancelFunc
}

type cancelListenConfigFunc func(params vo.ConfigParam) (err error)

func newWatcher(ctx context.Context, dataID string, group string, cancelListenConfig cancelListenConfigFunc) *Watcher {
	ctx, cancel := context.WithCancel(ctx)
	w := &Watcher{
		dataID:             dataID,
		group:              group,
		cancelListenConfig: cancelListenConfig,
		content:            make(chan string, 100),

		ctx:    ctx,
		cancel: cancel,
	}
	return w
}

func (w *Watcher) Next() ([]*config.KeyValue, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case content := <-w.content:
		k := w.dataID
		return []*config.KeyValue{
			{
				Key:    k,
				Value:  []byte(content),
				Format: strings.TrimPrefix(filepath.Ext(k), "."),
			},
		}, nil
	}
}

func (w *Watcher) Close() error {
	err := w.cancelListenConfig(vo.ConfigParam{
		DataId: w.dataID,
		Group:  w.group,
	})
	w.cancel()
	return err
}

func (w *Watcher) Stop() error {
	return w.Close()
}
