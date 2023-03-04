package storage

import (
	"context"
	"trading-service/common"
	"trading-service/modules/mastertx/model"
)

func (s *sqlStore) Create(ctx context.Context, data common.ExtendMasterTxData) error {
	db := s.db.Begin()

	if inData := data.MoneyInData(); inData != nil {
		newTx := model.MasterTx{
			SQLModel:     common.NewSQLModel(),
			UserId:       inData.UserId(),
			Amount:       inData.Amount(),
			Type:         "in",
			WalletId:     inData.WalletId(),
			AssetId:      inData.AssetId(),
			RelatedId:    data.GetRelatedId(),
			RelatedTable: data.GetRelatedTable(),
		}

		if err := db.Table(newTx.TableName()).Create(&newTx).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if outData := data.MoneyOutData(); outData != nil {
		newTx := model.MasterTx{
			SQLModel:     common.NewSQLModel(),
			UserId:       outData.UserId(),
			Amount:       outData.Amount(),
			Type:         "in",
			WalletId:     outData.WalletId(),
			AssetId:      outData.AssetId(),
			RelatedId:    data.GetRelatedId(),
			RelatedTable: data.GetRelatedTable(),
		}

		if err := db.Table(newTx.TableName()).Create(&newTx).Error; err != nil {
			db.Rollback()
			return common.ErrDB(err)
		}
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
