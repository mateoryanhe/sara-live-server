#!/usr/bin/env python3
"""Merge small module subdirectories into rank, cfg, cms, and liveroom."""
from __future__ import annotations

import os
import re
import shutil
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
MODULE = ROOT / "module"


def read_text(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def write_text(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8", newline="\n")


def replace_imports(content: str, mapping: dict[str, str]) -> str:
    for old, new in mapping.items():
        content = content.replace(f'"xr-game-server/module/{old}"', f'"xr-game-server/module/{new}"')
    return content


def create_rank_package() -> None:
    rank_dir = MODULE / "rank"
    rank_dir.mkdir(exist_ok=True)

    common = '''package rank

import "time"

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

func calcAge(birthday *time.Time) int {
	if birthday == nil || birthday.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	anniversary := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(anniversary) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

func pageRange(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
'''
    write_text(rank_dir / "common.go", common)

    init_go = '''package rank

func Init() {
	initRichRank()
	initAnchorRank()
	initGameConsumeRank()
}
'''
    write_text(rank_dir / "init.go", init_go)

    rich = read_text(MODULE / "richrank" / "rich_rank.go")
    rich = rich.replace("package richrank", "package rank")
    rich = rich.replace("type rankItem struct", "type richRankItem struct")
    rich = rich.replace("type rankSnapshot struct", "type richRankSnapshot struct")
    rich = rich.replace("*rankItem", "*richRankItem")
    rich = rich.replace("*rankSnapshot", "*richRankSnapshot")
    rich = rich.replace("var richRankCache", "var richRankCache")
    rich = rich.replace("dataRefreshDeadline", "richRankRefreshDeadline")
    rich = rich.replace("func Init()", "func initRichRank()")
    rich = rich.replace("markRichRankDataChanged()", "markRichRankDataChanged()")
    rich = rich.replace("func getSnapshot()", "func getRichRankSnapshot()")
    rich = rich.replace("getSnapshot()", "getRichRankSnapshot()")
    rich = rich.replace("[]*richRankItem) []*richRankItem", "[]*currencylogdao.DiamondConsumeStatRow) []*richRankItem")
    rich = re.sub(
        r"func buildRankItems\(rows \[\]\*currencylogdao\.DiamondConsumeStatRow\) \[\]\*richRankItem",
        "func buildRichRankItems(rows []*currencylogdao.DiamondConsumeStatRow) []*richRankItem",
        rich,
    )
    rich = rich.replace("buildRankItems(", "buildRichRankItems(")
    rich = rich.replace("func getRankListByPeriod(snapshot *richRankSnapshot", "func getRichRankListByPeriod(snapshot *richRankSnapshot")
    rich = rich.replace("getRankListByPeriod(snapshot", "getRichRankListByPeriod(snapshot")
    rich = rich.replace("func onRankListRefreshEvent", "func onRichRankListRefreshEvent")
    rich = rich.replace("onRankListRefreshEvent", "onRichRankListRefreshEvent")
    rich = re.sub(r"\nconst \(\n\tdefaultPageSize.*?richRankRefreshDelay = 10 \* time.Minute\n\)\n", "\n", rich, flags=re.S)
    rich = re.sub(r"\nfunc normalizePage\(.*?\n\}\n", "\n", rich, flags=re.S)
    rich = re.sub(r"\nfunc calcAge\(.*?\n\}\n", "\n", rich, flags=re.S)
    rich = re.sub(r"\nfunc pageRange\(.*?\n\}\n", "\n", rich, flags=re.S)
    rich = re.sub(r"\nfunc startOfDay\(.*?\n\}\n", "\n", rich, flags=re.S)
    write_text(rank_dir / "rich_rank.go", rich)

    anchor = read_text(MODULE / "anchorrank" / "anchor_rank.go")
    anchor = anchor.replace("package anchorrank", "package rank")
    anchor = anchor.replace("type rankItem struct", "type anchorRankItem struct")
    anchor = anchor.replace("type rankSnapshot struct", "type anchorRankSnapshot struct")
    anchor = anchor.replace("*rankItem", "*anchorRankItem")
    anchor = anchor.replace("*rankSnapshot", "*anchorRankSnapshot")
    anchor = anchor.replace("dataRefreshDeadline", "anchorRankRefreshDeadline")
    anchor = anchor.replace("func Init()", "func initAnchorRank()")
    anchor = anchor.replace("func getSnapshot()", "func getAnchorRankSnapshot()")
    anchor = anchor.replace("getSnapshot()", "getAnchorRankSnapshot()")
    anchor = re.sub(
        r"func buildRankItems\(rows \[\]\*liveroomdao\.AnchorRevenueStatRow\) \[\]\*anchorRankItem",
        "func buildAnchorRankItems(rows []*liveroomdao.AnchorRevenueStatRow) []*anchorRankItem",
        anchor,
    )
    anchor = anchor.replace("buildRankItems(", "buildAnchorRankItems(")
    anchor = anchor.replace("func getRankListByPeriod(snapshot *anchorRankSnapshot", "func getAnchorRankListByPeriod(snapshot *anchorRankSnapshot")
    anchor = anchor.replace("getRankListByPeriod(snapshot", "getAnchorRankListByPeriod(snapshot")
    anchor = anchor.replace("func onRankListRefreshEvent", "func onAnchorRankListRefreshEvent")
    anchor = anchor.replace("onRankListRefreshEvent", "onAnchorRankListRefreshEvent")
    anchor = re.sub(r"\nconst \(\n\tdefaultPageSize.*?anchorRankRefreshDelay = 10 \* time.Minute\n\)\n", "\n", anchor, flags=re.S)
    anchor = re.sub(r"\nfunc normalizePage\(.*?\n\}\n", "\n", anchor, flags=re.S)
    anchor = re.sub(r"\nfunc calcAge\(.*?\n\}\n", "\n", anchor, flags=re.S)
    anchor = re.sub(r"\nfunc pageRange\(.*?\n\}\n", "\n", anchor, flags=re.S)
    anchor = re.sub(r"\nfunc startOfDay\(.*?\n\}\n", "\n", anchor, flags=re.S)
    write_text(rank_dir / "anchor_rank.go", anchor)

    game = read_text(MODULE / "gameconsumrank" / "game_consume_rank.go")
    game = game.replace("package gameconsumrank", "package rank")
    game = game.replace("type rankItem struct", "type gameConsumeRankItem struct")
    game = game.replace("type rankSnapshot struct", "type gameConsumeRankSnapshot struct")
    game = game.replace("*rankItem", "*gameConsumeRankItem")
    game = game.replace("*rankSnapshot", "*gameConsumeRankSnapshot")
    game = game.replace("dataRefreshDeadline", "gameConsumeRankRefreshDeadline")
    game = game.replace("func Init()", "func initGameConsumeRank()")
    game = game.replace("func getSnapshot()", "func getGameConsumeRankSnapshot()")
    game = game.replace("getSnapshot()", "getGameConsumeRankSnapshot()")
    game = re.sub(
        r"func buildRankItems\(rows \[\]\*currencylogdao\.DiamondConsumeStatRow\) \[\]\*gameConsumeRankItem",
        "func buildGameConsumeRankItems(rows []*currencylogdao.DiamondConsumeStatRow) []*gameConsumeRankItem",
        game,
    )
    game = game.replace("buildRankItems(", "buildGameConsumeRankItems(")
    game = game.replace("func getRankListByPeriod(snapshot *gameConsumeRankSnapshot", "func getGameConsumeRankListByPeriod(snapshot *gameConsumeRankSnapshot")
    game = game.replace("getRankListByPeriod(snapshot", "getGameConsumeRankListByPeriod(snapshot")
    game = game.replace("func onRankListRefreshEvent", "func onGameConsumeRankListRefreshEvent")
    game = game.replace("onRankListRefreshEvent", "onGameConsumeRankListRefreshEvent")
    game = re.sub(r"\nconst \(\n\tdefaultPageSize.*?gameConsumeRankRefreshDelay  = 10 \* time.Minute\n\)\n", "\n", game, flags=re.S)
    game = re.sub(r"\nfunc normalizePage\(.*?\n\}\n", "\n", game, flags=re.S)
    game = re.sub(r"\nfunc calcAge\(.*?\n\}\n", "\n", game, flags=re.S)
    game = re.sub(r"\nfunc pageRange\(.*?\n\}\n", "\n", game, flags=re.S)
    game = re.sub(r"\nfunc startOfDay\(.*?\n\}\n", "\n", game, flags=re.S)
    write_text(rank_dir / "game_consume_rank.go", game)


def create_cfg_package() -> None:
    cfg_dir = MODULE / "cfg"
    cfg_dir.mkdir(exist_ok=True)

    write_text(cfg_dir / "format.go", '''package cfg

import "time"

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
''')

    account_memory = read_text(MODULE / "accountcfg" / "memory.go").replace("package accountcfg", "package cfg")
    account_memory = account_memory.replace("type cfgSnapshot struct", "type accountCfgSnapshot struct")
    account_memory = account_memory.replace("*cfgSnapshot", "*accountCfgSnapshot")
    account_memory = account_memory.replace("var cfgCache", "var accountCfgCache")
    account_memory = account_memory.replace("reloadCfgMemory", "reloadAccountCfgMemory")
    account_memory = account_memory.replace("getCfgCache", "getAccountCfgCache")
    account_memory = account_memory.replace("toCfgSnapshot", "toAccountCfgSnapshot")
    write_text(cfg_dir / "account_memory.go", account_memory)

    account_cms = read_text(MODULE / "accountcfg" / "cfg_cms.go").replace("package accountcfg", "package cfg")
    account_cms = account_cms.replace("toCfgItem(", "toAccountCfgItem(")
    account_cms = account_cms.replace("func toAccountCfgItem", "func toAccountCfgItem")
    account_cms = account_cms.replace("reloadCfgMemory()", "reloadAccountCfgMemory()")
    account_cms = re.sub(r"\nfunc formatTime\(.*?\n\}\n", "\n", account_cms, flags=re.S)
    account_cms = account_cms.replace("formatTime(", "formatCfgTime(")
    write_text(cfg_dir / "account_cms.go", account_cms)

    live_memory = read_text(MODULE / "livecfg" / "memory.go").replace("package livecfg", "package cfg")
    write_text(cfg_dir / "live_memory.go", live_memory)

    live_cms = read_text(MODULE / "livecfg" / "cfg_cms.go").replace("package livecfg", "package cfg")
    live_cms = live_cms.replace("formatLiveCfgTime(", "formatCfgTime(")
    live_cms = re.sub(r"\nfunc formatLiveCfgTime\(.*?\n\}\n", "\n", live_cms, flags=re.S)
    write_text(cfg_dir / "live_cms.go", live_cms)

    privacy_memory = read_text(MODULE / "privacypolicy" / "memory.go").replace("package privacypolicy", "package cfg")
    privacy_memory = privacy_memory.replace("type cfgSnapshot struct", "type privacyPolicyCfgSnapshot struct")
    privacy_memory = privacy_memory.replace("*cfgSnapshot", "*privacyPolicyCfgSnapshot")
    privacy_memory = privacy_memory.replace("var (\n\tcfgCache", "var (\n\tprivacyPolicyCfgCache")
    privacy_memory = privacy_memory.replace("emptyCfgSnapshot = &cfgSnapshot{}", "emptyPrivacyPolicyCfgSnapshot = &privacyPolicyCfgSnapshot{}")
    privacy_memory = privacy_memory.replace("emptyCfgSnapshot", "emptyPrivacyPolicyCfgSnapshot")
    privacy_memory = privacy_memory.replace("reloadCfgMemory", "reloadPrivacyPolicyCfgMemory")
    privacy_memory = privacy_memory.replace("getCfgCache", "getPrivacyPolicyCfgCache")
    privacy_memory = privacy_memory.replace("toCfgSnapshot", "toPrivacyPolicyCfgSnapshot")
    write_text(cfg_dir / "privacy_policy_memory.go", privacy_memory)

    write_text(cfg_dir / "privacy_policy_url.go", read_text(MODULE / "privacypolicy" / "url.go").replace("package privacypolicy", "package cfg"))

    privacy_cms = read_text(MODULE / "privacypolicy" / "cfg_cms.go").replace("package privacypolicy", "package cfg")
    privacy_cms = privacy_cms.replace("toCfgItem(", "toPrivacyPolicyCfgItem(")
    privacy_cms = privacy_cms.replace("reloadCfgMemory()", "reloadPrivacyPolicyCfgMemory()")
    privacy_cms = privacy_cms.replace("formatTime(", "formatCfgTime(")
    privacy_cms = re.sub(r"\nfunc formatTime\(.*?\n\}\n", "\n", privacy_cms, flags=re.S)
    write_text(cfg_dir / "privacy_policy_cms.go", privacy_cms)

    cs_cms = read_text(MODULE / "customerservice" / "cfg_cms.go").replace("package customerservice", "package cfg")
    cs_cms = cs_cms.replace("toCfgItem(", "toCustomerServiceCfgItem(")
    cs_cms = cs_cms.replace("formatTime(", "formatCfgTime(")
    cs_cms = re.sub(r"\nfunc formatTime\(.*?\n\}\n", "\n", cs_cms, flags=re.S)
    write_text(cfg_dir / "customer_service_cms.go", cs_cms)

    write_text(cfg_dir / "customer_service_app.go", read_text(MODULE / "customerservice" / "cfg_app.go").replace("package customerservice", "package cfg"))

    sim_memory = read_text(MODULE / "simulatorcpukeyword" / "memory.go").replace("package simulatorcpukeyword", "package cfg")
    sim_memory = sim_memory.replace("func Init()", "func initSimulatorCpuKeyword()")
    write_text(cfg_dir / "simulator_cpu_keyword_memory.go", sim_memory)

    write_text(cfg_dir / "simulator_cpu_keyword_cms.go", read_text(MODULE / "simulatorcpukeyword" / "cfg_cms.go").replace("package simulatorcpukeyword", "package cfg"))

    write_text(cfg_dir / "anchor_salary_cms.go", read_text(MODULE / "anchorsalarycfg" / "cfg_cms.go").replace("package anchorsalarycfg", "package cfg"))

    share_calc = read_text(MODULE / "liverevenuesharecfg" / "share_calc.go").replace("package liverevenuesharecfg", "package cfg")
    write_text(cfg_dir / "live_revenue_share_calc.go", share_calc)

    share_cms = read_text(MODULE / "liverevenuesharecfg" / "cfg_cms.go").replace("package liverevenuesharecfg", "package cfg")
    share_cms = share_cms.replace("formatCfgTime(", "formatCfgTime(")
    share_cms = re.sub(r"\nfunc formatCfgTime\(.*?\n\}\n", "\n", share_cms, flags=re.S)
    write_text(cfg_dir / "live_revenue_share_cms.go", share_cms)

    preload_init = read_text(MODULE / "preload" / "init.go").replace("package preload", "package cfg")
    preload_init = preload_init.replace("func Init()", "func initPreload()")
    write_text(cfg_dir / "preload_boot.go", preload_init)

    write_text(cfg_dir / "preload_recent_user.go", read_text(MODULE / "preload" / "recent_user.go").replace("package preload", "package cfg"))
    write_text(cfg_dir / "preload_cms.go", read_text(MODULE / "preload" / "cfg_cms.go").replace("package preload", "package cfg").replace("formatPreloadCfgTime(", "formatCfgTime("))

    init_go = '''package cfg

func Init() {
	initAccountCfg()
	initLiveCfg()
	initPrivacyPolicyCfg()
	initCustomerServiceCfg()
	initSimulatorCpuKeyword()
	initLiveRevenueShareCfg()
	initPreload()
}

func initAccountCfg() {
	reloadAccountCfgMemory()
}

func initLiveCfg() {
	reloadLiveCfgMemory()
}

func initPrivacyPolicyCfg() {
	reloadPrivacyPolicyCfgMemory()
}

func initCustomerServiceCfg() {
	cfgdaoInitCustomerService()
	cfgdaoReloadCustomerService()
}

func initLiveRevenueShareCfg() {
	cfgdaoInitLiveRevenueShare()
	cfgdaoReloadLiveRevenueShare()
}
'''
    init_go = init_go.replace("cfgdaoInitCustomerService", "cfgdao.InitCustomerServiceCfgDao")
    init_go = init_go.replace("cfgdaoReloadCustomerService", "cfgdao.ReloadCustomerServiceCfgCache")
    init_go = init_go.replace("cfgdaoInitLiveRevenueShare", "cfgdao.InitLiveRevenueShareCfgDao")
    init_go = init_go.replace("cfgdaoReloadLiveRevenueShare", "cfgdao.ReloadLiveRevenueShareCfgCache")
    init_go = 'package cfg\n\nimport "xr-game-server/dao/cfgdao"\n\n' + init_go.split("\n", 1)[1]
    write_text(cfg_dir / "init.go", init_go)


def merge_liveroom_satellites() -> None:
    util = '''package liveroom

import "strconv"

func parseUint64Filter(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
'''
    write_text(MODULE / "liveroom" / "cms_filter_util.go", util)

    mappings = [
        ("liverecord/app_list.go", "liveroom/live_record_app.go"),
        ("liverecord/cms_list.go", "liveroom/live_record_cms.go"),
        ("liverecord/cms_daily_effective_live_list.go", "liveroom/live_record_daily_effective_cms.go"),
        ("liverevenue/cms_list.go", "liveroom/revenue_log_cms.go"),
        ("incomesettlement/cms_list.go", "liveroom/income_settlement_cms.go"),
        ("livefollow/follow.go", "liveroom/follow.go"),
        ("livefollow/block.go", "liveroom/follow_block.go"),
        ("livefollow/push.go", "liveroom/follow_push.go"),
    ]
    for src_rel, dst_rel in mappings:
        src = MODULE / src_rel
        dst = MODULE / dst_rel
        content = read_text(src).replace("package liverecord", "package liveroom")
        content = content.replace("package liverevenue", "package liveroom")
        content = content.replace("package incomesettlement", "package liveroom")
        content = content.replace("package livefollow", "package liveroom")
        if dst_rel == "liveroom/live_record_cms.go":
            content = content.replace("func toCMSItem(", "func liveRecordToCMSItem(")
            content = content.replace("toCMSItem(", "liveRecordToCMSItem(")
            content = re.sub(r"\nfunc parseUint64Filter\(.*?\n\}\n", "\n", content, flags=re.S)
        if dst_rel == "liveroom/revenue_log_cms.go":
            content = content.replace("func toCMSItem(", "func revenueLogToCMSItem(")
            content = content.replace("toCMSItem(", "revenueLogToCMSItem(")
            content = re.sub(r"\nfunc parseUint64Filter\(.*?\n\}\n", "\n", content, flags=re.S)
        if dst_rel == "liveroom/income_settlement_cms.go":
            content = re.sub(r"\nfunc parseUint64Filter\(.*?\n\}\n", "\n", content, flags=re.S)
        if dst_rel in ("liveroom/follow.go", "liveroom/follow_block.go"):
            content = re.sub(r"\nfunc calcAge\(.*?\n\}\n", "\n", content, flags=re.S)
        write_text(dst, content)


def merge_cms_package() -> None:
    cms_dir = MODULE / "cms"
    cms_dir.mkdir(exist_ok=True)
    for src_dir, prefix in [("cmsuser", "user_"), ("cmsexport", "export_")]:
        src_path = MODULE / src_dir
        if not src_path.exists():
            continue
        for go_file in sorted(src_path.glob("*.go")):
            name = go_file.name
            if src_dir == "cmsuser":
                dst_name = name
            else:
                dst_name = name
            content = read_text(go_file).replace(f"package {src_dir}", "package cms")
            write_text(cms_dir / dst_name if src_dir == "cmsuser" else cms_dir / name, content)


def remove_old_dirs() -> None:
    old_dirs = [
        "richrank", "anchorrank", "gameconsumrank",
        "accountcfg", "livecfg", "privacypolicy", "customerservice",
        "simulatorcpukeyword", "anchorsalarycfg", "liverevenuesharecfg", "preload",
        "liverecord", "liverevenue", "incomesettlement", "livefollow",
        "cmsuser", "cmsexport",
    ]
    for name in old_dirs:
        path = MODULE / name
        if path.exists():
            shutil.rmtree(path)


def update_go_sources() -> None:
    import_map = {
        "richrank": "rank",
        "anchorrank": "rank",
        "gameconsumrank": "rank",
        "accountcfg": "cfg",
        "livecfg": "cfg",
        "privacypolicy": "cfg",
        "customerservice": "cfg",
        "simulatorcpukeyword": "cfg",
        "anchorsalarycfg": "cfg",
        "liverevenuesharecfg": "cfg",
        "preload": "cfg",
        "liverecord": "liveroom",
        "liverevenue": "liveroom",
        "incomesettlement": "liveroom",
        "livefollow": "liveroom",
        "cmsuser": "cms",
        "cmsexport": "cms",
    }
    symbol_map = {
        "richrank.": "rank.",
        "anchorrank.": "rank.",
        "gameconsumrank.": "rank.",
        "accountcfg.": "cfg.",
        "livecfg.": "cfg.",
        "privacypolicy.": "cfg.",
        "customerservice.": "cfg.",
        "simulatorcpukeyword.": "cfg.",
        "anchorsalarycfg.": "cfg.",
        "liverevenuesharecfg.": "cfg.",
        "preload.": "cfg.",
        "liverecord.": "liveroom.",
        "incomesettlement.": "liveroom.",
        "livefollow.": "liveroom.",
        "cmsuser.": "cms.",
        "cmsexport.": "cms.",
    }

    for root, _, files in os.walk(ROOT):
        if "vendor" in root.split(os.sep):
            continue
        for fname in files:
            if not fname.endswith(".go"):
                continue
            path = Path(root) / fname
            if path.parts[-2:] == ("scripts", "merge_modules.py"):
                continue
            content = read_text(path)
            original = content
            content = replace_imports(content, import_map)
            for old, new in symbol_map.items():
                content = content.replace(old, new)
            content = content.replace("liverevenue.GetCMSList", "liveroom.GetCMSList")
            if content != original:
                write_text(path, content)

    init_path = MODULE / "init.go"
    init = read_text(init_path)
    init = init.replace("\t\"xr-game-server/module/accountcfg\"\n", "")
    init = init.replace("\t\"xr-game-server/module/customerservice\"\n", "")
    init = init.replace("\t\"xr-game-server/module/gameconsumrank\"\n", "")
    init = init.replace("\t\"xr-game-server/module/livecfg\"\n", "")
    init = init.replace("\t\"xr-game-server/module/liverevenuesharecfg\"\n", "")
    init = init.replace("\t\"xr-game-server/module/preload\"\n", "")
    init = init.replace("\t\"xr-game-server/module/privacypolicy\"\n", "")
    init = init.replace("\t\"xr-game-server/module/richrank\"\n", "")
    init = init.replace("\t\"xr-game-server/module/simulatorcpukeyword\"\n", "")
    init = init.replace("\t\"xr-game-server/module/anchorrank\"\n", "")
    if "\"xr-game-server/module/cfg\"" not in init:
        init = init.replace(
            "\t\"xr-game-server/module/call\"\n",
            "\t\"xr-game-server/module/call\"\n\t\"xr-game-server/module/cfg\"\n",
        )
    if "\"xr-game-server/module/rank\"" not in init:
        init = init.replace(
            "\t\"xr-game-server/module/randomnick\"\n",
            "\t\"xr-game-server/module/randomnick\"\n\t\"xr-game-server/module/rank\"\n",
        )
    init = init.replace("\trichrank.Init()\n", "\trank.Init()\n")
    init = init.replace("\tgameconsumrank.Init()\n", "")
    init = init.replace("\tanchorrank.Init()\n", "")
    init = init.replace("\tpreload.Init()\n", "\tcfg.Init()\n")
    init = init.replace("\tlivecfg.Init()\n", "")
    init = init.replace("\tprivacypolicy.Init()\n", "")
    init = init.replace("\taccountcfg.Init()\n", "")
    init = init.replace("\tsimulatorcpukeyword.Init()\n", "")
    init = init.replace("\tcustomerservice.Init()\n", "")
    init = init.replace("\tliverevenuesharecfg.Init()\n", "")
    write_text(init_path, init)


def main() -> None:
    create_rank_package()
    create_cfg_package()
    merge_liveroom_satellites()
    merge_cms_package()
    remove_old_dirs()
    update_go_sources()
    print("module merge done")


if __name__ == "__main__":
    main()
