// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/fig"
	"github.com/xfali/neve-database/gobatiseve"
	"testing"
)

func TestGobatis(t *testing.T) {
	p := gobatiseve.NewProcessor()
	conf, err := fig.LoadYamlFile("../assets/config-example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	p.Init(conf, nil)
}
