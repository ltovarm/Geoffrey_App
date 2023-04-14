# .PHONY: run stop

# run: 
# 	docker-compose up -d --build

# stop: 
# 	docker-compose down

CC=g++
CFLAGS=-Wall -Wextra -I /usr/include
LDFLAGS=-L/usr/lib/x86_64-linux-gnu -lrabbitmq
OBJDIR:=obj
BUILDDIR:=bld
SENDER_SRCDIR := sender
SENDER_OBJDIR := $(OBJDIR)/sender
CONSUMER_OBJDIR := $(OBJDIR)/consumer
CONSUMER_SRCDIR := consumer
SENDER_SOURCES=$(SENDER_SRCDIR)/sender.cpp
CONSUMER_SOURCES=$(CONSUMER_SRCDIR)/consumer.cpp
SENDER_OBJECTS=$(SENDER_SOURCES:%.cpp=$(OBJDIR)/%.o)
CONSUMER_OBJECTS=$(CONSUMER_SOURCES:%.cpp=$(OBJDIR)/%.o)
SENDER_EXECUTABLE=sender_program
CONSUMER_EXECUTABLE=consumer_program

build: create_dir $(SENDER_SOURCES) $(CONSUMER_SOURCES) $(SENDER_EXECUTABLE) $(CONSUMER_EXECUTABLE) clean_object

create_dir:
	@echo "Creando directorios $(CONSUMER_OBJDIR) y $(SENDER_OBJDIR)..."
	mkdir -p $(CONSUMER_OBJDIR) $(SENDER_OBJDIR)
	mkdir -p $(BUILDDIR)

$(SENDER_EXECUTABLE): $(SENDER_OBJECTS)
	$(CC) $(SENDER_OBJECTS) -o ${BUILDDIR}/$@ $(LDFLAGS)

$(CONSUMER_EXECUTABLE): $(CONSUMER_OBJECTS)
	$(CC) $(CONSUMER_OBJECTS) -o ${BUILDDIR}/$@ $(LDFLAGS)

$(SENDER_OBJDIR)/%.o: $(SENDER_SRCDIR)/%.cpp
	$(CC) $(CFLAGS) -c $< -o $@

$(CONSUMER_OBJDIR)/%.o: $(CONSUMER_SRCDIR)/%.cpp
	$(CC) $(CFLAGS) -c $< -o $@

clean_object:
	rm -rf $(OBJDIR)
clean:
	rm -f $(SENDER_EXECUTABLE)$(CONSUMER_EXECUTABLE)