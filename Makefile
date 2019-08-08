test:
	go test -v 

testcover:
	go test -coverprofile=covprofile
	go tool cover -html=covprofile -o coverage.html

clean:
	rm -f covprofile
	rm -f coverage.html

