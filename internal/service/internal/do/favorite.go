// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorite is the golang structure of table favorite for DAO operations like Where/Data.
type Favorite struct {
	g.Meta       `orm:"table:favorite, do:true"`
	Id           interface{} //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	DeletedAt    *gtime.Time //
	FavoriteId   interface{} //
	FavoriteById interface{} //
}