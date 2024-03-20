package model

import (
	"fmt"
	otgorm "github.com/eddycjy/opentracing-gorm"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/pkg/setting"
	"time"
)

// Model 创建公共 model
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
	ModifiedBy string `json:"modified_by,omitempty"`
	CreatedOn  uint32 `json:"created_on,omitempty"`
	ModifiedOn uint32 `json:"modified_on,omitempty"`
	DeletedOn  uint32 `json:"deleted_on,omitempty"`
	IsDel      uint8  `json:"is_del,omitempty"`
}

const (
	STATE_OPEN  = 1
	STATE_CLOSE = 0
)

// NewDBEngine 创建DB实例
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	if global.ServerSetting.RunMode == "debug" {
		//启用Logger，显示详细日志
		db.LogMode(true)
		//// 禁用日志记录器，不显示任何日志
		//db.LogMode(false)
	}
	//gorm查找struct名对应数据库中的表名的时候会默认把你的struct中的大写字母转换为小写并加上“s”
	db.SingularTable(true) //创建生成的表名不带s，让gorm转义struct名字的时候不用加上“s”
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	//针对下面的三个 Callback 方法进行回调注册，才能够让我们的应用程序真正的使用上
	//db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	//db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	//db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)

	//sql追踪
	otgorm.AddGormCallbacks(db)

	return db, nil
}

// 新增行为 的回调
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		//scope.FieldByName，获取当前是否包含所需的字段。
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok { //ok 字段存在
			if createTimeField.IsBlank { //字段值为空 （字段不存在!=值为空）
				_ = createTimeField.Set(nowTime) //给该字段设置值
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

// 更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	//获取当前设置了标识 gorm:update_column 的字段属性。
	if _, ok := scope.Get("gorm:update_column"); !ok {
		//不存在，也就是没有自定义设置 update_column
		_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除行为的回调
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		//判断是否存在 DeletedOn 和 IsDel 字段
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField { //若存在
			now := time.Now().Unix()
			//执行 UPDATE 操作进行软删除（修改 DeletedOn 和 IsDel 的值）
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(), //获取当前所引用的表名
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now), //scope.AddToVars，该方法可以添加值作为SQL的参数，也可用于防范SQL注入
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()), //scope.CombinedConditionSql，完成 SQL 语句的组装
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			//不存在，执行 DELETE 进行硬删除
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}
func addExtraSpaceIfExist(str string) string {
	if str != "" { //没看懂？？？
		return " " + str
	}
	return ""
}
