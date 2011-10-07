CFLAGS=$(shell pkg-config --cflags libxml-2.0)
LDLIBS=$(shell pkg-config --libs libxml-2.0)

CFLAGS += -I/usr/local/include
LDLIBS += -lcppunit

OBJS = Chelpers_test.o Chelpers.o
DEPS = Chelpers.h

test: $(OBJS) $(DEPS)
	g++ -o test $(OBJS) $(LDLIBS) -lxml2 -lcppunit

%.o: %.c $(DEPS)
	gcc -c -o $@ $< $(CFLAGS)

%.o: %.cpp $(DEPS)
	g++ -c -o $@ $< $(CFLAGS)

clean:
	rm *.o test
