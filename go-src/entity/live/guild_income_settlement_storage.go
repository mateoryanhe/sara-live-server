package entity

import "xr-game-server/constants/db"

const (
	TbCoinMerchantGuildIncomeSettlementLog      db.TbName = "coin_merchant_guild_income_settlement_logs"
	TbCoinMerchantGuildIncomeSettlementDetail   db.TbName = "coin_merchant_guild_income_settlement_details"
	TbCoinMerchantGuildIncomeSettlementTransfer db.TbName = "coin_merchant_guild_income_settlement_transfers"
)

// GuildIncomeSettlementStorage 工会结算流水物理表分组（普通工会 / 币商工会）。
type GuildIncomeSettlementStorage uint8

const (
	GuildIncomeSettlementStorageNormal       GuildIncomeSettlementStorage = 0
	GuildIncomeSettlementStorageCoinMerchant GuildIncomeSettlementStorage = 1
)

type guildIncomeSettlementTables struct {
	Log      db.TbName
	Detail   db.TbName
	Transfer db.TbName
}

// Tables 返回该分组对应的主表、明细表、代付表。
func (s GuildIncomeSettlementStorage) Tables() (log, detail, transfer db.TbName) {
	t := s.tables()
	return t.Log, t.Detail, t.Transfer
}

func (s GuildIncomeSettlementStorage) tables() guildIncomeSettlementTables {
	if s == GuildIncomeSettlementStorageCoinMerchant {
		return guildIncomeSettlementTables{
			Log:      TbCoinMerchantGuildIncomeSettlementLog,
			Detail:   TbCoinMerchantGuildIncomeSettlementDetail,
			Transfer: TbCoinMerchantGuildIncomeSettlementTransfer,
		}
	}
	return guildIncomeSettlementTables{
		Log:      TbGuildIncomeSettlementLog,
		Detail:   TbGuildIncomeSettlementDetail,
		Transfer: TbGuildIncomeSettlementTransfer,
	}
}

func (r *GuildIncomeSettlementLog) storage() GuildIncomeSettlementStorage {
	if r == nil {
		return GuildIncomeSettlementStorageNormal
	}
	return r.Storage
}

func (r *GuildIncomeSettlementLog) logTable() db.TbName {
	return r.storage().tables().Log
}

func (r *GuildIncomeSettlementLog) detailTable() db.TbName {
	return r.storage().tables().Detail
}

func (r *GuildIncomeSettlementLog) transferTable() db.TbName {
	return r.storage().tables().Transfer
}

// GuildIncomeSettlementStorageForGuildType 将 CMS 工会类型筛选映射到结算表分组。
func GuildIncomeSettlementStorageForGuildType(guildType *uint8) *GuildIncomeSettlementStorage {
	if guildType == nil {
		return nil
	}
	if *guildType == LiveGuildTypeCoinMerchant {
		v := GuildIncomeSettlementStorageCoinMerchant
		return &v
	}
	v := GuildIncomeSettlementStorageNormal
	return &v
}
