GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH=amd64
GOFLAGS ?= $(GOFLAGS:)
MANIFEST_FILE ?= plugin.json

define GetFromPkg
$(shell node -p "require('./build_properties.json').$(1)")
endef

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

default: check-style test dist

check-style: check-style-server check-style-webapp

check-style-webapp: .npminstall
	@echo Checking for style guide compliance
	cd webapp && yarn run lint

check-style-server:
	@echo Running GOFMT

	@for package in $$(go list ./server/...); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -w -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \
	
fix-style: check-style-server
	@echo Checking for style guide compliance
	cd webapp && yarn run fix
	
test-server: vendor
	@echo Running server tests
	go test -v -coverprofile=coverage.txt ./...

test: test-server

cover: test-server
	go tool cover -html=coverage.txt -o coverage.html

.npminstall: webapp/package-lock.json
	@echo Getting dependencies using npm

	cd webapp && yarn install

vendor: server/glide.lock
	cd server && go get github.com/Masterminds/glide
	cd server && $(shell go env GOPATH)/bin/glide install

prequickdist: .distclean plugin.json
	@echo Updating plugin.json with timezones
	$(call AddTimeZoneOptions)
    
doquickdist: 
	@echo $(PLUGINNAME)
	@echo $(PACKAGENAME)
	@echo $(PLUGINVERSION)

	@echo Quick building plugin

	# Build and copy files from webapp
	cd webapp && yarn run build
	mkdir -p dist/$(PLUGINNAME)/webapp
	cp -r webapp/dist/* dist/$(PLUGINNAME)/webapp/

	# Build files from server
	 cd server && go get github.com/mitchellh/gox
	 $(shell go env GOPATH)/bin/gox -ldflags="-X main.SentryEnabled=$(call GetFromPkg,sentry.enabled) -X main.SentryDSN=$(call GetFromPkg,sentry.dsn)" -osarch='darwin/amd64 linux/amd64 windows/amd64' -gcflags='all=-N -l' -output 'dist/intermediate/plugin_{{.OS}}_{{.Arch}}' ./server

	# Copy plugin files
	cp plugin.json dist/$(PLUGINNAME)/

	# Copy server executables & compress plugin
	mkdir -p dist/$(PLUGINNAME)/server

	mv dist/intermediate/plugin_darwin_amd64 dist/$(PLUGINNAME)/server/plugin.exe
	cd dist && tar -zcvf $(PACKAGENAME)-darwin-amd64.tar.gz $(PLUGINNAME)/*

	mv dist/intermediate/plugin_linux_amd64 dist/$(PLUGINNAME)/server/plugin.exe
	cd dist && tar -zcvf $(PACKAGENAME)-linux-amd64.tar.gz $(PLUGINNAME)/*

	mv dist/intermediate/plugin_windows_amd64.exe dist/$(PLUGINNAME)/server/plugin.exe
	cd dist && tar -zcvf $(PACKAGENAME)-windows-amd64.tar.gz $(PLUGINNAME)/*

	# Clean up temp files
	rm -rf dist/$(PLUGINNAME)
	rm -rf dist/intermediate

	@echo Linux plugin built at: dist/$(PACKAGENAME)-linux-amd64.tar.gz
	@echo MacOS X plugin built at: dist/$(PACKAGENAME)-darwin-amd64.tar.gz
	@echo Windows plugin built at: dist/$(PACKAGENAME)-windows-amd64.tar.gz

postquickdist:
	@echo Remove data from plugin.json
	$(call RemoveTimeZoneOptions)
	
quickdist: prequickdist doquickdist postquickdist

dist: vendor .npminstall quickdist
	@echo Building plugin

run: .npminstall
	@echo Not yet implemented

stop:
	@echo Not yet implemented

.distclean:
	@echo Cleaning dist files

	rm -rf dist
	rm -rf webapp/dist
	rm -f server/plugin.exe

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
