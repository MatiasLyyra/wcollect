.PHONY: wcollect
wcollect:
	go build -o wcollect ./bin/

.PHONY: install
install:
	install ./wcollect /usr/local/bin