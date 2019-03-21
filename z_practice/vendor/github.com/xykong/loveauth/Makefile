ifeq ($(OS),Windows_NT)
    CCFLAGS += -D WIN32
    SWAGGER = ./scripts/swagger.exe
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
        SWAGGER = ./scripts/swagger-Linux
    endif
    ifeq ($(UNAME_S),Darwin)
        CCFLAGS += -D OSX
        SWAGGER = ./scripts/swagger-Darwin
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif

SOURCE_DIR=.

default : loveauth

release : version default test doc-only


FILES = $(foreach dir,$(DIRS),$(wildcard $(dir)/*.go))
SOURCES := $(shell find $(SOURCE_DIR) -name '*.go')

loveauth : $(SOURCES)
	make doc-only
	go build

.PHONY : clean all version test docker

docker :
	docker build -t loveauth .
	if (( `docker ps -a -q | wc -l` > 0 )); then docker rm -f `docker ps -a -q`; fi
	docker run -d -p 10022:22 --name "loveauth" -h "`cat VERSION`" loveauth

version :
	./scripts/version.sh VERSION ./cmd/version.go ./server/server.go README.md

test :
	go test `go list ./... | grep -v /vendor/`

clean :
	-rm loveauth

air :
	./scripts/air -c .air.conf

doc-only : $(SOURCES)
	$(SWAGGER) generate spec -o ./docs/swagger.json

doc : doc-only
	$(SWAGGER) serve ./docs/swagger.json -p 8888

publish : release
	git commit -a -m "publish version `cat VERSION`."
#	git push
	echo "publish version `cat VERSION` success."

deploy :
	ssh auth 'cd /disk/workspace/go/src/github.com/xykong/loveauth && git pull'
	ssh auth 'cd /disk/workspace/go/src/github.com/xykong/loveauth && go build'
	ssh auth 'cd /disk/workspace/go/src/github.com/xykong/loveauth && ./scripts/reload.sh'

deploy-build :
	ssh wol-build 'source ~/.zshrc; cd $${GOPATH}/src/github.com/xykong/loveauth && git pull'
	ssh wol-build 'source ~/.zshrc; cd $${GOPATH}/src/github.com/xykong/loveauth && go build'
	ssh wol-build 'source ~/.zshrc; cd $${GOPATH}/src/github.com/xykong/loveauth && ./scripts/reload.sh'

excels=~/Documents/workspace/LemonGrassMedia/proj4/RD/trunk/excels

excels:
	rm -rf templates/shop_activity.xlsx templates/shop.xlsx templates/maskingwords.txt templates/waring_config.xlsx
	cp $(excels)/shop_activity.xlsx $(excels)/shop.xlsx $(excels)/maskingwords.txt $(excels)/waring_config.xlsx templates/
	git add * || echo ok.
	git commit -a -m "Update excels to $$(svn info $(excels) | grep 'Last Changed Rev')" || echo ok.
