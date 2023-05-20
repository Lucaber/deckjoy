
REMOTE_TARGET := deck@192.168.1.156
REMOTE_TARGET_ROOT := root@192.168.1.156

.PHONY: all
all: zip upload cleanup-remote run-remote-steam

.PHONY: daemon
daemon: build upload cleanup-remote daemon-remote

.PHONY: build
build:
	go build -ldflags "-X 'github.com/lucaber/deckjoy/pkg/config.Version=0.0.1'" -o deckjoy main.go

.PHONY: zip
zip: build
	-rm deckjoy.zip
	mkdir -p build
	cp deckjoy build
	cp assets/* build
	cp README.md build
	cd build && zip ../deckjoy.zip *
	rm -r build

.PHONY: upload
upload:
	rsync deckjoy.zip $(REMOTE_TARGET):deckjoy.zip
	ssh $(REMOTE_TARGET) unzip -u -o deckjoy.zip -d deckjoy
	rsync deckjoy $(REMOTE_TARGET_ROOT):deckjoy

.PHONY: run-remote-steam
run-remote-steam:
	ssh $(REMOTE_TARGET) steam steam://rungameid/15557043747182084096

.PHONY: run-remote
run-remote:
	ssh $(REMOTE_TARGET) ./deckjoy gui

.PHONY: daemon-remote
daemon-remote:
	ssh $(REMOTE_TARGET_ROOT) ./deckjoy daemon

.PHONY: cleanup-remote
cleanup-remote:
	-ssh $(REMOTE_TARGET_ROOT) killall deckjoy
	sleep 1
	-ssh $(REMOTE_TARGET_ROOT) ./deckjoy cleanup
