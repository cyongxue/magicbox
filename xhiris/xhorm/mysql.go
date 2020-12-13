package xhorm

import (
	"database/sql"
	"fmt"
	"github.com/cyongxue/magicbox/xhiris/xhlog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"time"
)

type ORMEngine struct {
	Engine *xorm.Engine
}

type MysqlConfig struct {
	Host     string
	Port     uint16
	DBName   string
	Username string
	Password string
	MaxIdle  int
	MaxOpen  int
	Debug    bool
}

// NewEngine mysql连接初始化处理
func NewEngine(config *MysqlConfig) (*ORMEngine, error) {
	// 初始化engine
	var connectionStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		config.Username, config.Password, config.Host, config.Port, config.DBName)

	engine, err := xorm.NewEngine("mysql", connectionStr)
	if err != nil {
		xhlog.AppError(xhlog.OPUndef, map[string]interface{}{
			"msg":      "init mysql database failed",
			xhlog.Args: connectionStr,
		})
		return nil, err
	}
	if err := engine.Ping(); err != nil {
		xhlog.AppError(xhlog.OPUndef, map[string]interface{}{
			"msg":      "connect to database failed",
			xhlog.Args: connectionStr,
		})
		return nil, err
	}

	engine.SetMaxIdleConns(config.MaxIdle)
	engine.SetMaxOpenConns(config.MaxOpen)
	engine.ShowSQL(config.Debug)
	engine.ShowExecTime(config.Debug)

	ORMEngine := ORMEngine{
		Engine: engine,
	}
	return &ORMEngine, nil
}

// ORMEngine.ExecWithCtx 用于支持insert、update、del等写入操作
func (o *ORMEngine) ExecWithCtx(ctx iris.Context, query *builder.Builder) (sql.Result, error) {
	startTime := time.Now()

	cmd, args, err := query.ToSQL()
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			xhlog.ErrorMsg: fmt.Sprintf("query.ToSQL() fail, %s", err.Error()),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, err
	}

	res, err := o.Engine.Exec(cmd, args)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return res, err
	}

	rowsAffected, _ := res.RowsAffected()
	lastInsertId, _ := res.LastInsertId()
	xhlog.Info(ctx, xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		"rowsAffected": rowsAffected,
		"lastInsertId": lastInsertId,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return res, err
}

// ORMEngine.ExecWithCtx 用于支持insert、update、del等写入操作
func (o *ORMEngine) ExecWithCtxArgs(ctx iris.Context, cmd string, args ...interface{}) (sql.Result, error) {
	startTime := time.Now()

	res, err := o.Engine.Exec(cmd, args)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return res, err
	}

	rowsAffected, _ := res.RowsAffected()
	lastInsertId, _ := res.LastInsertId()
	xhlog.Info(ctx, xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		"rowsAffected": rowsAffected,
		"lastInsertId": lastInsertId,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return res, err
}

// ORMEngine.ExecNoCtx 用于支持insert、update、del等写入操作
func (o *ORMEngine) ExecNoCtx(query *builder.Builder) (sql.Result, error) {
	startTime := time.Now()

	cmd, args, err := query.ToSQL()
	if err != nil {
		xhlog.AppError(xhlog.OPMysql, map[string]interface{}{
			xhlog.ErrorMsg: fmt.Sprintf("query.ToSQL() fail, %s", err.Error()),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, err
	}

	res, err := o.Engine.Exec(cmd, args)
	if err != nil {
		xhlog.AppError(xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return res, err
	}

	rowsAffected, _ := res.RowsAffected()
	lastInsertId, _ := res.LastInsertId()
	xhlog.AppInfo(xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		"rowsAffected": rowsAffected,
		"lastInsertId": lastInsertId,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return res, err
}

// ORMEngine.GetWithCtx 用于单条记录查询
func (o *ORMEngine) GetWithCtx(ctx iris.Context, result interface{}, query *builder.Builder) (interface{}, bool, error) {
	startTime := time.Now()

	cmd, args, err := query.ToSQL()
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			xhlog.ErrorMsg: fmt.Sprintf("query.ToSQL() fail, %s", err.Error()),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, false, err
	}

	has, err := o.Engine.SQL(cmd, args...).Get(result)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, false, err
	}

	xhlog.Info(ctx, xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		"has":          has,
		xhlog.Result:   result,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return result, has, nil
}

// ORMEngine.GetNoCtx 用于单条记录查询
func (o *ORMEngine) GetNoCtx(result interface{}, query *builder.Builder) (interface{}, bool, error) {
	startTime := time.Now()

	cmd, args, err := query.ToSQL()
	if err != nil {
		xhlog.AppError(xhlog.OPMysql, map[string]interface{}{
			xhlog.ErrorMsg: fmt.Sprintf("query.ToSQL() fail, %s", err.Error()),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, false, err
	}

	has, err := o.Engine.SQL(cmd, args...).Get(result)
	if err != nil {
		xhlog.AppError(xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, false, err
	}

	xhlog.AppInfo(xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		"has":          has,
		xhlog.Result:   result,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return result, has, nil
}

// ORMEngine.FindWithCtx 用于多条记录查询
func (o *ORMEngine) FindWithCtx(ctx iris.Context, result interface{}, query *builder.Builder) (interface{}, error) {
	startTime := time.Now()

	cmd, args, err := query.ToSQL()
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			xhlog.ErrorMsg: fmt.Sprintf("query.ToSQL() fail, %s", err.Error()),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, err
	}

	err = o.Engine.SQL(cmd, args...).Find(result)
	if err != nil {
		xhlog.Error(ctx, xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, err
	}

	xhlog.Info(ctx, xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		xhlog.Result:   result,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return result, nil
}

// ORMEngine.FindNoCtx 用于多条记录查询
func (o *ORMEngine) FindNoCtx(result interface{}, query *builder.Builder) (interface{}, error) {
	startTime := time.Now()

	cmd, args, err := query.ToSQL()
	if err != nil {
		xhlog.AppError(xhlog.OPMysql, map[string]interface{}{
			xhlog.ErrorMsg: fmt.Sprintf("query.ToSQL() fail, %s", err.Error()),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, err
	}

	err = o.Engine.SQL(cmd, args...).Find(result)
	if err != nil {
		xhlog.AppError(xhlog.OPMysql, map[string]interface{}{
			"sql":          cmd,
			xhlog.Args:     args,
			xhlog.ErrorMsg: err.Error(),
			xhlog.Res:      "fail",
			xhlog.ProcTime: xhlog.GetProcTime(startTime),
		})
		return nil, err
	}

	xhlog.AppInfo(xhlog.OPMysql, map[string]interface{}{
		"sql":          cmd,
		xhlog.Args:     args,
		xhlog.Res:      "success",
		xhlog.Result:   result,
		xhlog.ProcTime: xhlog.GetProcTime(startTime),
	})
	return result, nil
}

// todo: 事物相关的待实现
