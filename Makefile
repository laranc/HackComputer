SCRIPT = 

default:
	go build -o ./build/HackComputer ./cmd/. 

run:
	go run ./cmd/. $(SCRIPT)
