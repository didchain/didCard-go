BINDIR=bin

all: i a

a:
	gomobile bind -v -o $(BINDIR)/dss.aar -target=android github.com/didchain/didCard-go/android
i:
	gomobile bind -v -o $(BINDIR)/iosLib.framework -target=ios github.com/didchain/didCard-go/ios

clean:
	gomobile clean
	rm $(BINDIR)/*
