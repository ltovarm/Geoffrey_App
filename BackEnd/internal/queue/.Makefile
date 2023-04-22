run-docker: 
	docker-compose up -d --build

stop-docker: 
	docker-compose down

CC=g++
CFLAGS=-Wall -fPIC
WSDIR := /app
SRCDIR = $(WSDIR)/src
LIBDIR = $(WSDIR)/lib
OBJDIR = $(WSDIR)/obj
BUILDDIR = $(WSDIR)/bld
EXECUTABLES=$(SOURCES:$(SRCDIR)/%.cpp=$(BUILDDIR)/%)
INCLUDE := -I$(WSDIR)/lib -L$(WSDIR)/lib -lqueuermq
SOURCES=$(wildcard $(SRCDIR)/*.cpp)
OBJECTS=$(SOURCES:$(SRCDIR)/%.cpp=$(OBJDIR)/%.o)
LIBRARIES=$(wildcard $(LIBDIR)/*.cpp)
LIBOBJECTS=$(LIBRARIES:%.cpp=%.o)
TARGET_LIB = $(LIBDIR)/lib$(notdir $(LIBRARIES:%.cpp=%.so))
LDFLAGS=-shared -Wl,-soname,lib$(notdir $(LIBRARIES:%.cpp=%.so))
LDRABBIT=-L/usr/lib/x86_64-linux-gnu -lrabbitmq

# .PHONY: all 
# all:
# 	echo "--->($(EXECUTABLES))"

build: create_dirs $(LIBOBJECTS) $(TARGET_LIB) $(EXECUTABLES) clean_object
create_dirs:
	@echo "Creating directories..."
	mkdir -p $(OBJDIR)
	mkdir -p $(LIBDIR)
	mkdir -p $(BUILDDIR)


$(LIBOBJECTS): $(LIBRARIES)
	$(CC) $(CFLAGS) -c $< -o $@
$(TARGET_LIB): $(LIBOBJECTS)
	$(CC) $(LDFLAGS) -o $@ $< $(LDRABBIT)
$(OBJECTS): $(SOURCES)
	$(CC) $(CFLAGS) $(INCLUDE) -c $< -o $@
$(EXECUTABLES): $(BUILDDIR)/% : $(SRCDIR)/%.cpp
	$(CC) $(CFLAGS) -o $@ $< $(INCLUDE)
clean_object:
	rm -rf $(OBJDIR) $(LIBDIR)/*.o
clean:
	rm -rf $(OBJDIR) $(BUILDDIR) $(LIBDIR)/*.so