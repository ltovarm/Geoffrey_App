# .PHONY: run stop

# run: 
# 	docker-compose up -d --build

# stop: 
# 	docker-compose down

CC=g++
CFLAGS=-Wall -fPIC
SRCDIR := src
LIBDIR := lib
OBJDIR:=obj
BUILDDIR:=bld
EXECUTABLES=$(SOURCES:$(SRCDIR)/%.cpp=$(BUILDDIR)/%)
INCLUDE := -I/home/ltovar/workspace/Geoffrey_App/BackEnd/internal/queue/lib -L/home/ltovar/workspace/Geoffrey_App/BackEnd/internal/queue/lib -lqueuermq
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

build: create_dirs $(LIBOBJECTS) $(TARGET_LIB) $(EXECUTABLES) clean
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
	rm -rf $(OBJDIR)
clean:
	rm -rf $(OBJDIR) $(LIBDIR)/*.o