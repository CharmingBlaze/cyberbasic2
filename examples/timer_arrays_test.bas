DIM a(3, 4) AS Integer
a(0, 0) = 10
a(1, 2) = 99
Print "a(0,0) = ", a(0, 0)
Print "a(1,2) = ", a(1, 2)

ResetTimer()
Print "Timer ~0: ", Timer()
Sleep 100
Print "Timer after Sleep 100ms: ", Timer()

Print "Random(): ", Random()
Print "Random(10): ", Random(10)
Print "Int(3.7) = ", Int(3.7)
Print "Done."