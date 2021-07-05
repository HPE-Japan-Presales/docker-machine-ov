# Go Prams
GOCMD=go
GOBUILD=$(GOCMD) build -trimpath
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
CURRENT_VERSION=$(shell git describe --tags --abbrev=0)
BUILD_TARGET="cmd/machine-driver.go"
BUILD_PATH="./bin/"
BUILD_BASE_NAME=docker-machine-driver-ov

all: test build
testsum:
	@echo "Execute test code..."
	@OV_DEBUG=TRUE gotestsum -f testname
test:
	@echo "Execute test code..."
ifdef TEST_TARGET
	@echo "== Test: $(TEST_TARGET) ================"
	@OV_DEBUG=TRUE $(GOTEST) -v ./driver -run $(TEST_TARGET)
else
	@echo "== Test: All ======================"
	@OV_DEBUG=TRUE $(GOTEST) -v ./driver
endif
build:
	@$(GOCLEAN) all
	@echo Version:$(CURRENT_VERSION)
	@mkdir -p $(BUILD_PATH)
	@echo "== Build for Windows amd64"
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_PATH)$(BUILD_BASE_NAME) -ldflags "-X main.version=$(CURRENT_VERSION)" $(BUILD_TARGET)
	@tar zcvf $(BUILD_PATH)$(BUILD_BASE_NAME)-$(CURRENT_VERSION)-windows-amd64.tar.gz -C $(BUILD_PATH) $(BUILD_BASE_NAME)
	@echo "== Build for OSX amd64"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_PATH)$(BUILD_BASE_NAME) -ldflags "-X main.version=$(CURRENT_VERSION)" $(BUILD_TARGET)
	@tar zcvf $(BUILD_PATH)$(BUILD_BASE_NAME)-$(CURRENT_VERSION)-darwin-amd64.tar.gz -C $(BUILD_PATH) $(BUILD_BASE_NAME)
	@echo "== Build for Linux amd64"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $(BUILD_PATH)$(BUILD_BASE_NAME) -ldflags "-X main.version=$(CURRENT_VERSION)" $(BUILD_TARGET)
	@tar zcvf $(BUILD_PATH)$(BUILD_BASE_NAME)-$(CURRENT_VERSION)-linux-amd64.tar.gz -C $(BUILD_PATH) $(BUILD_BASE_NAME)
	@rm $(BUILD_PATH)$(BUILD_BASE_NAME)