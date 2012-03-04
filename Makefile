DIRS = help xpath xml html

all:
	set -e; for d in $(DIRS); do make -C $$d ; done

test:
	set -e; for d in $(DIRS); do make test -C $$d ; done

install:
	set -e; for d in $(DIRS); do make install -C $$d ; done

clean:
	set -e; for d in $(DIRS); do make clean -C $$d ; done

