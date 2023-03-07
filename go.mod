module example.com/main

go 1.19

replace example.com/test => ./modules

require example.com/test v0.0.0-00010101000000-000000000000

require golang.org/x/crypto v0.7.0 // indirect
