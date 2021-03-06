package mysql

import (
	"database/sql"
	"pinkacg/models"
)

func GetFollowUserById(followId int64, uid int64) (num []*models.Follow, err error) {
	sqlStr := `select user_id,follow_id from follows where user_id=? and follow_id=?`
	err = db.Select(&num, sqlStr, uid, followId)
	return
}

func GetFollowsUserById(uid int64) (follows []*models.Follow, err error) {
	sqlStr := `select user_id,follow_id from follows where user_id=?`
	err = db.Select(&follows, sqlStr, uid)
	return
}

func GetFollowsUserByFollowId(uid int64) (follows []*models.Follow, err error) {
	sqlStr := `select user_id,follow_id from follows where follow_id=?`
	err = db.Select(&follows, sqlStr, uid)
	return
}

func DeleteFollowUserById(followId int64, uid int64) (res sql.Result, err error) {
	// 查询关注数量
	sqlStr4 := `select user_id,follow_id from follows where user_id=?`
	var followsNum []*models.Follow
	err = db.Select(&followsNum, sqlStr4, uid)
	if err != nil {
		return nil, err
	}
	// 查询粉丝数量
	sqlStr5 := `select user_id,follow_id from follows where follow_id=?`
	var fansNum []*models.Follow
	err = db.Select(&fansNum, sqlStr5, followId)
	if err != nil {
		return nil, err
	}

	// 开启事务
	tx, err := db.Begin()
	// 删除关注数据
	sqlStr := `delete from follows where user_id=? and follow_id=?`
	res, err = tx.Exec(sqlStr, uid, followId)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	// 更新user表当前用户关注数量
	sqlStr2 := `update users set follows=?,update_time=NOW() where user_id=?`
	res, err = tx.Exec(sqlStr2, len(followsNum)-1, uid)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	// 更新user表被关注用户粉丝数量
	sqlStr3 := `update users set fans=?,update_time=NOW() where user_id=?`
	res, err = tx.Exec(sqlStr3, len(fansNum)-1, followId)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit()
	return
}

// InsertFollowUserById 插入关注
func InsertFollowUserById(followId int64, uid int64) (res sql.Result, err error) {
	// 查询关注数量
	sqlStr4 := `select user_id,follow_id from follows where user_id=?`
	var followsNum []*models.Follow
	err = db.Select(&followsNum, sqlStr4, uid)
	if err != nil {
		return nil, err
	}
	// 查询粉丝数量
	sqlStr5 := `select user_id,follow_id from follows where follow_id=?`
	var fansNum []*models.Follow
	err = db.Select(&fansNum, sqlStr5, followId)
	if err != nil {
		return nil, err
	}

	// 开启事务
	tx, err := db.Begin()
	// 插入关注数据
	sqlStr := `insert into follows (user_id,follow_id,update_time,create_time) values(?,?,NOW(),NOW())`
	res, err = tx.Exec(sqlStr, uid, followId)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	// 更新user表当前用户关注数量
	sqlStr2 := `update users set follows=?,update_time=NOW() where user_id=?`
	res, err = tx.Exec(sqlStr2, len(followsNum)+1, uid)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	// 更新user表被关注用户粉丝数量
	sqlStr3 := `update users set fans=?,update_time=NOW() where user_id=?`
	res, err = tx.Exec(sqlStr3, len(fansNum)+1, followId)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit()
	return
}
