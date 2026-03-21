// Smoke: v2 assets module (key/value store).
// cyberbasic --lint examples/smoke_assets_v2.bas
assets.set("ping", 1)
LET v = assets.get("ping")
PRINT v
