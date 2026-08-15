package cfgdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/snowflake"
	"xr-game-server/entity/user"
)

func CountAllRandomNicknames() (int, error) {
	n, err := g.DB().Model(string(entity.TbRandomNickname)).Count()
	return int(n), err
}

func CountRandomNicknamesByLang(lang uint8) (int, error) {
	n, err := g.DB().Model(string(entity.TbRandomNickname)).Where("lang", lang).Count()
	return int(n), err
}

func LoadAllRandomNicknames() ([]*entity.RandomNickname, error) {
	var rows []*entity.RandomNickname
	err := g.DB().Model(string(entity.TbRandomNickname)).Order("lang asc, id asc").Scan(&rows)
	return rows, err
}

func DeleteRandomNicknamesByLang(lang uint8) error {
	_, err := g.DB().Model(string(entity.TbRandomNickname)).Where("lang", lang).Delete()
	return err
}

func ListRandomNicknamesByLang(lang uint8) ([]string, error) {
	var rows []struct {
		Nickname string
	}
	err := g.DB().Model(string(entity.TbRandomNickname)).
		Fields("nickname").
		Where("lang", lang).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Nickname != "" {
			out = append(out, row.Nickname)
		}
	}
	return out, nil
}

func BatchInsertRandomNicknames(lang uint8, names []string) error {
	if len(names) == 0 {
		return nil
	}
	now := time.Now()
	rows := make([]*entity.RandomNickname, 0, len(names))
	for _, name := range names {
		row := &entity.RandomNickname{
			Nickname: name,
			Lang:     lang,
		}
		row.ID = snowflake.GetId()
		row.CreatedAt = now
		row.UpdatedAt = now
		rows = append(rows, row)
	}
	_, err := g.DB().Model(string(entity.TbRandomNickname)).Data(rows).Batch(len(rows)).Insert()
	return err
}

func ReplaceRandomNicknamesByLang(lang uint8, names []string) error {
	if err := DeleteRandomNicknamesByLang(lang); err != nil {
		return err
	}
	return BatchInsertRandomNicknames(lang, names)
}
