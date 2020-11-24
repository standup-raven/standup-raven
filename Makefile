GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH=amd64
GOFLAGS ?= $(GOFLAGS:)
MANIFEST_FILE ?= plugin.json

define GetPluginId
$(shell node -p "require('./plugin.json').id")
endef

define GetPluginVersion
$(shell node -p "'v' + require('./plugin.json').version")
endef

define AddTimeZoneOptions
$(shell node -e 
"
let fs = require('fs');
try {
	let manifest = fs.readFileSync('plugin.json', 'utf8'); 
	manifest = JSON.parse(manifest);
	let timezones = fs.readFileSync('timezones.json', 'utf8'); 
	timezones = JSON.parse(timezones); 
	manifest.settings_schema.settings[0].options=timezones; 
	let json = JSON.stringify(manifest, null, 2);
	fs.writeFileSync('plugin.json', json, 'utf8'); 
} catch (err) {
	console.log(err);
};"
)
endef

define RemoveTimeZoneOptions
$(shell node -e 
"
let fs = require('fs');
try {
	let manifest = fs.readFileSync('plugin.json', 'utf8'); 
	manifest = JSON.parse(manifest);
	manifest.settings_schema.settings[0].options=[]; 
	let json = JSON.stringify(manifest, null, 2);
	fs.writeFileSync('plugin.json', json, 'utf8'); 
} catch (err) {
	console.log(err);
};"
)
endef


PLUGINNAME=$(call GetPluginId)
PLUGINVERSION=$(call GetPluginVersion)
PACKAGENAME=mattermost-plugin-$(PLUGINNAME)-$(PLUGINVERSION)

.PHONY: default build test run clean stop check-style check-style-server .distclean dist fix-style release

.SILENT: default build test run clean stop check-style check-style-server .distclean dist fix-style release inithashes buildwebapp buildserver

default: check-style test dist

check-style: check-style-server check-style-webapp

check-style-webapp: .webinstall
	@echo Checking for style guide compliance
	cd webapp && yarn run lintjs
	cd webapp && yarn run lintstyle

check-style-server:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
    		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
    		exit 1; \
    	fi; \
    
	@echo Running golangci-lint
	golangci-lint run ./server/...
	
fix-style: check-style-server
	@echo Checking for style guide compliance
	cd webapp && yarn run fixjs
	cd webapp && yarn run fixstyle
	
test-server: vendor
	@echo Running server tests
	go test -v -coverprofile=coverage.txt ./...

test: test-server

coverage: test-server
	go tool cover -html=coverage.txt -o coverage.html

.webinstall: webapp/yarn.lock
	@echo Getting webapp dependencies

	cd webapp && yarn install

vendor: go.sum
	@echo "Downloading server dependencies"
	go mod download

inithashes:
ifeq (,$(wildcard ./server.sha))
	@echo "Initializing server hash file"
	@$(call UpdateServerHash)
endif
ifeq (,$(wildcard ./webapp.sha))
	@echo "Initializing webapp hash file"
	@$(call UpdateWebappHash)
endif

prequickdist: plugin.json
	@echo Updating plugin.json with timezones
	$(call AddTimeZoneOptions)
    
doquickdist: inithashes buildwebapp buildserver package
	@echo $(PLUGINNAME)
	@echo $(PACKAGENAME)
	@echo $(PLUGINVERSION)

	@echo Quick building plugin

buildserver:
	cp server.sha server.old.sha
	echo "Updating server hash"
	$(call UpdateServerHash)
	FILES_MATCH=true;\
	if cmp -s "server.sha" "server.old.sha"; then\
		FILES_MATCH=true;\
	else\
		FILES_MATCH=false;\
	fi;\
	DIST_DIR="./dist";\
	export DIST_EXISTS=true;\
	if [ -d "DIST_DIR" ]; then\
		export DIST_EXISTS=true;\
	else\
		export DIST_EXISTS=false;\
	fi;\
	if $$FILES_MATCH && $$DIST_EXISTS; then\
		echo "Skipping server build as nothing updated since last build.";\
	else\
		# Build files from server\
		# We need to disable gomodules when installing gox to prevent `go get` from updating go.mod file.\
		# See this for more details -\
		# 	https://stackoverflow.com/questions/56842385/using-go-get-to-download-binaries-without-adding-them-to-go-mod\
		cd server;\
		GO111MODULE=off go get github.com/mitchellh/gox;\
		cd ..;\
		$(shell go env GOPATH)/bin/gox -ldflags="-X 'main.PluginVersion=$(PLUGINVERSION)' -X 'main.SentryServerDSN=$(SERVER_DSN)' -X 'main.SentryWebappDSN=$(WEBAPP_DSN)' -X 'main.EncodedPluginIcon=data:image/svg+xml;base64,`base64 webapp/src/assets/images/logo.svg`' " -osarch='darwin/amd64 linux/amd64 windows/amd64' -gcflags='all=-N -l' -output 'dist/intermediate/plugin_{{.OS}}_{{.Arch}}' ./server;\
	fi
	rm server.old.sha

buildwebapp:
	cp webapp.sha webapp.old.sha
	echo "Updating webapp hash"
	$(call UpdateWebappHash)
	FILES_MATCH=true;\
	if cmp -s "webapp.sha" "webapp.old.sha"; then\
    	FILES_MATCH=true;\
    else\
    	FILES_MATCH=false;\
    fi;\
    DIST_DIR="./dist";\
    export DIST_EXISTS=true;\
    if [ -d "DIST_DIR" ]; then\
    	export DIST_EXISTS=true;\
    else\
    	export DIST_EXISTS=false;\
    fi;\
    if $$FILES_MATCH && $$DIST_EXISTS; then\
    	echo "Skipping webapp build as nothing updated since last build.";\
    else\
    	cd webapp;\
    	yarn run build;\
    	cd ..;\
    	mkdir -p dist/$(PLUGINNAME)/webapp;\
    	cp -r webapp/dist/* dist/$(PLUGINNAME)/webapp/;\
    fi
	rm webapp.old.sha

package:
	mkdir -p dist/$(PLUGINNAME)
	
	# Copy plugin files
	cp plugin.json dist/$(PLUGINNAME)/

	# Copy server executables & compress plugin
	mkdir -p dist/$(PLUGINNAME)/server

	mv dist/intermediate/plugin_darwin_amd64 dist/$(PLUGINNAME)/server/plugin.exe
	cd dist && tar -zcvf $(PACKAGENAME)-darwin-amd64.tar.gz $(PLUGINNAME)/*

#	mv dist/intermediate/plugin_linux_amd64 dist/$(PLUGINNAME)/server/plugin.exe
#	cd dist && tar -zcvf $(PACKAGENAME)-linux-amd64.tar.gz $(PLUGINNAME)/*
#
#	mv dist/intermediate/plugin_windows_amd64.exe dist/$(PLUGINNAME)/server/plugin.exe
#	cd dist && tar -zcvf $(PACKAGENAME)-windows-amd64.tar.gz $(PLUGINNAME)/*

	# Clean up temp files
	#rm -rf dist/$(PLUGINNAME)
	#rm -rf dist/intermediate

	@echo Linux plugin built at: dist/$(PACKAGENAME)-linux-amd64.tar.gz
	@echo MacOS X plugin built at: dist/$(PACKAGENAME)-darwin-amd64.tar.gz
	@echo Windows plugin built at: dist/$(PACKAGENAME)-windows-amd64.tar.gz

postquickdist:
	@echo Remove data from plugin.json
	$(call RemoveTimeZoneOptions)
	
quickdist: prequickdist doquickdist postquickdist

dist: vendor .webinstall quickdist
	@echo Building plugin

run: .webinstall
	@echo Not yet implemented

stop:
	@echo Not yet implemented

clean: .distclean
	@echo Cleaning plugin

	rm -rf webapp/node_modules
	rm -rf webapp/.npminstall

# deploy installs the plugin to a (development) server, using the API if appropriate environment
# variables are defined, or copying the files directly to a sibling mattermost-server directory
.PHONY: deploy
deploy:
	@echo "Installing plugin via API"

	@echo "Authenticating admin user..." && \
	TOKEN=`http --print h POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/users/login login_id=$(MM_ADMIN_USERNAME) password=$(MM_ADMIN_PASSWORD) X-Requested-With:"XMLHttpRequest" | grep Token | cut -f2 -d' '` && \
	http GET $(MM_SERVICESETTINGS_SITEURL)/api/v4/users/me Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Deleting existing plugin..." && \
	http DELETE $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins/$(PLUGINNAME) Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Uploading plugin..." && \
	http --check-status --form POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins plugin@dist/$(PACKAGENAME)-$(PLATFORM)-amd64.tar.gz Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Enabling uploaded plugin..." && \
	http POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins/$(PLUGINNAME)/enable Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Logging out admin user..." && \
	http POST $(MM_SERVICESETTINGS_SITEURL)/api/v4/users/logout Authorization:"Bearer $$TOKEN" > /dev/null && \
	echo "Plugin uploaded successfully"

release: dist
	@echo "Installing ghr"
	@go get -u github.com/tcnksm/ghr
	@echo "Create new tag"
	$(shell git tag $(PLUGINVERSION))
	@echo "Uploading artifacts"
	@ghr -t $(GITHUB_TOKEN) -u $(ORG_NAME) -r $(REPO_NAME) $(PLUGINVERSION) dist/
