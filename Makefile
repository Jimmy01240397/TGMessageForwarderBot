builder := go
builddir := bin
exe := $(builddir)/tgmsgforwarderbot
config := .env.sample
#install := $(builddir)/install.sh
#systemd := $(builddir)/piproxyclient.service

all: $(exe) $(builddir)/$(config) #$(install) $(systemd)

$(builddir)/$(config): $(config)
		cp $(config) $(builddir)/$(config)

#$(install): install.sh
#		cp install.sh $(install)

#$(systemd): piproxyclient.service
#		cp piproxyclient.service $(systemd)

$(exe): main.go go.mod go.sum models router utils
		$(builder) build -o $(exe) $<

.PHONY = clean

clean: 
		rm -r $(builddir)
