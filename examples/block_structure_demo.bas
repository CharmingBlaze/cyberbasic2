// Block structure demo: END IF, ENDFUNCTION, END SUB, ELSEIF, END FUNCTION, END SUB, END MODULE

FUNCTION calculate(x, y, op)
    IF op = "add" THEN
        RETURN x + y
    ELSEIF op = "multiply" THEN
        RETURN x * y
    ELSE
        RETURN 0
    END IF
ENDFUNCTION

SUB PrintHeader(title)
    PRINT "=========="
    PRINT title
    PRINT "=========="
END SUB

MODULE Helpers
    FUNCTION double(n)
        RETURN n * 2
    END FUNCTION
    SUB greet(name)
        PRINT "Hello, "; name
    END SUB
END MODULE

// Main
PRINT calculate(3, 4, "add")
PRINT calculate(2, 5, "multiply")
PrintHeader("Block Structure Demo")
PRINT "Done."
