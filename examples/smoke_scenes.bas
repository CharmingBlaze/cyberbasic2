// Smoke: flat scene API (CreateScene / LoadScene / GetCurrentScene).
// cyberbasic --lint examples/smoke_scenes.bas
CreateScene("demo")
LoadScene("demo")
VAR cur = GetCurrentScene()
Print cur
