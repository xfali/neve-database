// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package gobatiseve

import (
	"errors"
	"github.com/xfali/fig"
	"github.com/xfali/gobatis"
	"github.com/xfali/gobatis/datasource"
	"github.com/xfali/gobatis/factory"
	"github.com/xfali/neve-core"
	"github.com/xfali/neve-core/container"
	"github.com/xfali/pagehelper"
	"github.com/xfali/xlog"
	"strings"
	"sync"
	"time"
)

const (
	BuildinValueDataSources = "DataSources"
)

type DataSource struct {
	DriverName string
	DriverInfo string

	MaxConn     int
	MaxIdleConn int
	//millisecond
	ConnMaxLifetime int
}

type FactoryCreatorWrapper func(f func(source *DataSource) (factory.Factory, error)) func(source *DataSource) (factory.Factory, error)

type Processor struct {
	logger        xlog.Logger
	facWrapper    FactoryCreatorWrapper
	dataSources   sync.Map
	usePageHelper bool
	gobatisLog    string
}

type Opt func(*Processor)

func NewProcessor(opts ...Opt) *Processor {
	ret := &Processor{
		logger:     xlog.GetLogger(),
		facWrapper: defaultWrapper,
	}

	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func OptSetLogger(logger xlog.Logger) Opt {
	return func(processor *Processor) {
		processor.logger = logger
	}
}

func OptFactoryCreatorWrapper(wrapper FactoryCreatorWrapper) Opt {
	return func(processor *Processor) {
		if wrapper != nil {
			processor.facWrapper = wrapper
		}
	}
}

func (p *Processor) Init(conf fig.Properties, container container.Container) error {
	dss := map[string]*DataSource{}
	err := conf.GetValue(BuildinValueDataSources, &dss)
	if err != nil {
		return err
	}
	if len(dss) == 0 {
		p.logger.Errorln("No Database")
		return nil
	}
	p.usePageHelper = conf.Get("gobatis.pagehelper.enable", "") == "true"
	p.gobatisLog = strings.ToUpper(conf.Get("gobatis.log.level", "DEBUG"))

	for k, v := range dss {
		fac, err := p.facWrapper(p.createFactory)(v)
		if err != nil {
			p.logger.Errorln("init db failed")
			return err
		}
		sm := gobatis.NewSessionManager(fac)
		p.dataSources.Store(k, sm)
		//添加到注入容器
		container.RegisterByName(k, sm)
	}

	// scan mapper files
	mapper := conf.Get("gobatis.mapper.dir", "")
	if mapper != "" {
		return gobatis.ScanMapperFile(neve.GetResource(mapper))
	}
	return nil
}

func (p *Processor) Classify(o interface{}) (bool, error) {
	switch v := o.(type) {
	case Component:
		err := p.parseBean(v)
		return true, err
	}
	return false, nil
}

func (p *Processor) parseBean(comp Component) error {
	name := comp.DataSource()
	if v, ok := p.dataSources.Load(name); ok {
		comp.SetSessionManager(v.(*gobatis.SessionManager))
	}
	p.logger.Errorln("DataSource Name found: ", name)
	return errors.New("DataSource Name found. ")
}

func (p *Processor) Process() error {
	return nil
}

func (p *Processor) createFactory(v *DataSource) (factory.Factory, error) {
	fac, err := gobatis.CreateFactory(
		gobatis.SetMaxConn(v.MaxConn),
		gobatis.SetMaxIdleConn(v.MaxIdleConn),
		gobatis.SetConnMaxLifetime(time.Duration(v.ConnMaxLifetime)*time.Millisecond),
		gobatis.SetLog(func(level int, format string, args ...interface{}) {
			p.selectLog()
		}),
		gobatis.SetDataSource(&datasource.CommonDataSource{
			Name: v.DriverName,
			Info: v.DriverInfo,
		}))
	if fac == nil || err != nil {
		return nil, err
	}
	if p.usePageHelper {
		fac = pagehelper.New(fac)
	}
	return fac, err
}

func (p *Processor) Close() error {
	p.dataSources.Range(func(key, value interface{}) bool {
		value.(*gobatis.SessionManager).Close()
		return true
	})
	return nil
}

func defaultWrapper(f func(source *DataSource) (factory.Factory, error)) func(source *DataSource) (factory.Factory, error) {
	return f
}

func (p *Processor) selectLog() func(level int, fmt string, o ...interface{}) {
	switch p.gobatisLog {
	case "DEBUG":
		return func(level int, fmt string, o ...interface{}) {
			p.logger.Debugf(fmt, o...)
		}
	case "INFO":
		return func(level int, fmt string, o ...interface{}) {
			p.logger.Infof(fmt, o...)
		}
	case "WARN":
		return func(level int, fmt string, o ...interface{}) {
			p.logger.Warnf(fmt, o...)
		}
	case "ERROR":
		return func(level int, fmt string, o ...interface{}) {
			p.logger.Errorf(fmt, o...)
		}
	}
	return func(level int, fmt string, o ...interface{}) {
		p.logger.Debugf(fmt, o...)
	}
}
