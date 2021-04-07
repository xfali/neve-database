// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/fig"
	"github.com/xfali/gobatis"
	"github.com/xfali/neve-core"
	"github.com/xfali/neve-database/gobatiseve"
	"testing"
)

type MyTest struct {
	SessMgr *gobatis.SessionManager `inject:"testDB"`
}

type test struct {
	Name string `xfield:"name"`
}

func (t *MyTest) test() {
	sess := t.SessMgr.NewSession()
	var ret []test
	sess.Select("select * from test").Param().Result(&ret)
}

func TestGobatis(t *testing.T) {
	p := gobatiseve.NewProcessor()
	conf, err := fig.LoadYamlFile("../assets/config-example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	p.Init(conf, nil)
}

func TestGobatisInject(t *testing.T) {
	app := neve.NewFileConfigApplication("assets/config-example.yaml")
	app.RegisterBean(gobatiseve.NewProcessor())
	app.RegisterBean(&MyTest{})
	//Register other beans...
	//app.RegisterBean(test.NewMyTest())
	app.Run()
}
