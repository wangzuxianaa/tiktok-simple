package repository

type Favorite struct {
	Id         int64 `gorm:"primarykey"` //视频id 唯一标识
	OperatorId int64 `gorm:"not null"`
	DeleteFlag string
}

func (f *Favorite) CreateFavorite() error {
	err := db.Debug().Create(f).Error
	return err
}

func (v *Video) FindVideoByVideoId() error { //通过视频的id视频从数据库中查找出来
	tx := db.Debug().Where("id = ?", v.Id).Find(&v)
	err := tx.Error
	if err != nil {
		return err
	}
	return nil
}

func (f *Favorite) DeleteFavoriteByDeleteFlag() error { //删除操作要通过操作用户id和视频id确定--->播放地址唯一的 操作者id唯一的 将两者组合成一个字段
	err := db.Debug().Where("delete_flag = ?", f.DeleteFlag).Delete(&f).Error //是否能这样删除？
	return err
}

func (f *Favorite) FindvideosByOperationId() ([]*Favorite, error) {
	var favoritevideo []*Favorite
	//是否能这样查找？
	tx := db.Debug().Where("operator_id = ?", f.OperatorId).Find(&favoritevideo)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return favoritevideo, nil
}
