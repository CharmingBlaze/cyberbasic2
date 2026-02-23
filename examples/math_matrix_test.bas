Print "Floor(2.7)=", Floor(2.7)
Print "Ceil(2.2)=", Ceil(2.2)
Print "Round(2.5)=", Round(2.5)
Print "Min(3, 7)=", Min(3, 7)
Print "Max(3, 7)=", Max(3, 7)
Print "Clamp(15, 0, 10)=", Clamp(15, 0, 10)
Print "Pow(2, 10)=", Pow(2, 10)
Print "Len(hello)=", Len("hello")
Print "Left(hello, 2)=", Left("hello", 2)
Print "Right(hello, 2)=", Right("hello", 2)
Print "Mid(hello, 2, 2)=", Mid("hello", 2, 2)

x = 0
REPEAT
    x = x + 1
    Print "repeat x=", x
UNTIL x >= 3

DIM A(2, 2) AS Float
DIM B(2, 2) AS Float
DIM R(2, 2) AS Float
A(0, 0) = 1
A(0, 1) = 2
A(1, 0) = 3
A(1, 1) = 4
B(0, 0) = 1
B(0, 1) = 0
B(1, 0) = 0
B(1, 1) = 1
MatMul(R, A, B)
Print "R(0,0)=", R(0, 0), " R(0,1)=", R(0, 1)
Print "R(1,0)=", R(1, 0), " R(1,1)=", R(1, 1)
Print "Done."