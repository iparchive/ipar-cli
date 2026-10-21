ipar_name := ipar
ipar_src  := cmd/ipar/main.go
ipar_bin  := $(BINARY_DIR)/$(ipar_name)

BUILD_TARGETS += build-ipar
CLEAN_TARGETS += clean-ipar
INSTALL_TARGETS += install-ipar

.PHONY: build-ipar clean-ipar run-ipar

build-ipar: $(ipar_bin)

$(ipar_bin): $(GO_FILES)
	@echo "🔨 Building $(ipar_name)..."
	@go build -ldflags "$(LDFLAGS)" -o $(ipar_bin) $(ipar_src)

run-ipar: build-ipar
	@./$(ipar_bin)

clean-ipar:
	@echo "Removing $(ipar_name) binaries..."
	@rm -f $(ipar_bin) ~/bin/$(shell basename $(ipar_bin))

install-ipar:
	@echo "Installing $(ipar_bin) binaries..."
	@cp $(ipar_bin) ~/bin
