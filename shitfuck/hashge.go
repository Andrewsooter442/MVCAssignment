// hashgen.go
package main

import (
    "fmt"
    "log"
    "os"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run hashgen.go <password>")
        return
    }

    password := os.Args[1]
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Password:", password)
    fmt.Println("Hashed Password:", string(hashedPassword))
}
