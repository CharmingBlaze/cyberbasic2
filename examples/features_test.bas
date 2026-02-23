// Comments: // and we have /* block */ too
Print "Sin(0)=", Sin(0)
Print "Cos(0)=", Cos(0)
Print "Lerp(0, 10, 0.5)=", Lerp(0, 10, 0.5)
Print "Sqrt(4)=", Sqrt(4)
Print "Noise(1, 2)=", Noise(1, 2)

let x = 2
SELECT CASE x
CASE 1
    Print "one"
CASE 2
    Print "two"
CASE 3
    Print "three"
CASE ELSE
    Print "other"
END SELECT

// File I/O: OpenFile(path, mode) mode 0=read 1=write 2=append
// h = OpenFile("out.txt", 1)
// WriteLine(h, "hello")
// CloseFile(h)

Quit
Print "never"