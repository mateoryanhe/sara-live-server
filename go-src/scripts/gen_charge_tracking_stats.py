from pathlib import Path

root = Path(__file__).resolve().parents[1]
src_tpl = root / "entity/stat/daily_mini_game_exposure_stat.go"
text_tpl = src_tpl.read_text(encoding="utf-8")

events = [
    ("MiniGameExposure", "ChargePanelEnter", "mini_game_exposure", "charge_panel_enter", "充值面板曝光"),
    ("MiniGameExposure", "ChargeClick", "mini_game_exposure", "charge_click", "充值档位点击建单"),
    ("MiniGameExposure", "ChargeResultServer", "mini_game_exposure", "charge_result_server", "充值到账成功"),
]

for pascal, new_pascal, old_snake, new_snake, comment in events:
    pairs = [
        (pascal, new_pascal),
        (old_snake, new_snake),
        ("半屏游戏窗口曝光", comment),
        ("曝光次数", "次数"),
    ]
    for period in ("daily", "weekly", "monthly"):
        for layer, sub in (
            ("entity/stat", f"{period}_{old_snake}_stat.go"),
            ("dao/statdao", f"{period}_{old_snake}_stat.go"),
        ):
            t = (root / layer / sub).read_text(encoding="utf-8")
            for a, b in pairs:
                t = t.replace(a, b)
            dst = root / layer / f"{period}_{new_snake}_stat.go"
            dst.write_text(t, encoding="utf-8")
            print(dst.relative_to(root))
