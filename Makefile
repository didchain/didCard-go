BINDIR=bin

all: i a

a:
	gomobile bind -v -x -o $(BINDIR)/dss.aar -target=android github.com/didchain/didCard-go/android
	#gomobile bind -v -x -a ../../bls-eth-go-binary/android/obj/local/armeabi-v7a/libbls384_256.a -o $(BINDIR)/dss.aar -target=android github.com/didchain/didCard-go/android
i:
	gomobile bind -v -o $(BINDIR)/iosLib.framework -target=ios github.com/didchain/didCard-go/ios

clean:
	gomobile clean
	rm $(BINDIR)/*
