DIRS = util help xpath xml html libxml

all: install

build:
	set -e; for d in $(DIRS); do make -C $$d ; done

test:
	set -e; for d in $(DIRS); do make test -C $$d ; done
	set -e; cd libxml; make runtest

install:
	set -e; for d in $(DIRS); do make install -C $$d ; done

clean:
	set -e; for d in $(DIRS); do make clean -C $$d ; done

nuke:
	set -e; for d in $(DIRS); do make nuke -C $$d ; done
bench:
	set -e; for d in $(DIRS); do make bench -C $$d; done

