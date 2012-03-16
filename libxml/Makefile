subdirs = help tree xpath

all: install

install:
	set -e; for d in $(subdirs); do make -C $$d install ; done
	make -f libxml.mak install

nuke:
	make -f libxml.mak nuke
	set -e; for d in $(subdirs); do make -C $$d nuke ; done

clean:
	make -f libxml.mak clean
	set -e; for d in $(subdirs); do make -C $$d clean ; done

runtest:
	set -e; cd test; make test

bench:
	#echo "no benchmark tests for the old gokogiri"
