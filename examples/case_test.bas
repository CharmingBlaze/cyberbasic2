DIM x As Integer
x = 10
Print x
If x > 5 Then
    Print 99
EndIf
DIM a As FLOAT
DIM b As Boolean
a = 1.5
b = true
If a > 1 AND b Then
    Print 42
EndIf
bullet.CreateWorld("w", 0, -9.81, 0)
bullet.Step("w", 0.1)
