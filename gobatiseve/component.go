// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package gobatiseve

import "github.com/xfali/gobatis"

type Component interface {
	// Return: 需要注入的数据源名称
	DataSource() string
	// Param: o - 注入SessionManager的目的对象
	// Return: 注入成功返回true，否则返回false
	SetSessionManager(manager *gobatis.SessionManager) error
}
